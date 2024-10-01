-- name: GetNumberOfAdminUsers :one
SELECT COUNT(*) WHERE EXISTS (
    SELECT * FROM users
    WHERE role = $1
);

-- name: CreateUser :one
INSERT INTO users (username, role, full_name, email, password)
VALUES ( @username, @role, @full_name, @email, @password)
RETURNING users.id;