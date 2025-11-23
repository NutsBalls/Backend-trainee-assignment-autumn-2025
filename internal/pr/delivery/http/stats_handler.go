package http

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/delivery/http/dto"
)

func (h *Handler) GetUserStats(c echo.Context) error {
	ctx := c.Request().Context()

	stats, err := h.statsUC.GetUserAssignmentStats(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(
			"INTERNAL_ERROR",
			"failed to get user stats: "+err.Error(),
		))
	}

	out := make([]dto.UserAssignmentStats, len(stats))
	for i, s := range stats {
		out[i] = dto.UserAssignmentStats{
			UserID:           s.UserID,
			Username:         s.Username,
			TeamName:         s.TeamName,
			AssignmentsCount: s.AssignmentsCount,
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"users": out,
	})
}

func (h *Handler) GetPRStats(c echo.Context) error {
	ctx := c.Request().Context()

	stats, err := h.statsUC.GetPRStats(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(
			"INTERNAL_ERROR",
			"failed to get PR stats: "+err.Error(),
		))
	}

	out := dto.PRStats{
		TotalPRs:  stats.TotalPRs,
		OpenPRs:   stats.OpenPRs,
		MergedPRs: stats.MergedPRs,
	}

	return c.JSON(http.StatusOK, out)
}

func (h *Handler) GetReviewerWorkload(c echo.Context) error {
	ctx := c.Request().Context()

	workload, err := h.statsUC.GetReviewerWorkload(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(
			"INTERNAL_ERROR",
			"failed to get reviewer workload: "+err.Error(),
		))
	}

	out := make([]dto.ReviewerWorkload, len(workload))
	for i, w := range workload {
		out[i] = dto.ReviewerWorkload{
			UserID:       w.UserID,
			Username:     w.Username,
			TeamName:     w.TeamName,
			OpenPRsCount: w.OpenPRsCount,
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"reviewers": out,
	})
}
