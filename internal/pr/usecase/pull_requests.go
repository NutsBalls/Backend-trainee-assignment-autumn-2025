package usecase

import "github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/domain"

type CreatePRRequest struct {
	PullRequestID   string
	PullRequestName string
	AuthorID        string
}

type MergePRRequest struct {
	PullRequestID string
}

type ReassignReviewerRequest struct {
	PullRequestID string
	OldReviewerID string
}

type ReassignReviewerResponse struct {
	PullRequest *domain.PullRequest
	ReplacedBy  string
}
