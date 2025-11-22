package dto

import "github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/domain"

type CreateTeamRequest struct {
	TeamName string       `json:"team_name" validate:"required"`
	Members  []TeamMember `json:"members" validate:"required,min=1"`
}

type TeamMember struct {
	UserID   string `json:"user_id" validate:"required"`
	Username string `json:"username" validate:"required"`
	IsActive bool   `json:"is_active"`
}

type TeamResponse struct {
	Team Team `json:"team"`
}

type Team struct {
	TeamName string       `json:"team_name"`
	Members  []TeamMember `json:"members"`
}

func ToTeamResponse(team *domain.Team) TeamResponse {
	members := make([]TeamMember, len(team.Members))
	for i, m := range team.Members {
		members[i] = TeamMember{
			UserID:   m.UserID,
			Username: m.Username,
			IsActive: m.IsActive,
		}
	}

	return TeamResponse{
		Team: Team{
			TeamName: team.TeamName,
			Members:  members,
		},
	}
}
