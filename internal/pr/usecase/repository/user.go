package repository

import (
	"context"

	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/domain"
)

//go:generate mockgen -destination=../mocks/mock_user_repository.go -package=mocks github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/usecase/repository UserRepository
type UserRepository interface {
	UpsertUser(ctx context.Context, user *domain.User) error
	GetUser(ctx context.Context, userID string) (*domain.User, error)
	GetUsersByTeam(ctx context.Context, teamName string) ([]domain.User, error)
	SetUserIsActive(ctx context.Context, userID string, isActive bool) (*domain.User, error)
	UserExists(ctx context.Context, userID string) (bool, error)
}
