package http

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/delivery/http/dto"
	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/usecase"
)

func (h *Handler) CreatePR(c echo.Context) error {
	var req dto.CreatePRRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
			dto.ErrCodeInvalidInput,
			"invalid JSON: "+err.Error(),
		))
	}

	if req.PullRequestID == "" {
		return c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
			dto.ErrCodeInvalidInput,
			"pull_request_id is required",
		))
	}
	if req.PullRequestName == "" {
		return c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
			dto.ErrCodeInvalidInput,
			"pull_request_name is required",
		))
	}
	if req.AuthorID == "" {
		return c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
			dto.ErrCodeInvalidInput,
			"author_id is required",
		))
	}

	usecaseReq := usecase.CreatePRRequest{
		PullRequestID:   req.PullRequestID,
		PullRequestName: req.PullRequestName,
		AuthorID:        req.AuthorID,
	}

	pr, err := h.prUC.CreatePR(c.Request().Context(), usecaseReq)
	if err != nil {
		return mapDomainError(c, err)
	}

	response := dto.ToPRResponse(pr)
	return c.JSON(http.StatusCreated, response)
}

func (h *Handler) MergePR(c echo.Context) error {
	var req dto.MergePRRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
			dto.ErrCodeInvalidInput,
			"invalid JSON: "+err.Error(),
		))
	}

	if req.PullRequestID == "" {
		return c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
			dto.ErrCodeInvalidInput,
			"pull_request_id is required",
		))
	}

	usecaseReq := usecase.MergePRRequest{
		PullRequestID: req.PullRequestID,
	}

	pr, err := h.prUC.MergePR(c.Request().Context(), usecaseReq)
	if err != nil {
		return mapDomainError(c, err)
	}

	response := dto.ToPRResponse(pr)
	return c.JSON(http.StatusOK, response)
}

func (h *Handler) ReassignReviewer(c echo.Context) error {
	var req dto.ReassignReviewerRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
			dto.ErrCodeInvalidInput,
			"invalid JSON: "+err.Error(),
		))
	}

	if req.PullRequestID == "" {
		return c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
			dto.ErrCodeInvalidInput,
			"pull_request_id is required",
		))
	}
	if req.OldReviewerID == "" {
		return c.JSON(http.StatusBadRequest, dto.NewErrorResponse(
			dto.ErrCodeInvalidInput,
			"old_reviewer_id is required",
		))
	}

	usecaseReq := usecase.ReassignReviewerRequest{
		PullRequestID: req.PullRequestID,
		OldReviewerID: req.OldReviewerID,
	}

	result, err := h.prUC.ReassignReviewer(c.Request().Context(), usecaseReq)
	if err != nil {
		return mapDomainError(c, err)
	}

	response := dto.ToReassignResponse(result.PullRequest, result.ReplacedBy)
	return c.JSON(http.StatusOK, response)
}
