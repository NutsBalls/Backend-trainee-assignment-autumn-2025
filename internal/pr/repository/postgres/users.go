package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/domain"
	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/repository/postgres/sqlc"
	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	queries *sqlc.Queries
}

func NewUserRepository(queries *sqlc.Queries) *UserRepository {
	return &UserRepository{queries: queries}
}

func (r *UserRepository) UpsertUser(ctx context.Context, user *domain.User) error {
	exists, err := r.queries.UserExists(ctx, user.UserID)
	if err != nil {
		return err
	}

	if exists {
		err = r.queries.UpdateUser(ctx, sqlc.UpdateUserParams{
			UserID:   user.UserID,
			Username: user.Username,
			TeamName: user.TeamName,
			IsActive: user.IsActive,
		})
	} else {
		_, err = r.queries.InsertUser(ctx, sqlc.InsertUserParams{
			UserID:   user.UserID,
			Username: user.Username,
			TeamName: user.TeamName,
			IsActive: user.IsActive,
		})
	}

	if err != nil {
		return fmt.Errorf("upsert user: %w", err)
	}
	return nil
}

func (r *UserRepository) GetUser(ctx context.Context, userID string) (*domain.User, error) {
	user, err := r.queries.GetUser(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("get user: %w", err)
	}

	return &domain.User{
		UserID:   user.UserID,
		Username: user.Username,
		TeamName: user.TeamName,
		IsActive: user.IsActive,
	}, nil
}

func (r *UserRepository) GetUsersByTeam(ctx context.Context, teamName string) ([]domain.User, error) {
	users, err := r.queries.GetUsersByTeam(ctx, teamName)
	if err != nil {
		return nil, fmt.Errorf("get users by team: %w", err)
	}

	result := make([]domain.User, len(users))
	for i, u := range users {
		result[i] = domain.User{
			UserID:   u.UserID,
			Username: u.Username,
			TeamName: u.TeamName,
			IsActive: u.IsActive,
		}
	}
	return result, nil
}

func (r *UserRepository) SetUserIsActive(ctx context.Context, userID string, isActive bool) (*domain.User, error) {
	user, err := r.queries.SetUserActivity(ctx, sqlc.SetUserActivityParams{
		UserID:   userID,
		IsActive: isActive,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("set user activity: %w", err)
	}

	return &domain.User{
		UserID:   user.UserID,
		Username: user.Username,
		TeamName: user.TeamName,
		IsActive: user.IsActive,
	}, nil
}

func (r *UserRepository) UserExists(ctx context.Context, userID string) (bool, error) {
	exists, err := r.queries.UserExists(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("check user exists: %w", err)
	}
	return exists, nil
}
