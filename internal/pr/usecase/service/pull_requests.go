package service

import (
	"context"
	"fmt"

	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/domain"
	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/usecase/dto"
	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/usecase/repository"
)

type PRService struct {
	uow repository.UnitOfWork
}

func NewPRService(uow repository.UnitOfWork) *PRService {
	return &PRService{uow: uow}
}

func (s *PRService) CreatePR(ctx context.Context, req dto.CreatePRRequest) (*domain.PullRequest, error) {
	if req.PullRequestID == "" {
		return nil, fmt.Errorf("pull_request_id is required")
	}
	if req.PullRequestName == "" {
		return nil, fmt.Errorf("pull_request_name is required")
	}
	if req.AuthorID == "" {
		return nil, fmt.Errorf("author_id is required")
	}

	exists, err := s.uow.PullRequests().PRExists(ctx, req.PullRequestID)
	if err != nil {
		return nil, fmt.Errorf("check PR exists: %w", err)
	}
	if exists {
		return nil, domain.ErrPRAlreadyExists
	}

	author, err := s.uow.Users().GetUser(ctx, req.AuthorID)
	if err != nil {
		return nil, err
	}

	var createdPR *domain.PullRequest
	err = s.uow.WithinTransaction(ctx, func(txCtx context.Context) error {
		pr := &domain.PullRequest{
			PullRequestID:   req.PullRequestID,
			PullRequestName: req.PullRequestName,
			AuthorID:        req.AuthorID,
			Status:          domain.PRStatusOpen,
		}
		if err := s.uow.PullRequests().CreatePR(txCtx, pr); err != nil {
			return fmt.Errorf("create PR: %w", err)
		}

		candidates, err := s.uow.Reviewers().FindCandidatesForNewPR(
			txCtx,
			author.TeamName,
			req.AuthorID,
		)
		if err != nil {
			return fmt.Errorf("find candidates: %w", err)
		}

		for _, candidateID := range candidates {
			if err := s.uow.Reviewers().AssignReviewer(txCtx, req.PullRequestID, candidateID); err != nil {
				return fmt.Errorf("assign reviewer %s: %w", candidateID, err)
			}
		}

		createdPR, err = s.uow.PullRequests().GetPRWithReviewers(txCtx, req.PullRequestID)
		return err
	})

	if err != nil {
		return nil, err
	}

	return createdPR, nil
}

func (s *PRService) MergePR(ctx context.Context, req dto.MergePRRequest) (*domain.PullRequest, error) {
	if req.PullRequestID == "" {
		return nil, fmt.Errorf("pull_request_id is required")
	}

	pr, err := s.uow.PullRequests().MergePR(ctx, req.PullRequestID)
	if err != nil {
		return nil, err
	}

	return pr, nil
}

func (s *PRService) ReassignReviewer(ctx context.Context, req dto.ReassignReviewerRequest) (*dto.ReassignReviewerResponse, error) {
	if req.PullRequestID == "" {
		return nil, fmt.Errorf("pull_request_id is required")
	}
	if req.OldReviewerID == "" {
		return nil, fmt.Errorf("old_reviewer_id is required")
	}

	pr, err := s.uow.PullRequests().GetPR(ctx, req.PullRequestID)
	if err != nil {
		return nil, err
	}

	if pr.Status == domain.PRStatusMerged {
		return nil, domain.ErrPRMerged
	}

	isAssigned, err := s.uow.Reviewers().IsReviewerAssigned(ctx, req.PullRequestID, req.OldReviewerID)
	if err != nil {
		return nil, fmt.Errorf("check reviewer assigned: %w", err)
	}
	if !isAssigned {
		return nil, domain.ErrReviewerNotAssigned
	}

	oldReviewer, err := s.uow.Users().GetUser(ctx, req.OldReviewerID)
	if err != nil {
		return nil, err
	}

	authorID, err := s.uow.PullRequests().GetPRAuthorID(ctx, req.PullRequestID)
	if err != nil {
		return nil, err
	}

	candidates, err := s.uow.Reviewers().FindCandidatesForReassignment(
		ctx,
		oldReviewer.TeamName,
		authorID,
		req.PullRequestID,
	)
	if err != nil {
		return nil, fmt.Errorf("find replacement candidates: %w", err)
	}

	if len(candidates) == 0 {
		return nil, domain.ErrNoCandidates
	}

	newReviewerID := candidates[0]

	if err := s.uow.Reviewers().ReplaceReviewer(ctx, req.PullRequestID, req.OldReviewerID, newReviewerID); err != nil {
		return nil, fmt.Errorf("replace reviewer: %w", err)
	}

	updatedPR, err := s.uow.PullRequests().GetPRWithReviewers(ctx, req.PullRequestID)
	if err != nil {
		return nil, fmt.Errorf("get updated PR: %w", err)
	}

	return &dto.ReassignReviewerResponse{
		PullRequest: updatedPR,
		ReplacedBy:  newReviewerID,
	}, nil
}

func (s *PRService) GetReviewerPRs(ctx context.Context, reviewerID string) ([]domain.PullRequestShort, error) {
	if reviewerID == "" {
		return nil, fmt.Errorf("reviewer_id is required")
	}

	prs, err := s.uow.Reviewers().ListPRsByReviewer(ctx, reviewerID)
	if err != nil {
		return nil, fmt.Errorf("list PRs by reviewer: %w", err)
	}

	return prs, nil
}
