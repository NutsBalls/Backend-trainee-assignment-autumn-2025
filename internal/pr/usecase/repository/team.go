package repository

import (
	"context"

	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/domain"
)

type TeamRepository interface {
	CreateTeam(ctx context.Context, teamName string) error
	GetTeam(ctx context.Context, teamName string) (*domain.Team, error)
	TeamExists(ctx context.Context, teamName string) (bool, error)
}
