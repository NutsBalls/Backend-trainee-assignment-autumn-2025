package domain

import "time"

type Team struct {
	TeamName string
	Members  []User
}

type User struct {
	UserID   string
	Username string
	TeamName string
	IsActive bool
}

type PullRequest struct {
	PullRequestID     string
	PullRequestName   string
	AuthorID          string
	Status            PRStatus
	AssignedReviewers []string
	CreatedAt         *time.Time
	MergedAt          *time.Time
}

type PRStatus string

const (
	PRStatusOpen   PRStatus = "OPEN"
	PRStatusMerged PRStatus = "MERGED"
)

type PullRequestShort struct {
	PullRequestID   string
	PullRequestName string
	AuthorID        string
	Status          PRStatus
}

type UserAssignmentStats struct {
	UserID           string
	Username         string
	TeamName         string
	AssignmentsCount int64
}

type PRStats struct {
	TotalPRs  int64
	OpenPRs   int64
	MergedPRs int64
}

type ReviewerWorkload struct {
	UserID       string
	Username     string
	TeamName     string
	OpenPRsCount int64
}
