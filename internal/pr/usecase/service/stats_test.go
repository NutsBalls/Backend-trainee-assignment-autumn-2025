package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/domain"
	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/usecase/mocks"
)

func TestStatsService_GetUserAssignmentStats(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStatsRepo := mocks.NewMockStatsRepository(ctrl)
	service := NewStatsService(mockStatsRepo)
	ctx := context.Background()

	t.Run("success - return user assignment stats", func(t *testing.T) {
		mockStatsRepo.EXPECT().
			GetUserAssignmentStats(ctx).
			Return([]domain.UserAssignmentStats{
				{UserID: "u1", Username: "Alice", TeamName: "backend", AssignmentsCount: 3},
				{UserID: "u2", Username: "Bob", TeamName: "backend", AssignmentsCount: 1},
			}, nil).
			Times(1)

		result, err := service.GetUserAssignmentStats(ctx)
		require.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "u1", result[0].UserID)
		assert.Equal(t, int64(3), result[0].AssignmentsCount)
	})

	t.Run("error - repository error", func(t *testing.T) {
		mockStatsRepo.EXPECT().
			GetUserAssignmentStats(ctx).
			Return(nil, errors.New("db error")).
			Times(1)

		result, err := service.GetUserAssignmentStats(ctx)
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "db error")
	})
}

func TestStatsService_GetPRStats(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStatsRepo := mocks.NewMockStatsRepository(ctrl)
	service := NewStatsService(mockStatsRepo)
	ctx := context.Background()

	t.Run("success - return PR stats", func(t *testing.T) {
		mockStatsRepo.EXPECT().
			GetPRStats(ctx).
			Return(&domain.PRStats{TotalPRs: 5, OpenPRs: 2, MergedPRs: 3}, nil).
			Times(1)

		result, err := service.GetPRStats(ctx)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(5), result.TotalPRs)
		assert.Equal(t, int64(2), result.OpenPRs)
		assert.Equal(t, int64(3), result.MergedPRs)
	})

	t.Run("error - repository error", func(t *testing.T) {
		mockStatsRepo.EXPECT().
			GetPRStats(ctx).
			Return(nil, errors.New("db error")).
			Times(1)

		result, err := service.GetPRStats(ctx)
		require.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestStatsService_GetReviewerWorkload(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStatsRepo := mocks.NewMockStatsRepository(ctrl)
	service := NewStatsService(mockStatsRepo)
	ctx := context.Background()

	t.Run("success - return reviewer workload", func(t *testing.T) {
		mockStatsRepo.EXPECT().
			GetReviewerWorkload(ctx).
			Return([]domain.ReviewerWorkload{
				{UserID: "u1", Username: "Alice", TeamName: "backend", OpenPRsCount: 2},
				{UserID: "u2", Username: "Bob", TeamName: "backend", OpenPRsCount: 1},
			}, nil).
			Times(1)

		result, err := service.GetReviewerWorkload(ctx)
		require.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, int64(2), result[0].OpenPRsCount)
		assert.Equal(t, "Alice", result[0].Username)
	})

	t.Run("error - repository error", func(t *testing.T) {
		mockStatsRepo.EXPECT().
			GetReviewerWorkload(ctx).
			Return(nil, errors.New("db error")).
			Times(1)

		result, err := service.GetReviewerWorkload(ctx)
		require.Error(t, err)
		assert.Nil(t, result)
	})
}
