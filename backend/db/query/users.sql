-- name: GetNumberOfAdminUsers :one
SELECT COUNT(*) FROM users WHERE role = $1;

-- name: CreateUser :one
INSERT INTO users (username, role, full_name, email, password)
VALUES ( @username, @role, @full_name, @email, @password)
RETURNING users.id;

-- name: GetTotalNumOfUsers :one
SELECT COUNT(*) FROM users;

-- name: ListAllUsers :many
SELECT id AS user_id, username, role, email, full_name, password_changed, updated_at, created_at from users
LIMIT $1
OFFSET $2;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: UpdateUser :exec
UPDATE users 
SET username = $2, full_name = $3, role = $4, email = $5, updated_at = $6
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: UpdateUserPassword :exec
UPDATE users
 SET password = $2, password_changed = $3
WHERE id = $1;