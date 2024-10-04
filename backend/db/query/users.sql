-- name: GetNumberOfAdminUsers :one
SELECT COUNT(*) FROM users WHERE role = $1;

-- name: CreateUser :one
INSERT INTO users (username, role, full_name, email, password)
VALUES ( @username, @role, @full_name, @email, @password)
RETURNING users.id;

-- name: GetTotalNumOfUsers :one
SELECT COUNT(*) FROM users;

-- name: ListAllUsers :many
SELECT 
    id AS user_id,
    username,
    role,
    email,
    full_name,
    password_changed,
    to_char(updated_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"') AS updated_at,
    to_char(created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"') AS created_at
FROM users
LIMIT $1
OFFSET $2;



-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: UpdateUser :one
UPDATE users 
 SET username = $2, full_name = $3, role = $4, email = $5, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: UpdateUserPassword :exec
UPDATE users
 SET password = $2, password_changed = $3
WHERE id = $1;