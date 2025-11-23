package http

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/delivery/http/dto"
	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/usecase"
)

func (h *Handler) SetUserIsActive(c echo.Context) error {
	var req dto.SetUserIsActiveRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
			dto.ErrCodeInvalidInput,
			"invalid JSON: "+err.Error(),
		))
	}

	if req.UserID == "" {
		return c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
			dto.ErrCodeInvalidInput,
			"user_id is required",
		))
	}

	usecaseReq := usecase.SetUserIsActiveRequest{
		UserID:   req.UserID,
		IsActive: req.IsActive,
	}

	user, err := h.userUC.SetIsActive(c.Request().Context(), usecaseReq)
	if err != nil {
		return mapDomainError(c, err)
	}

	response := dto.ToUserResponse(user)
	return c.JSON(http.StatusOK, response)
}

func (h *Handler) GetReviewerPRs(c echo.Context) error {
	userID := c.QueryParam("user_id")
	if userID == "" {
		return c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
			dto.ErrCodeInvalidInput,
			"user_id query parameter is required",
		))
	}

	prs, err := h.prUC.GetReviewerPRs(c.Request().Context(), userID)
	if err != nil {
		return mapDomainError(c, err)
	}

	response := dto.ToGetReviewerPRsResponse(userID, prs)
	return c.JSON(http.StatusOK, response)
}
