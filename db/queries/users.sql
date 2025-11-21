-- name: InsertUser :one
INSERT INTO users (user_id, username, team_name, is_active)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UserExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE user_id = $1);

-- name: UpdateUser :exec
UPDATE users
SET
    username = $2,
    team_name = $3,
    is_active = $4
WHERE user_id = $1;

-- name: GetUser :one
SELECT user_id, username, team_name, is_active
FROM users
WHERE user_id = $1;

-- name: GetUsersByTeam :many
SELECT user_id, username, team_name, is_active
FROM users
WHERE team_name = $1
ORDER BY user_id;

-- name: SetUserActivity :one
UPDATE users
SET is_active = $2
WHERE user_id = $1
RETURNING *;

-- name: GetActiveCandidatesForPR :many
SELECT user_id, username
FROM users
WHERE team_name = $1 
  AND is_active = true 
  AND user_id != $2
ORDER BY RANDOM()
LIMIT 2;

-- name: GetActiveCandidatesForReassignment :many
SELECT user_id
FROM users
WHERE team_name = $1 
  AND is_active = true 
  AND user_id != $2
  AND user_id NOT IN (
    SELECT reviewer_id 
    FROM assigned_reviewers 
    WHERE pr_id = $3
  )
ORDER BY RANDOM()
LIMIT 1;