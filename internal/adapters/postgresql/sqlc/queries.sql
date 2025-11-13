-- name: GetUser :one
SELECT * FROM users
WHERE id = $1
LIMIT 1;

-- name: InsertRefreshToken :exec
INSERT INTO auth_tokens (
    user_id,token_hash,expires_at 
) VALUES ($1,$2,$3 );
