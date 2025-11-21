package repository

import "context"

type Transactor interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type UnitOfWork interface {
	Transactor

	Teams() TeamRepository
	Users() UserRepository
	PullRequests() PRRepository
	Reviewers() ReviewerRepository
}
