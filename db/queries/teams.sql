-- name: CreateTeam :one
INSERT INTO teams (team_name)
VALUES ($1)
RETURNING team_name;

-- name: GetTeam :one
SELECT team_name
FROM teams
WHERE team_name = $1;

-- name: TeamExists :one
SELECT EXISTS(SELECT 1 FROM teams WHERE team_name = $1);
