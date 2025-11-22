package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/domain"
	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/usecase"
	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/usecase/mocks"
)

func TestPRService_CreatePR(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUOW := mocks.NewMockUnitOfWork(ctrl)
	mockPRRepo := mocks.NewMockPRRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockReviewerRepo := mocks.NewMockReviewerRepository(ctrl)

	mockUOW.EXPECT().PullRequests().Return(mockPRRepo).AnyTimes()
	mockUOW.EXPECT().Users().Return(mockUserRepo).AnyTimes()
	mockUOW.EXPECT().Reviewers().Return(mockReviewerRepo).AnyTimes()

	service := NewPRService(mockUOW)
	ctx := context.Background()

	t.Run("success - create PR with 2 reviewers", func(t *testing.T) {
		req := usecase.CreatePRRequest{
			PullRequestID:   "pr-1001",
			PullRequestName: "Add authentication",
			AuthorID:        "u1",
		}

		author := &domain.User{
			UserID:   "u1",
			Username: "Alice",
			TeamName: "backend",
			IsActive: true,
		}

		candidates := []string{"u2", "u3"}

		now := time.Now()
		expectedPR := &domain.PullRequest{
			PullRequestID:     "pr-1001",
			PullRequestName:   "Add authentication",
			AuthorID:          "u1",
			Status:            domain.PRStatusOpen,
			AssignedReviewers: []string{"u2", "u3"},
			CreatedAt:         &now,
		}

		mockPRRepo.EXPECT().PRExists(ctx, "pr-1001").Return(false, nil)
		mockUserRepo.EXPECT().GetUser(ctx, "u1").Return(author, nil)

		mockUOW.EXPECT().WithinTransaction(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			})

		mockPRRepo.EXPECT().CreatePR(ctx, gomock.Any()).Return(nil)
		mockReviewerRepo.EXPECT().
			FindCandidatesForNewPR(ctx, "backend", "u1").
			Return(candidates, nil)
		mockReviewerRepo.EXPECT().AssignReviewer(ctx, "pr-1001", "u2").Return(nil)
		mockReviewerRepo.EXPECT().AssignReviewer(ctx, "pr-1001", "u3").Return(nil)
		mockPRRepo.EXPECT().GetPRWithReviewers(ctx, "pr-1001").Return(expectedPR, nil)

		result, err := service.CreatePR(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "pr-1001", result.PullRequestID)
		assert.Equal(t, "Add authentication", result.PullRequestName)
		assert.Equal(t, domain.PRStatusOpen, result.Status)
		assert.Len(t, result.AssignedReviewers, 2)
		assert.Contains(t, result.AssignedReviewers, "u2")
		assert.Contains(t, result.AssignedReviewers, "u3")
	})

	t.Run("success - create PR with 1 reviewer", func(t *testing.T) {
		req := usecase.CreatePRRequest{
			PullRequestID:   "pr-1002",
			PullRequestName: "Fix bug",
			AuthorID:        "u1",
		}

		author := &domain.User{UserID: "u1", TeamName: "backend", IsActive: true}
		candidates := []string{"u2"}

		now := time.Now()
		expectedPR := &domain.PullRequest{
			PullRequestID:     "pr-1002",
			PullRequestName:   "Fix bug",
			AuthorID:          "u1",
			Status:            domain.PRStatusOpen,
			AssignedReviewers: []string{"u2"},
			CreatedAt:         &now,
		}

		mockPRRepo.EXPECT().PRExists(ctx, "pr-1002").Return(false, nil)
		mockUserRepo.EXPECT().GetUser(ctx, "u1").Return(author, nil)
		mockUOW.EXPECT().WithinTransaction(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			})
		mockPRRepo.EXPECT().CreatePR(ctx, gomock.Any()).Return(nil)
		mockReviewerRepo.EXPECT().
			FindCandidatesForNewPR(ctx, "backend", "u1").
			Return(candidates, nil)
		mockReviewerRepo.EXPECT().AssignReviewer(ctx, "pr-1002", "u2").Return(nil)
		mockPRRepo.EXPECT().GetPRWithReviewers(ctx, "pr-1002").Return(expectedPR, nil)

		result, err := service.CreatePR(ctx, req)

		require.NoError(t, err)
		assert.Len(t, result.AssignedReviewers, 1)
	})

	t.Run("success - create PR with 0 reviewers", func(t *testing.T) {
		req := usecase.CreatePRRequest{
			PullRequestID:   "pr-1003",
			PullRequestName: "Update docs",
			AuthorID:        "u1",
		}

		author := &domain.User{UserID: "u1", TeamName: "backend", IsActive: true}
		candidates := []string{}

		now := time.Now()
		expectedPR := &domain.PullRequest{
			PullRequestID:     "pr-1003",
			PullRequestName:   "Update docs",
			AuthorID:          "u1",
			Status:            domain.PRStatusOpen,
			AssignedReviewers: []string{},
			CreatedAt:         &now,
		}

		mockPRRepo.EXPECT().PRExists(ctx, "pr-1003").Return(false, nil)
		mockUserRepo.EXPECT().GetUser(ctx, "u1").Return(author, nil)
		mockUOW.EXPECT().WithinTransaction(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			})
		mockPRRepo.EXPECT().CreatePR(ctx, gomock.Any()).Return(nil)
		mockReviewerRepo.EXPECT().
			FindCandidatesForNewPR(ctx, "backend", "u1").
			Return(candidates, nil)
		mockPRRepo.EXPECT().GetPRWithReviewers(ctx, "pr-1003").Return(expectedPR, nil)

		result, err := service.CreatePR(ctx, req)

		require.NoError(t, err)
		assert.Len(t, result.AssignedReviewers, 0)
	})

	t.Run("error - PR already exists", func(t *testing.T) {
		req := usecase.CreatePRRequest{
			PullRequestID:   "pr-1001",
			PullRequestName: "Add auth",
			AuthorID:        "u1",
		}

		mockPRRepo.EXPECT().PRExists(ctx, "pr-1001").Return(true, nil)

		result, err := service.CreatePR(ctx, req)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, domain.ErrPRAlreadyExists)
	})

	t.Run("error - author not found", func(t *testing.T) {
		req := usecase.CreatePRRequest{
			PullRequestID:   "pr-1001",
			PullRequestName: "Add auth",
			AuthorID:        "nonexistent",
		}

		mockPRRepo.EXPECT().PRExists(ctx, "pr-1001").Return(false, nil)
		mockUserRepo.EXPECT().GetUser(ctx, "nonexistent").Return(nil, domain.ErrUserNotFound)

		result, err := service.CreatePR(ctx, req)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, domain.ErrUserNotFound)
	})

	t.Run("error - empty pull request ID", func(t *testing.T) {
		req := usecase.CreatePRRequest{
			PullRequestID:   "",
			PullRequestName: "Add auth",
			AuthorID:        "u1",
		}

		result, err := service.CreatePR(ctx, req)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "pull_request_id is required")
	})

	t.Run("error - empty pull request name", func(t *testing.T) {
		req := usecase.CreatePRRequest{
			PullRequestID:   "pr-1001",
			PullRequestName: "",
			AuthorID:        "u1",
		}

		result, err := service.CreatePR(ctx, req)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "pull_request_name is required")
	})

	t.Run("error - empty author ID", func(t *testing.T) {
		req := usecase.CreatePRRequest{
			PullRequestID:   "pr-1001",
			PullRequestName: "Add auth",
			AuthorID:        "",
		}

		result, err := service.CreatePR(ctx, req)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "author_id is required")
	})
}

func TestPRService_MergePR(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUOW := mocks.NewMockUnitOfWork(ctrl)
	mockPRRepo := mocks.NewMockPRRepository(ctrl)

	mockUOW.EXPECT().PullRequests().Return(mockPRRepo).AnyTimes()

	service := NewPRService(mockUOW)
	ctx := context.Background()

	t.Run("success - merge PR", func(t *testing.T) {
		req := usecase.MergePRRequest{
			PullRequestID: "pr-1001",
		}

		now := time.Now()
		expectedPR := &domain.PullRequest{
			PullRequestID:     "pr-1001",
			PullRequestName:   "Add auth",
			AuthorID:          "u1",
			Status:            domain.PRStatusMerged,
			AssignedReviewers: []string{"u2", "u3"},
			CreatedAt:         &now,
			MergedAt:          &now,
		}

		mockPRRepo.EXPECT().
			MergePR(ctx, "pr-1001").
			Return(expectedPR, nil).
			Times(1)

		result, err := service.MergePR(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "pr-1001", result.PullRequestID)
		assert.Equal(t, domain.PRStatusMerged, result.Status)
		assert.NotNil(t, result.MergedAt)
	})

	t.Run("success - idempotent merge (already merged)", func(t *testing.T) {
		req := usecase.MergePRRequest{
			PullRequestID: "pr-1001",
		}

		now := time.Now()
		alreadyMergedPR := &domain.PullRequest{
			PullRequestID: "pr-1001",
			Status:        domain.PRStatusMerged,
			MergedAt:      &now,
		}

		mockPRRepo.EXPECT().
			MergePR(ctx, "pr-1001").
			Return(alreadyMergedPR, nil).
			Times(1)

		result, err := service.MergePR(ctx, req)

		require.NoError(t, err)
		assert.Equal(t, domain.PRStatusMerged, result.Status)
	})

	t.Run("error - PR not found", func(t *testing.T) {
		req := usecase.MergePRRequest{
			PullRequestID: "nonexistent",
		}

		mockPRRepo.EXPECT().
			MergePR(ctx, "nonexistent").
			Return(nil, domain.ErrPRNotFound).
			Times(1)

		result, err := service.MergePR(ctx, req)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, domain.ErrPRNotFound)
	})

	t.Run("error - empty pull request ID", func(t *testing.T) {
		req := usecase.MergePRRequest{
			PullRequestID: "",
		}

		result, err := service.MergePR(ctx, req)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "pull_request_id is required")
	})
}

func TestPRService_ReassignReviewer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUOW := mocks.NewMockUnitOfWork(ctrl)
	mockPRRepo := mocks.NewMockPRRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockReviewerRepo := mocks.NewMockReviewerRepository(ctrl)

	mockUOW.EXPECT().PullRequests().Return(mockPRRepo).AnyTimes()
	mockUOW.EXPECT().Users().Return(mockUserRepo).AnyTimes()
	mockUOW.EXPECT().Reviewers().Return(mockReviewerRepo).AnyTimes()

	service := NewPRService(mockUOW)
	ctx := context.Background()

	t.Run("success - reassign reviewer", func(t *testing.T) {
		req := usecase.ReassignReviewerRequest{
			PullRequestID: "pr-1001",
			OldReviewerID: "u2",
		}

		now := time.Now()
		openPR := &domain.PullRequest{
			PullRequestID: "pr-1001",
			AuthorID:      "u1",
			Status:        domain.PRStatusOpen,
			CreatedAt:     &now,
		}

		oldReviewer := &domain.User{
			UserID:   "u2",
			TeamName: "backend",
		}

		candidates := []string{"u4"}

		updatedPR := &domain.PullRequest{
			PullRequestID:     "pr-1001",
			AuthorID:          "u1",
			Status:            domain.PRStatusOpen,
			AssignedReviewers: []string{"u3", "u4"},
			CreatedAt:         &now,
		}

		mockPRRepo.EXPECT().GetPR(ctx, "pr-1001").Return(openPR, nil)
		mockReviewerRepo.EXPECT().IsReviewerAssigned(ctx, "pr-1001", "u2").Return(true, nil)
		mockUserRepo.EXPECT().GetUser(ctx, "u2").Return(oldReviewer, nil)
		mockPRRepo.EXPECT().GetPRAuthorID(ctx, "pr-1001").Return("u1", nil)
		mockReviewerRepo.EXPECT().
			FindCandidatesForReassignment(ctx, "backend", "u1", "pr-1001").
			Return(candidates, nil)
		mockReviewerRepo.EXPECT().ReplaceReviewer(ctx, "pr-1001", "u2", "u4").Return(nil)
		mockPRRepo.EXPECT().GetPRWithReviewers(ctx, "pr-1001").Return(updatedPR, nil)

		result, err := service.ReassignReviewer(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "u4", result.ReplacedBy)
		assert.Contains(t, result.PullRequest.AssignedReviewers, "u4")
		assert.NotContains(t, result.PullRequest.AssignedReviewers, "u2")
	})

	t.Run("error - PR is merged", func(t *testing.T) {
		req := usecase.ReassignReviewerRequest{
			PullRequestID: "pr-1001",
			OldReviewerID: "u2",
		}

		now := time.Now()
		mergedPR := &domain.PullRequest{
			PullRequestID: "pr-1001",
			Status:        domain.PRStatusMerged,
			MergedAt:      &now,
		}

		mockPRRepo.EXPECT().GetPR(ctx, "pr-1001").Return(mergedPR, nil)

		result, err := service.ReassignReviewer(ctx, req)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, domain.ErrPRMerged)
	})

	t.Run("error - reviewer not assigned", func(t *testing.T) {
		req := usecase.ReassignReviewerRequest{
			PullRequestID: "pr-1001",
			OldReviewerID: "u5",
		}

		openPR := &domain.PullRequest{
			PullRequestID: "pr-1001",
			Status:        domain.PRStatusOpen,
		}

		mockPRRepo.EXPECT().GetPR(ctx, "pr-1001").Return(openPR, nil)
		mockReviewerRepo.EXPECT().IsReviewerAssigned(ctx, "pr-1001", "u5").Return(false, nil)

		result, err := service.ReassignReviewer(ctx, req)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, domain.ErrReviewerNotAssigned)
	})

	t.Run("error - no candidates available", func(t *testing.T) {
		req := usecase.ReassignReviewerRequest{
			PullRequestID: "pr-1001",
			OldReviewerID: "u2",
		}

		openPR := &domain.PullRequest{
			PullRequestID: "pr-1001",
			AuthorID:      "u1",
			Status:        domain.PRStatusOpen,
		}

		oldReviewer := &domain.User{
			UserID:   "u2",
			TeamName: "backend",
		}

		candidates := []string{}

		mockPRRepo.EXPECT().GetPR(ctx, "pr-1001").Return(openPR, nil)
		mockReviewerRepo.EXPECT().IsReviewerAssigned(ctx, "pr-1001", "u2").Return(true, nil)
		mockUserRepo.EXPECT().GetUser(ctx, "u2").Return(oldReviewer, nil)
		mockPRRepo.EXPECT().GetPRAuthorID(ctx, "pr-1001").Return("u1", nil)
		mockReviewerRepo.EXPECT().
			FindCandidatesForReassignment(ctx, "backend", "u1", "pr-1001").
			Return(candidates, nil)

		result, err := service.ReassignReviewer(ctx, req)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, domain.ErrNoCandidates)
	})

	t.Run("error - PR not found", func(t *testing.T) {
		req := usecase.ReassignReviewerRequest{
			PullRequestID: "nonexistent",
			OldReviewerID: "u2",
		}

		mockPRRepo.EXPECT().GetPR(ctx, "nonexistent").Return(nil, domain.ErrPRNotFound)

		result, err := service.ReassignReviewer(ctx, req)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, domain.ErrPRNotFound)
	})

	t.Run("error - empty pull request ID", func(t *testing.T) {
		req := usecase.ReassignReviewerRequest{
			PullRequestID: "",
			OldReviewerID: "u2",
		}

		result, err := service.ReassignReviewer(ctx, req)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "pull_request_id is required")
	})

	t.Run("error - empty old reviewer ID", func(t *testing.T) {
		req := usecase.ReassignReviewerRequest{
			PullRequestID: "pr-1001",
			OldReviewerID: "",
		}

		result, err := service.ReassignReviewer(ctx, req)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "old_reviewer_id is required")
	})
}

func TestPRService_GetReviewerPRs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUOW := mocks.NewMockUnitOfWork(ctrl)
	mockReviewerRepo := mocks.NewMockReviewerRepository(ctrl)

	mockUOW.EXPECT().Reviewers().Return(mockReviewerRepo).AnyTimes()

	service := NewPRService(mockUOW)
	ctx := context.Background()

	t.Run("success - get reviewer PRs", func(t *testing.T) {
		expectedPRs := []domain.PullRequestShort{
			{
				PullRequestID:   "pr-1001",
				PullRequestName: "Add auth",
				AuthorID:        "u1",
				Status:          domain.PRStatusOpen,
			},
			{
				PullRequestID:   "pr-1002",
				PullRequestName: "Fix bug",
				AuthorID:        "u3",
				Status:          domain.PRStatusMerged,
			},
		}

		mockReviewerRepo.EXPECT().
			ListPRsByReviewer(ctx, "u2").
			Return(expectedPRs, nil).
			Times(1)

		result, err := service.GetReviewerPRs(ctx, "u2")

		require.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "pr-1001", result[0].PullRequestID)
		assert.Equal(t, "pr-1002", result[1].PullRequestID)
	})

	t.Run("success - no PRs for reviewer", func(t *testing.T) {
		mockReviewerRepo.EXPECT().
			ListPRsByReviewer(ctx, "u5").
			Return([]domain.PullRequestShort{}, nil).
			Times(1)

		result, err := service.GetReviewerPRs(ctx, "u5")

		require.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("error - empty reviewer ID", func(t *testing.T) {
		result, err := service.GetReviewerPRs(ctx, "")

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "reviewer_id is required")
	})

	t.Run("error - database error", func(t *testing.T) {
		dbErr := errors.New("database connection failed")
		mockReviewerRepo.EXPECT().
			ListPRsByReviewer(ctx, "u2").
			Return(nil, dbErr).
			Times(1)

		result, err := service.GetReviewerPRs(ctx, "u2")

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "list PRs by reviewer")
	})
}
