package http

import (
	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/usecase"
)

type Handler struct {
	teamUC  usecase.TeamUseCase
	userUC  usecase.UserUseCase
	prUC    usecase.PRUseCase
	statsUC usecase.StatsUseCase
}

func NewHandler(teamUC usecase.TeamUseCase, userUC usecase.UserUseCase, prUC usecase.PRUseCase, statsUC usecase.StatsUseCase) *Handler {
	return &Handler{
		teamUC:  teamUC,
		userUC:  userUC,
		prUC:    prUC,
		statsUC: statsUC,
	}
}
