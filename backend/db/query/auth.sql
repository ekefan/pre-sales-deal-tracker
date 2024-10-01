-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1
LIMIT 1;