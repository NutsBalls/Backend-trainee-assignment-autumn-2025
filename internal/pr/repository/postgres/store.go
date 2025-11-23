package postgres

import (
	"context"
	"fmt"
	"log"

	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/repository/postgres/sqlc"
	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/usecase/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	pool         *pgxpool.Pool
	teamRepo     *TeamRepository
	userRepo     *UserRepository
	prRepo       *PRRepository
	reviewerRepo *ReviewerRepository
}

func NewStore(pool *pgxpool.Pool) *Store {
	queries := sqlc.New(pool)

	return &Store{
		pool:         pool,
		teamRepo:     NewTeamRepository(queries),
		userRepo:     NewUserRepository(queries),
		prRepo:       NewPRRepository(queries),
		reviewerRepo: NewReviewerRepository(queries),
	}
}

func (s *Store) Teams() repository.TeamRepository {
	return s.teamRepo
}

func (s *Store) Users() repository.UserRepository {
	return s.userRepo
}

func (s *Store) PullRequests() repository.PRRepository {
	return s.prRepo
}

func (s *Store) Reviewers() repository.ReviewerRepository {
	return s.reviewerRepo
}

func (s *Store) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
			log.Printf("failed to rollback transaction: %v", err)
		}
	}()

	txCtx := injectTx(ctx, tx)

	if err := fn(txCtx); err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

type txKey struct{}

func injectTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}
