package usecase

import (
	"context"

	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/domain"
)

type PRUseCase interface {
	CreatePR(ctx context.Context, req CreatePRRequest) (*domain.PullRequest, error)
	MergePR(ctx context.Context, req MergePRRequest) (*domain.PullRequest, error)
	ReassignReviewer(ctx context.Context, req ReassignReviewerRequest) (*ReassignReviewerResponse, error)
	GetReviewerPRs(ctx context.Context, reviewerID string) ([]domain.PullRequestShort, error)
}

type TeamUseCase interface {
	CreateTeam(ctx context.Context, req CreateTeamRequest) (*domain.Team, error)
	GetTeam(ctx context.Context, teamName string) (*domain.Team, error)
}

type UserUseCase interface {
	SetIsActive(ctx context.Context, req SetUserIsActiveRequest) (*domain.User, error)
}
