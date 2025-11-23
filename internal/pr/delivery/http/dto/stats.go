package dto

type UserAssignmentStats struct {
	UserID           string `json:"user_id"`
	Username         string `json:"username"`
	TeamName         string `json:"team_name"`
	AssignmentsCount int64  `json:"assignments_count"`
}

type PRStats struct {
	TotalPRs  int64 `json:"total_prs"`
	OpenPRs   int64 `json:"open_prs"`
	MergedPRs int64 `json:"merged_prs"`
}

type ReviewerWorkload struct {
	UserID       string `json:"user_id"`
	Username     string `json:"username"`
	TeamName     string `json:"team_name"`
	OpenPRsCount int64  `json:"open_prs_count"`
}
