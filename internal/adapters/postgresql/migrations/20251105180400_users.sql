-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(100) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Insert some test data
INSERT INTO users (email, username, first_name, last_name, is_active) VALUES
    ('john.doe@example.com', 'johndoe', 'John', 'Doe', true),
    ('jane.smith@example.com', 'janesmith', 'Jane', 'Smith', true),
    ('bob.wilson@example.com', 'bobwilson', 'Bob', 'Wilson', false),
    ('alice.johnson@example.com', 'alicej', 'Alice', 'Johnson', true),
    ('charlie.brown@example.com', 'charlieb', 'Charlie', 'Brown', true);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
