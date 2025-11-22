package repository

import "context"

type Transactor interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

//go:generate mockgen -destination=../mocks/mock_unit_of_work.go -package=mocks github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/usecase/repository UnitOfWork
type UnitOfWork interface {
	Transactor

	Teams() TeamRepository
	Users() UserRepository
	PullRequests() PRRepository
	Reviewers() ReviewerRepository
}
