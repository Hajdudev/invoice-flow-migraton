-- +goose Up
-- +goose StatementBegin

-- Create the auth_tokens table with security and integrity best practices.
CREATE TABLE IF NOT EXISTS auth_tokens (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash BYTEA NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL
);

-- +goose StatementEnd

-- +goose StatementBegin
-- Add an index for fast token lookups, which is critical for performance.
CREATE INDEX IF NOT EXISTS idx_auth_tokens_token_hash ON auth_tokens(token_hash);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS auth_tokens;
-- The index will be dropped automatically when the table is dropped.
-- +goose StatementEnd
