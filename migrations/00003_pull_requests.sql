-- +goose Up
CREATE TABLE pull_requests (
    pull_request_id TEXT PRIMARY KEY,
    pull_request_name TEXT NOT NULL,
    author_id TEXT NOT NULL REFERENCES users(user_id),
    status TEXT NOT NULL CHECK (status IN ('OPEN','MERGED')),
    created_at TIMESTAMP DEFAULT NOW(),
    merged_at TIMESTAMP
);

-- +goose Down
DROP TABLE pull_requests;
