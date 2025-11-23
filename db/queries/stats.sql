-- name: GetUserAssignmentStats :many
SELECT 
    u.user_id,
    u.username,
    u.team_name,
    COUNT(ar.pr_id) as assignments_count
FROM users u
LEFT JOIN assigned_reviewers ar ON u.user_id = ar.reviewer_id
GROUP BY u.user_id, u.username, u.team_name
ORDER BY assignments_count DESC, u.user_id;

-- name: GetPRStats :one
SELECT 
    COUNT(*) as total_prs,
    COUNT(*) FILTER (WHERE status = 'OPEN') as open_prs,
    COUNT(*) FILTER (WHERE status = 'MERGED') as merged_prs
FROM pull_requests;

-- name: GetReviewerWorkload :many
SELECT 
    u.user_id,
    u.username,
    u.team_name,
    COUNT(ar.pr_id) as open_prs_count
FROM users u
LEFT JOIN assigned_reviewers ar ON u.user_id = ar.reviewer_id
LEFT JOIN pull_requests pr ON ar.pr_id = pr.pull_request_id AND pr.status = 'OPEN'
WHERE u.is_active = true
GROUP BY u.user_id, u.username, u.team_name
ORDER BY open_prs_count DESC;