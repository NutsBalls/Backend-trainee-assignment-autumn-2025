package http

import (
	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/usecase"
)

type Handler struct {
	teamUC usecase.TeamUseCase
	userUC usecase.UserUseCase
	prUC   usecase.PRUseCase
}

func NewHandler(teamUC usecase.TeamUseCase, userUC usecase.UserUseCase, prUC usecase.PRUseCase) *Handler {
	return &Handler{
		teamUC: teamUC,
		userUC: userUC,
		prUC:   prUC,
	}
}
