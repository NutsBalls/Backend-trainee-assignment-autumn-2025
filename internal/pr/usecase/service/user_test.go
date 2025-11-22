package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/domain"
	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/usecase/dto"
	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/usecase/mocks"
)

func TestUserService_SetIsActive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUOW := mocks.NewMockUnitOfWork(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	mockUOW.EXPECT().Users().Return(mockUserRepo).AnyTimes()

	service := NewUserService(mockUOW)
	ctx := context.Background()

	t.Run("success - set user active", func(t *testing.T) {
		req := dto.SetUserIsActiveRequest{
			UserID:   "u1",
			IsActive: true,
		}

		expectedUser := &domain.User{
			UserID:   "u1",
			Username: "Alice",
			TeamName: "backend",
			IsActive: true,
		}

		mockUserRepo.EXPECT().
			SetUserIsActive(ctx, "u1", true).
			Return(expectedUser, nil).
			Times(1)

		result, err := service.SetIsActive(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "u1", result.UserID)
		assert.Equal(t, "Alice", result.Username)
		assert.True(t, result.IsActive)
	})

	t.Run("success - set user inactive", func(t *testing.T) {
		req := dto.SetUserIsActiveRequest{
			UserID:   "u2",
			IsActive: false,
		}

		expectedUser := &domain.User{
			UserID:   "u2",
			Username: "Bob",
			TeamName: "backend",
			IsActive: false,
		}

		mockUserRepo.EXPECT().
			SetUserIsActive(ctx, "u2", false).
			Return(expectedUser, nil).
			Times(1)

		result, err := service.SetIsActive(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "u2", result.UserID)
		assert.False(t, result.IsActive)
	})

	t.Run("error - user not found", func(t *testing.T) {
		req := dto.SetUserIsActiveRequest{
			UserID:   "nonexistent",
			IsActive: false,
		}

		mockUserRepo.EXPECT().
			SetUserIsActive(ctx, "nonexistent", false).
			Return(nil, domain.ErrUserNotFound).
			Times(1)

		result, err := service.SetIsActive(ctx, req)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, domain.ErrUserNotFound)
	})

	t.Run("error - empty user ID", func(t *testing.T) {
		req := dto.SetUserIsActiveRequest{
			UserID:   "",
			IsActive: true,
		}

		result, err := service.SetIsActive(ctx, req)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "user_id is required")
	})

	t.Run("error - database error", func(t *testing.T) {
		req := dto.SetUserIsActiveRequest{
			UserID:   "u1",
			IsActive: true,
		}

		dbErr := errors.New("database connection lost")
		mockUserRepo.EXPECT().
			SetUserIsActive(ctx, "u1", true).
			Return(nil, dbErr).
			Times(1)

		result, err := service.SetIsActive(ctx, req)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "database connection lost")
	})
}
