package http

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/delivery/http/dto"
	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/domain"
)

func mapDomainError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, domain.ErrTeamAlreadyExists):
		return c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
			dto.ErrCodeTeamExists,
			"team_name already exists",
		))

	case errors.Is(err, domain.ErrTeamNotFound):
		return c.JSON(http.StatusNotFound, dto.NewErrorResponse(
			dto.ErrCodeNotFound,
			"team not found",
		))

	case errors.Is(err, domain.ErrUserNotFound):
		return c.JSON(http.StatusNotFound, dto.NewErrorResponse(
			dto.ErrCodeNotFound,
			"user not found",
		))

	case errors.Is(err, domain.ErrPRAlreadyExists):
		return c.JSON(http.StatusConflict, dto.NewErrorResponse(
			dto.ErrCodePRExists,
			"PR id already exists",
		))

	case errors.Is(err, domain.ErrPRNotFound):
		return c.JSON(http.StatusNotFound, dto.NewErrorResponse(
			dto.ErrCodeNotFound,
			"pull request not found",
		))

	case errors.Is(err, domain.ErrPRMerged):
		return c.JSON(http.StatusConflict, dto.NewErrorResponse(
			dto.ErrCodePRMerged,
			"cannot reassign on merged PR",
		))

	case errors.Is(err, domain.ErrReviewerNotAssigned):
		return c.JSON(http.StatusConflict, dto.NewErrorResponse(
			dto.ErrCodeNotAssigned,
			"reviewer is not assigned to this PR",
		))

	case errors.Is(err, domain.ErrNoCandidates):
		return c.JSON(http.StatusConflict, dto.NewErrorResponse(
			dto.ErrCodeNoCandidate,
			"no active replacement candidate in team",
		))

	default:
		return c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(
			"INTERNAL_ERROR",
			err.Error(),
		))
	}
}
