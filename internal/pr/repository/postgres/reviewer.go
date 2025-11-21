package postgres

import (
	"context"
	"fmt"

	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/domain"
	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/repository/postgres/sqlc"
)

type ReviewerRepository struct {
	queries *sqlc.Queries
}

func NewReviewerRepository(queries *sqlc.Queries) *ReviewerRepository {
	return &ReviewerRepository{queries: queries}
}

func (r *ReviewerRepository) AssignReviewer(ctx context.Context, prID, reviewerID string) error {
	err := r.queries.AddReviewer(ctx, sqlc.AddReviewerParams{
		PrID:       prID,
		ReviewerID: reviewerID,
	})
	if err != nil {
		return fmt.Errorf("assign reviewer: %w", err)
	}
	return nil
}

func (r *ReviewerRepository) ReplaceReviewer(ctx context.Context, prID, oldReviewerID, newReviewerID string) error {
	err := r.queries.ReplaceReviewer(ctx, sqlc.ReplaceReviewerParams{
		PrID:         prID,
		ReviewerID:   oldReviewerID,
		ReviewerID_2: newReviewerID,
	})
	if err != nil {
		return fmt.Errorf("replace reviewer: %w", err)
	}
	return nil
}

func (r *ReviewerRepository) IsReviewerAssigned(ctx context.Context, prID, reviewerID string) (bool, error) {
	assigned, err := r.queries.IsReviewerAssigned(ctx, sqlc.IsReviewerAssignedParams{
		PrID:       prID,
		ReviewerID: reviewerID,
	})
	if err != nil {
		return false, fmt.Errorf("check reviewer assigned: %w", err)
	}
	return assigned, nil
}

func (r *ReviewerRepository) GetAssignedReviewers(ctx context.Context, prID string) ([]string, error) {
	reviewers, err := r.queries.GetAssignedReviewers(ctx, prID)
	if err != nil {
		return nil, fmt.Errorf("get assigned reviewers: %w", err)
	}
	return reviewers, nil
}

func (r *ReviewerRepository) FindCandidatesForNewPR(ctx context.Context, teamName, authorID string) ([]string, error) {
	users, err := r.queries.GetActiveCandidatesForPR(ctx, sqlc.GetActiveCandidatesForPRParams{
		TeamName: teamName,
		UserID:   authorID,
	})
	if err != nil {
		return nil, fmt.Errorf("find candidates for new PR: %w", err)
	}

	result := make([]string, len(users))
	for i, u := range users {
		result[i] = u.UserID
	}
	return result, nil
}

func (r *ReviewerRepository) FindCandidatesForReassignment(ctx context.Context, teamName, authorID, prID string) ([]string, error) {
	candidates, err := r.queries.GetActiveCandidatesForReassignment(ctx, sqlc.GetActiveCandidatesForReassignmentParams{
		TeamName: teamName,
		UserID:   authorID,
		PrID:     prID,
	})
	if err != nil {
		return nil, fmt.Errorf("find candidates for reassignment: %w", err)
	}
	return candidates, nil
}

func (r *ReviewerRepository) ListPRsByReviewer(ctx context.Context, reviewerID string) ([]domain.PullRequestShort, error) {
	prs, err := r.queries.ListPullRequestsByReviewer(ctx, reviewerID)
	if err != nil {
		return nil, fmt.Errorf("list PRs by reviewer: %w", err)
	}

	result := make([]domain.PullRequestShort, len(prs))
	for i, pr := range prs {
		result[i] = domain.PullRequestShort{
			PullRequestID:   pr.PullRequestID,
			PullRequestName: pr.PullRequestName,
			AuthorID:        pr.AuthorID,
			Status:          domain.PRStatus(pr.Status),
		}
	}
	return result, nil
}
