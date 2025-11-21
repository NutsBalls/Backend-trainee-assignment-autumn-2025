package usecase

import (
	"context"

	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/domain"
	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/usecase/dto"
)

type PRUseCase interface {
	CreatePR(ctx context.Context, req dto.CreatePRRequest) (*domain.PullRequest, error)
	MergePR(ctx context.Context, req dto.MergePRRequest) (*domain.PullRequest, error)
	ReassignReviewer(ctx context.Context, req dto.ReassignReviewerRequest) (*dto.ReassignReviewerResponse, error)
	GetReviewerPRs(ctx context.Context, reviewerID string) ([]domain.PullRequestShort, error)
}

type TeamUseCase interface {
	CreateTeam(ctx context.Context, req dto.CreateTeamRequest) (*domain.Team, error)
	GetTeam(ctx context.Context, teamName string) (*domain.Team, error)
}

type UserUseCase interface {
	SetIsActive(ctx context.Context, req dto.SetUserIsActiveRequest) (*domain.User, error)
}
