package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/domain"
	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/repository/postgres/sqlc"
	"github.com/jackc/pgx/v5"
)

type PRRepository struct {
	queries *sqlc.Queries
}

func NewPRRepository(queries *sqlc.Queries) *PRRepository {
	return &PRRepository{queries: queries}
}

func (r *PRRepository) CreatePR(ctx context.Context, pr *domain.PullRequest) error {
	_, err := r.queries.CreatePullRequest(ctx, sqlc.CreatePullRequestParams{
		PullRequestID:   pr.PullRequestID,
		PullRequestName: pr.PullRequestName,
		AuthorID:        pr.AuthorID,
	})
	if err != nil {
		if isPgUniqueViolation(err) {
			return domain.ErrPRAlreadyExists
		}
		return fmt.Errorf("create PR: %w", err)
	}
	return nil
}

func (r *PRRepository) GetPR(ctx context.Context, prID string) (*domain.PullRequest, error) {
	pr, err := r.queries.GetPullRequest(ctx, prID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrPRNotFound
		}
		return nil, fmt.Errorf("get PR: %w", err)
	}

	return &domain.PullRequest{
		PullRequestID:   pr.PullRequestID,
		PullRequestName: pr.PullRequestName,
		AuthorID:        pr.AuthorID,
		Status:          domain.PRStatus(pr.Status),
		CreatedAt:       &pr.CreatedAt,
		MergedAt:        pr.MergedAt,
	}, nil
}

func (r *PRRepository) GetPRWithReviewers(ctx context.Context, prID string) (*domain.PullRequest, error) {
	pr, err := r.GetPR(ctx, prID)
	if err != nil {
		return nil, err
	}

	reviewers, err := r.queries.GetAssignedReviewers(ctx, prID)
	if err != nil {
		return nil, err
	}

	pr.AssignedReviewers = reviewers
	return pr, nil
}

func (r *PRRepository) PRExists(ctx context.Context, prID string) (bool, error) {
	exists, err := r.queries.PRExists(ctx, prID)
	if err != nil {
		return false, fmt.Errorf("check PR exists: %w", err)
	}
	return exists, nil
}

func (r *PRRepository) MergePR(ctx context.Context, prID string) (*domain.PullRequest, error) {
	pr, err := r.queries.MergePullRequest(ctx, prID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrPRNotFound
		}
		return nil, fmt.Errorf("merge PR: %w", err)
	}

	reviewers, err := r.queries.GetAssignedReviewers(ctx, prID)
	if err != nil {
		return nil, err
	}

	return &domain.PullRequest{
		PullRequestID:     pr.PullRequestID,
		PullRequestName:   pr.PullRequestName,
		AuthorID:          pr.AuthorID,
		Status:            domain.PRStatus(pr.Status),
		AssignedReviewers: reviewers,
		CreatedAt:         &pr.CreatedAt,
		MergedAt:          pr.MergedAt,
	}, nil
}

func (r *PRRepository) GetPRAuthorID(ctx context.Context, prID string) (string, error) {
	authorID, err := r.queries.GetPRAuthorId(ctx, prID)
	if err != nil {
		return "", err
	}

	return authorID, nil
}
