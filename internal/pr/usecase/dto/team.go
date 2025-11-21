package dto

type CreateTeamRequest struct {
	TeamName string
	Members  []CreateTeamMember
}

type CreateTeamMember struct {
	UserID   string
	Username string
	IsActive bool
}
