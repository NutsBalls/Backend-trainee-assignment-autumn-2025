-- name: CreatePullRequest :one
INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status)
VALUES ($1, $2, $3, 'OPEN')
RETURNING *;

-- name: PRExists :one
SELECT EXISTS(SELECT 1 FROM pull_requests WHERE pull_request_id = $1);

-- name: GetPRAuthorId :one
SELECT author_id
FROM pull_requests
WHERE pull_request_id = $1;

-- name: GetPullRequest :one
SELECT *
FROM pull_requests
WHERE pull_request_id = $1;

-- name: MergePullRequest :one
UPDATE pull_requests
SET status = 'MERGED', 
    merged_at = COALESCE(merged_at, NOW())
WHERE pull_request_id = $1
RETURNING *;

-- name: ListPullRequestsByReviewer :many
SELECT pr.pull_request_id, pr.pull_request_name, pr.author_id, pr.status
FROM pull_requests pr
JOIN assigned_reviewers ar ON pr.pull_request_id = ar.pr_id
WHERE ar.reviewer_id = $1
ORDER BY pr.created_at DESC;