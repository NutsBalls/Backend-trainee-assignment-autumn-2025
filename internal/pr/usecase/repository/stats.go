package repository

import (
	"context"

	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/domain"
)

//go:generate mockgen -destination=../mocks/mock_stats_repository.go -package=mocks github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/usecase/repository StatsRepository
type StatsRepository interface {
	GetUserAssignmentStats(ctx context.Context) ([]domain.UserAssignmentStats, error)
	GetPRStats(ctx context.Context) (*domain.PRStats, error)
	GetReviewerWorkload(ctx context.Context) ([]domain.ReviewerWorkload, error)
}
