-- +goose Up
CREATE TABLE assigned_reviewers (
    pr_id TEXT NOT NULL REFERENCES pull_requests(pull_request_id) ON DELETE CASCADE,
    reviewer_id TEXT NOT NULL REFERENCES users(user_id),
    PRIMARY KEY (pr_id, reviewer_id)
);

-- +goose Down
DROP TABLE assigned_reviewers;
