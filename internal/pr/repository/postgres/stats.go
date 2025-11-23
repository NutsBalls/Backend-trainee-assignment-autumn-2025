package postgres

import (
	"context"
	"fmt"

	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/domain"
	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/repository/postgres/sqlc"
)

type StatsRepository struct {
	queries *sqlc.Queries
}

func NewStatsRepository(queries *sqlc.Queries) *StatsRepository {
	return &StatsRepository{queries: queries}
}

func (r *StatsRepository) GetUserAssignmentStats(ctx context.Context) ([]domain.UserAssignmentStats, error) {
	rows, err := r.queries.GetUserAssignmentStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("get user assignment stats: %w", err)
	}

	result := make([]domain.UserAssignmentStats, len(rows))
	for i, row := range rows {
		result[i] = domain.UserAssignmentStats{
			UserID:           row.UserID,
			Username:         row.Username,
			TeamName:         row.TeamName,
			AssignmentsCount: row.AssignmentsCount,
		}
	}
	return result, nil
}

func (r *StatsRepository) GetPRStats(ctx context.Context) (*domain.PRStats, error) {
	stats, err := r.queries.GetPRStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("get PR stats: %w", err)
	}

	return &domain.PRStats{
		TotalPRs:  stats.TotalPrs,
		OpenPRs:   stats.OpenPrs,
		MergedPRs: stats.MergedPrs,
	}, nil
}

func (r *StatsRepository) GetReviewerWorkload(ctx context.Context) ([]domain.ReviewerWorkload, error) {
	rows, err := r.queries.GetReviewerWorkload(ctx)
	if err != nil {
		return nil, fmt.Errorf("get reviewer workload: %w", err)
	}

	result := make([]domain.ReviewerWorkload, len(rows))
	for i, row := range rows {
		result[i] = domain.ReviewerWorkload{
			UserID:       row.UserID,
			Username:     row.Username,
			TeamName:     row.TeamName,
			OpenPRsCount: row.OpenPrsCount,
		}
	}
	return result, nil
}
