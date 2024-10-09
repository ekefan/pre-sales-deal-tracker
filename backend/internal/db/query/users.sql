-- name: GetNumberOfAdminUsers :one
SELECT COUNT(*) FROM users WHERE role = $1;

-- name: CreateUser :one
INSERT INTO users (username, role, full_name, email, password)
VALUES ( @username, @role, @full_name, @email, @password)
RETURNING users.id;


-- name: CreateMasterUser :one
INSERT INTO users (username, role, full_name, email, password, is_master)
VALUES (@username, @role, @full_name, @email, @password, true)
RETURNING users.id;

-- name: GetMasterUser :one
SELECT id as user_id FROM users
WHERE is_master = true;

-- name: GetUserFullName :one
SELECT users.full_name FROM users
WHERE users.id = $1;

-- name: GetUserPaginated :one
WITH user_data AS (
    SELECT 
        id as user_id,
        username,
        role,
        email,
        full_name,
        password_changed,
        to_char(updated_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"') AS updated_at,
        to_char(created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"') AS created_at
    FROM users
    ORDER BY id
    LIMIT $1 OFFSET $2
)
SELECT 
    (SELECT COUNT(*) FROM users) AS total_users,
    json_agg(user_data) AS users
FROM user_data;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: UpdateUser :execrows
UPDATE users 
 SET username = $2, full_name = $3, role = $4, email = $5, updated_at = NOW()
WHERE id = $1;

-- name: DeleteUser :execrows
DELETE FROM users
WHERE users.id = $1 AND (
    users.is_master != true
    );

-- name: UpdateUserPassword :exec
UPDATE users
 SET password = $2, password_changed = $3
WHERE id = $1;