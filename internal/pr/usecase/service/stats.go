package service

import (
	"context"

	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/domain"
	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/usecase/repository"
)

type StatsService struct {
	statsRepo repository.StatsRepository
}

func NewStatsService(statsRepo repository.StatsRepository) *StatsService {
	return &StatsService{
		statsRepo: statsRepo,
	}
}

func (s *StatsService) GetUserAssignmentStats(ctx context.Context) ([]domain.UserAssignmentStats, error) {
	return s.statsRepo.GetUserAssignmentStats(ctx)
}

func (s *StatsService) GetPRStats(ctx context.Context) (*domain.PRStats, error) {
	return s.statsRepo.GetPRStats(ctx)
}

func (s *StatsService) GetReviewerWorkload(ctx context.Context) ([]domain.ReviewerWorkload, error) {
	return s.statsRepo.GetReviewerWorkload(ctx)
}
