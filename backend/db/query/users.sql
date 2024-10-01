-- name: GetNumberOfAdminUsers :one
SELECT COUNT(*) FROM users WHERE role = $1;

-- name: CreateUser :one
INSERT INTO users (username, role, full_name, email, password)
VALUES ( @username, @role, @full_name, @email, @password)
RETURNING users.id;

-- name: GetTotalNumOfUsers :one
SELECT COUNT(*) FROM users;

-- name: ListAllUsers :many
SELECT * from users;


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