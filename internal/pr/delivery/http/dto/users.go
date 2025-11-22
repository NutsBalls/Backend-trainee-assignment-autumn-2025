package dto

import "github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/domain"

type SetUserIsActiveRequest struct {
	UserID   string `json:"user_id" validate:"required"`
	IsActive bool   `json:"is_active"`
}

type UserResponse struct {
	User User `json:"user"`
}

type User struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

func ToUserResponse(user *domain.User) UserResponse {
	return UserResponse{
		User: User{
			UserID:   user.UserID,
			Username: user.Username,
			TeamName: user.TeamName,
			IsActive: user.IsActive,
		},
	}
}
