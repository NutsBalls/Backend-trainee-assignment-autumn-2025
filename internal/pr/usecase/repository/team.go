package repository

import (
	"context"

	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/domain"
)

//go:generate mockgen -destination=../mocks/mock_team_repository.go -package=mocks github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/usecase/repository TeamRepository
type TeamRepository interface {
	CreateTeam(ctx context.Context, teamName string) error
	GetTeam(ctx context.Context, teamName string) (*domain.Team, error)
	TeamExists(ctx context.Context, teamName string) (bool, error)
}
