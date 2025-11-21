-- name: AddReviewer :exec
INSERT INTO assigned_reviewers (pr_id, reviewer_id)
VALUES ($1, $2);

-- name: RemoveReviewer :exec
DELETE FROM assigned_reviewers
WHERE pr_id = $1 AND reviewer_id = $2;

-- name: IsReviewerAssigned :one
SELECT EXISTS (
  SELECT 1
  FROM assigned_reviewers
  WHERE pr_id = $1 AND reviewer_id = $2
) AS is_assigned;

-- name: ReplaceReviewer :exec
UPDATE assigned_reviewers
SET reviewer_id = $3
WHERE pr_id = $1 AND reviewer_id = $2;

-- name: GetAssignedReviewers :many
SELECT reviewer_id
FROM assigned_reviewers
WHERE pr_id = $1
ORDER BY reviewer_id;