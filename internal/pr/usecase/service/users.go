package service

import (
	"context"
	"fmt"

	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/domain"
	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/usecase/dto"
	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/usecase/repository"
)

type UserService struct {
	uow repository.UnitOfWork
}

func NewUserService(uow repository.UnitOfWork) *UserService {
	return &UserService{uow: uow}
}

func (s *UserService) SetIsActive(ctx context.Context, req dto.SetUserIsActiveRequest) (*domain.User, error) {
	if req.UserID == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	user, err := s.uow.Users().SetUserIsActive(ctx, req.UserID, req.IsActive)
	if err != nil {
		return nil, err
	}

	return user, nil
}
