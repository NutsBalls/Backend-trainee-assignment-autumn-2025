package http

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/delivery/http/dto"
	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/usecase"
)

func (h *Handler) CreateTeam(c echo.Context) error {
	var req dto.CreateTeamRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
			dto.ErrCodeInvalidInput,
			"invalid JSON: "+err.Error(),
		))
	}

	if req.TeamName == "" {
		return c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
			dto.ErrCodeInvalidInput,
			"team_name is required",
		))
	}
	if len(req.Members) == 0 {
		return c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
			dto.ErrCodeInvalidInput,
			"members are required",
		))
	}

	for i, member := range req.Members {
		if member.UserID == "" {
			return c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
				dto.ErrCodeInvalidInput,
				"member user_id is required at index "+string(rune(i)),
			))
		}
		if member.Username == "" {
			return c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
				dto.ErrCodeInvalidInput,
				"member username is required at index "+string(rune(i)),
			))
		}
	}

	usecaseReq := usecase.CreateTeamRequest{
		TeamName: req.TeamName,
		Members:  make([]usecase.CreateTeamMember, len(req.Members)),
	}
	for i, m := range req.Members {
		usecaseReq.Members[i] = usecase.CreateTeamMember{
			UserID:   m.UserID,
			Username: m.Username,
			IsActive: m.IsActive,
		}
	}

	team, err := h.teamUC.CreateTeam(c.Request().Context(), usecaseReq)
	if err != nil {
		return mapDomainError(c, err)
	}

	response := dto.ToTeamResponse(team)
	return c.JSON(http.StatusCreated, response)
}

func (h *Handler) GetTeam(c echo.Context) error {
	teamName := c.QueryParam("team_name")
	if teamName == "" {
		return c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
			dto.ErrCodeInvalidInput,
			"team_name query parameter is required",
		))
	}

	team, err := h.teamUC.GetTeam(c.Request().Context(), teamName)
	if err != nil {
		return mapDomainError(c, err)
	}

	response := dto.ToTeamResponse(team)
	return c.JSON(http.StatusOK, response)
}
