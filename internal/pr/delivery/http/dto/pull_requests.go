package dto

import (
	"time"

	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/domain"
)

type CreatePRRequest struct {
	PullRequestID   string `json:"pull_request_id" validate:"required"`
	PullRequestName string `json:"pull_request_name" validate:"required"`
	AuthorID        string `json:"author_id" validate:"required"`
}

type MergePRRequest struct {
	PullRequestID string `json:"pull_request_id" validate:"required"`
}

type ReassignReviewerRequest struct {
	PullRequestID string `json:"pull_request_id" validate:"required"`
	OldReviewerID string `json:"old_reviewer_id" validate:"required"`
}

type PRResponse struct {
	PR PullRequest `json:"pr"`
}

type PullRequest struct {
	PullRequestID     string   `json:"pull_request_id"`
	PullRequestName   string   `json:"pull_request_name"`
	AuthorID          string   `json:"author_id"`
	Status            string   `json:"status"`
	AssignedReviewers []string `json:"assigned_reviewers"`
	CreatedAt         *string  `json:"createdAt,omitempty"`
	MergedAt          *string  `json:"mergedAt,omitempty"`
}

type ReassignReviewerResponse struct {
	PR         PullRequest `json:"pr"`
	ReplacedBy string      `json:"replaced_by"`
}

type GetReviewerPRsResponse struct {
	UserID       string             `json:"user_id"`
	PullRequests []PullRequestShort `json:"pull_requests"`
}

type PullRequestShort struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
	Status          string `json:"status"`
}

func ToPRResponse(pr *domain.PullRequest) PRResponse {
	return PRResponse{
		PR: toPullRequest(pr),
	}
}

func toPullRequest(pr *domain.PullRequest) PullRequest {
	var createdAt, mergedAt *string

	if !pr.CreatedAt.IsZero() {
		t := pr.CreatedAt.Format(time.RFC3339)
		createdAt = &t
	}

	if pr.MergedAt != nil && !pr.MergedAt.IsZero() {
		t := pr.MergedAt.Format(time.RFC3339)
		mergedAt = &t
	}

	return PullRequest{
		PullRequestID:     pr.PullRequestID,
		PullRequestName:   pr.PullRequestName,
		AuthorID:          pr.AuthorID,
		Status:            string(pr.Status),
		AssignedReviewers: pr.AssignedReviewers,
		CreatedAt:         createdAt,
		MergedAt:          mergedAt,
	}
}

func ToReassignResponse(pr *domain.PullRequest, replacedBy string) ReassignReviewerResponse {
	return ReassignReviewerResponse{
		PR:         toPullRequest(pr),
		ReplacedBy: replacedBy,
	}
}

func ToGetReviewerPRsResponse(userID string, prs []domain.PullRequestShort) GetReviewerPRsResponse {
	prList := make([]PullRequestShort, len(prs))
	for i, pr := range prs {
		prList[i] = PullRequestShort{
			PullRequestID:   pr.PullRequestID,
			PullRequestName: pr.PullRequestName,
			AuthorID:        pr.AuthorID,
			Status:          string(pr.Status),
		}
	}

	return GetReviewerPRsResponse{
		UserID:       userID,
		PullRequests: prList,
	}
}
