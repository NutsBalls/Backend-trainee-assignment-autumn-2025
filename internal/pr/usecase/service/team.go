package service

import (
	"context"
	"fmt"

	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/domain"
	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/usecase/dto"
	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/usecase/repository"
)

type TeamService struct {
	uow repository.UnitOfWork
}

func NewTeamService(uow repository.UnitOfWork) *TeamService {
	return &TeamService{uow: uow}
}

func (s *TeamService) CreateTeam(ctx context.Context, req dto.CreateTeamRequest) (*domain.Team, error) {
	if req.TeamName == "" {
		return nil, fmt.Errorf("team_name is required")
	}
	if len(req.Members) == 0 {
		return nil, fmt.Errorf("members are required")
	}

	exists, err := s.uow.Teams().TeamExists(ctx, req.TeamName)
	if err != nil {
		return nil, fmt.Errorf("check team exists: %w", err)
	}
	if exists {
		return nil, domain.ErrTeamAlreadyExists
	}

	err = s.uow.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := s.uow.Teams().CreateTeam(txCtx, req.TeamName); err != nil {
			return fmt.Errorf("create team: %w", err)
		}

		for _, member := range req.Members {
			user := &domain.User{
				UserID:   member.UserID,
				Username: member.Username,
				TeamName: req.TeamName,
				IsActive: member.IsActive,
			}
			if err := s.uow.Users().UpsertUser(txCtx, user); err != nil {
				return fmt.Errorf("upsert user %s: %w", member.UserID, err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	team, err := s.uow.Teams().GetTeam(ctx, req.TeamName)
	if err != nil {
		return nil, fmt.Errorf("get created team: %w", err)
	}

	return team, nil
}

func (s *TeamService) GetTeam(ctx context.Context, teamName string) (*domain.Team, error) {
	if teamName == "" {
		return nil, fmt.Errorf("team_name is required")
	}

	team, err := s.uow.Teams().GetTeam(ctx, teamName)
	if err != nil {
		return nil, err
	}

	return team, nil
}
