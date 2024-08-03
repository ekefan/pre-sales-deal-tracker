-- name: CreateDeal :one
INSERT INTO deals (
    pitch_id, sales_rep_name, customer_name, service_to_render, status, status_tag, current_pitch_request
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;


-- name: CreateNewUser :one
INSERT INTO users (
    username, role, full_name, email, password
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetUserForUpdate :one
SELECT * FROM users
WHERE id = $1
LIMIT 1
FOR UPDATE;

-- name: AdminUpdateUser :one
UPDATE users
    set full_name = $2, email = $3, password = $4, username = $5, updated_at = $6
WHERE id = $1
RETURNING *;

-- name: AdminGetDealForUpdate :one
SELECT * FROM deals
WHERE id = $1
LIMIT 1
FOR UPDATE;

-- name: AdminUpdateDeal :one
UPDATE deals
    set service_to_render = $2, status = $3,
    status_tag = $4, current_pitch_request = $5, updated_at = $7,
    closed_at = $6
WHERE id = $1
RETURNING *;

-- name: AdminUserExists :one
SELECT EXISTS (
    SELECT 1 
    FROM users
    WHERE id = $1
);

-- name: AdminDeleteUser :exec
DELETE FROM users
WHERE id = $1;


-- name: AdminDealExists :one
SELECT EXISTS (
    SELECT 1
    FROM deals
    WHERE id = $1
);

-- name: AdminDeleteDeal :exec
DELETE FROM deals
WHERE id = $1;

-- name: AdminViewUsers :many
SELECT * FROM users
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: AdminViewDeals :many
SELECT * FROM deals
ORDER BY id
LIMIT $1
OFFSET $2;
