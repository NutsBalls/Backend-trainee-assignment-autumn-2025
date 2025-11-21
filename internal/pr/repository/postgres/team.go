package postgres

import (
	"context"
	"fmt"

	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/domain"
	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/repository/postgres/sqlc"
)

type TeamRepository struct {
	queries *sqlc.Queries
}

func NewTeamRepository(queries *sqlc.Queries) *TeamRepository {
	return &TeamRepository{queries: queries}
}

func (r *TeamRepository) CreateTeam(ctx context.Context, teamName string) error {
	_, err := r.queries.CreateTeam(ctx, teamName)
	if err != nil {
		if isPgUniqueViolation(err) {
			return domain.ErrTeamAlreadyExists
		}
		return fmt.Errorf("create team: %w", err)
	}
	return nil
}

func (r *TeamRepository) TeamExists(ctx context.Context, teamName string) (bool, error) {
	exists, err := r.queries.TeamExists(ctx, teamName)
	if err != nil {
		return false, fmt.Errorf("check team exists: %w", err)
	}
	return exists, nil
}

func (r *TeamRepository) GetTeam(ctx context.Context, teamName string) (*domain.Team, error) {
	exists, err := r.TeamExists(ctx, teamName)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, domain.ErrTeamNotFound
	}

	users, err := r.queries.GetUsersByTeam(ctx, teamName)
	if err != nil {
		return nil, fmt.Errorf("get team members: %w", err)
	}

	members := make([]domain.User, len(users))
	for i, u := range users {
		members[i] = domain.User{
			UserID:   u.UserID,
			Username: u.Username,
			TeamName: u.TeamName,
			IsActive: u.IsActive,
		}
	}

	return &domain.Team{
		TeamName: teamName,
		Members:  members,
	}, nil
}
