-- name: UpdateUser :one
UPDATE users
    set username = $2, password = $1
WHERE id = $1
RETURNING *;

-- name: CreatePitchRequest :one
INSERT INTO pitch_requests (
    sales_rep_id, status, customer_name, pitch_tag, customer_request
) VALUES (
    $1, $2, $3, $4, $5
) 
RETURNING *;

-- name: GetPitchRequestForUpdate :one
SELECT * FROM pitch_requests
WHERE id = $1
LIMIT 1
FOR UPDATE;

-- name: UpdatePitchRequest :one
UPDATE pitch_requests
    set status = $2, pitch_tag = $3, customer_request = $4, admin_viewed = $5
WHERE id = $1
RETURNING *;

-- name: DeletePitchRequest :exec
DELETE FROM pitch_requests
WHERE sales_rep_id = $1;