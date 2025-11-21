package interfaces

import (
	"context"

	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/domain"
)

type PRRepository interface {
	CreatePR(ctx context.Context, pr *domain.PullRequest) error
	GetPR(ctx context.Context, prID string) (*domain.PullRequest, error)
	GetPRWithReviewers(ctx context.Context, prID string) (*domain.PullRequest, error)
	PRExists(ctx context.Context, prID string) (bool, error)
	MergePR(ctx context.Context, prID string) (*domain.PullRequest, error)
	GetPRAuthorID(ctx context.Context, prID string) (string, error)
}

type ReviewerRepository interface {
	AssignReviewer(ctx context.Context, prID, reviewerID string) error
	ReplaceReviewer(ctx context.Context, prID, oldReviewerID, newReviewerID string) error
	IsReviewerAssigned(ctx context.Context, prID, reviewerID string) (bool, error)
	GetAssignedReviewers(ctx context.Context, prID string) ([]string, error)
	FindCandidatesForNewPR(ctx context.Context, teamName, authorID string) ([]string, error)
	FindCandidatesForReassignment(ctx context.Context, teamName, authorID, prID string) ([]string, error)
	ListPRsByReviewer(ctx context.Context, reviewerID string) ([]domain.PullRequestShort, error)
}
