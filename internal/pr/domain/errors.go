package domain

import "errors"

var (
	ErrTeamAlreadyExists = errors.New("team already exists")
	ErrTeamNotFound      = errors.New("team not found")

	ErrUserNotFound = errors.New("user not found")

	ErrPRAlreadyExists     = errors.New("pull request already exists")
	ErrPRNotFound          = errors.New("pull request not found")
	ErrPRMerged            = errors.New("cannot modify merged pull request")
	ErrReviewerNotAssigned = errors.New("reviewer is not assigned to this PR")
	ErrNoCandidates        = errors.New("no active replacement candidate in team")
)
