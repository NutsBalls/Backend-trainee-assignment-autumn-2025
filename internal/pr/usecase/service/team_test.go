package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/domain"
	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/usecase"
	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/usecase/mocks"
)

func TestTeamService_CreateTeam(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUOW := mocks.NewMockUnitOfWork(ctrl)
	mockTeamRepo := mocks.NewMockTeamRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	mockUOW.EXPECT().Teams().Return(mockTeamRepo).AnyTimes()
	mockUOW.EXPECT().Users().Return(mockUserRepo).AnyTimes()

	service := NewTeamService(mockUOW)
	ctx := context.Background()

	t.Run("success - create team with members", func(t *testing.T) {
		req := usecase.CreateTeamRequest{
			TeamName: "backend",
			Members: []usecase.CreateTeamMember{
				{UserID: "u1", Username: "Alice", IsActive: true},
				{UserID: "u2", Username: "Bob", IsActive: true},
			},
		}

		mockTeamRepo.EXPECT().
			TeamExists(ctx, "backend").
			Return(false, nil).
			Times(1)

		mockUOW.EXPECT().
			WithinTransaction(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			}).
			Times(1)

		mockTeamRepo.EXPECT().
			CreateTeam(ctx, "backend").
			Return(nil).
			Times(1)

		mockUserRepo.EXPECT().
			UpsertUser(ctx, &domain.User{
				UserID:   "u1",
				Username: "Alice",
				TeamName: "backend",
				IsActive: true,
			}).
			Return(nil).
			Times(1)

		mockUserRepo.EXPECT().
			UpsertUser(ctx, &domain.User{
				UserID:   "u2",
				Username: "Bob",
				TeamName: "backend",
				IsActive: true,
			}).
			Return(nil).
			Times(1)

		mockTeamRepo.EXPECT().
			GetTeam(ctx, "backend").
			Return(&domain.Team{
				TeamName: "backend",
				Members: []domain.User{
					{UserID: "u1", Username: "Alice", TeamName: "backend", IsActive: true},
					{UserID: "u2", Username: "Bob", TeamName: "backend", IsActive: true},
				},
			}, nil).
			Times(1)

		result, err := service.CreateTeam(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "backend", result.TeamName)
		assert.Len(t, result.Members, 2)
		assert.Equal(t, "Alice", result.Members[0].Username)
		assert.Equal(t, "Bob", result.Members[1].Username)
	})

	t.Run("error - team already exists", func(t *testing.T) {
		req := usecase.CreateTeamRequest{
			TeamName: "backend",
			Members: []usecase.CreateTeamMember{
				{UserID: "u1", Username: "Alice", IsActive: true},
			},
		}

		mockTeamRepo.EXPECT().
			TeamExists(ctx, "backend").
			Return(true, nil).
			Times(1)

		result, err := service.CreateTeam(ctx, req)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, domain.ErrTeamAlreadyExists)
	})

	t.Run("error - empty team name", func(t *testing.T) {
		req := usecase.CreateTeamRequest{
			TeamName: "",
			Members: []usecase.CreateTeamMember{
				{UserID: "u1", Username: "Alice", IsActive: true},
			},
		}

		result, err := service.CreateTeam(ctx, req)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "team_name is required")
	})

	t.Run("error - no members", func(t *testing.T) {
		req := usecase.CreateTeamRequest{
			TeamName: "backend",
			Members:  []usecase.CreateTeamMember{},
		}

		result, err := service.CreateTeam(ctx, req)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "members are required")
	})

	t.Run("error - database error on TeamExists", func(t *testing.T) {
		req := usecase.CreateTeamRequest{
			TeamName: "backend",
			Members: []usecase.CreateTeamMember{
				{UserID: "u1", Username: "Alice", IsActive: true},
			},
		}

		dbErr := errors.New("database connection failed")
		mockTeamRepo.EXPECT().
			TeamExists(ctx, "backend").
			Return(false, dbErr).
			Times(1)

		result, err := service.CreateTeam(ctx, req)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "database connection failed")
	})

	t.Run("error - CreateTeam fails in transaction", func(t *testing.T) {
		req := usecase.CreateTeamRequest{
			TeamName: "backend",
			Members: []usecase.CreateTeamMember{
				{UserID: "u1", Username: "Alice", IsActive: true},
			},
		}

		mockTeamRepo.EXPECT().TeamExists(ctx, "backend").Return(false, nil)

		createErr := errors.New("create team failed")
		mockUOW.EXPECT().
			WithinTransaction(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			})

		mockTeamRepo.EXPECT().
			CreateTeam(ctx, "backend").
			Return(createErr)

		result, err := service.CreateTeam(ctx, req)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "create team failed")
	})

	t.Run("error - UpsertUser fails in transaction", func(t *testing.T) {
		req := usecase.CreateTeamRequest{
			TeamName: "backend",
			Members: []usecase.CreateTeamMember{
				{UserID: "u1", Username: "Alice", IsActive: true},
			},
		}

		mockTeamRepo.EXPECT().TeamExists(ctx, "backend").Return(false, nil)

		upsertErr := errors.New("upsert user failed")
		mockUOW.EXPECT().
			WithinTransaction(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			})

		mockTeamRepo.EXPECT().CreateTeam(ctx, "backend").Return(nil)
		mockUserRepo.EXPECT().
			UpsertUser(ctx, gomock.Any()).
			Return(upsertErr)

		result, err := service.CreateTeam(ctx, req)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "upsert user")
	})
}

func TestTeamService_GetTeam(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUOW := mocks.NewMockUnitOfWork(ctrl)
	mockTeamRepo := mocks.NewMockTeamRepository(ctrl)

	mockUOW.EXPECT().Teams().Return(mockTeamRepo).AnyTimes()

	service := NewTeamService(mockUOW)
	ctx := context.Background()

	t.Run("success - get existing team", func(t *testing.T) {
		expectedTeam := &domain.Team{
			TeamName: "backend",
			Members: []domain.User{
				{UserID: "u1", Username: "Alice", TeamName: "backend", IsActive: true},
				{UserID: "u2", Username: "Bob", TeamName: "backend", IsActive: true},
			},
		}

		mockTeamRepo.EXPECT().
			GetTeam(ctx, "backend").
			Return(expectedTeam, nil).
			Times(1)

		result, err := service.GetTeam(ctx, "backend")

		require.NoError(t, err)
		assert.Equal(t, expectedTeam, result)
		assert.Equal(t, "backend", result.TeamName)
		assert.Len(t, result.Members, 2)
	})

	t.Run("error - team not found", func(t *testing.T) {
		mockTeamRepo.EXPECT().
			GetTeam(ctx, "nonexistent").
			Return(nil, domain.ErrTeamNotFound).
			Times(1)

		result, err := service.GetTeam(ctx, "nonexistent")

		require.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, domain.ErrTeamNotFound)
	})

	t.Run("error - empty team name", func(t *testing.T) {
		result, err := service.GetTeam(ctx, "")

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "team_name is required")
	})

	t.Run("error - database error", func(t *testing.T) {
		dbErr := errors.New("database error")
		mockTeamRepo.EXPECT().
			GetTeam(ctx, "backend").
			Return(nil, dbErr).
			Times(1)

		result, err := service.GetTeam(ctx, "backend")

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "database error")
	})
}
