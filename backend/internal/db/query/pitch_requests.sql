-- name: CreatePitchRequest :execrows
INSERT INTO pitch_requests (user_id, customer_name, customer_request, admin_task, admin_deadline)
VALUES (@user_id, @customer_name, @customer_request, @admin_task, @admin_deadline);

-- name: GetPitchRequestsPaginated :one
WITH pitch_data AS (
    SELECT
        pitch_requests.id as pitch_id, 
        pitch_requests.user_id,
        pitch_requests.customer_name, 
        pitch_requests.customer_request, 
        pitch_requests.admin_task, 
        pitch_requests.admin_deadline, 
        pitch_requests.admin_viewed, 
        to_char(pitch_requests.updated_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"') AS updated_at,
        to_char(pitch_requests.created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"') AS created_at
    FROM pitch_requests WHERE pitch_requests.user_id = $1
    LIMIT $2 OFFSET $3
)
SELECT
    (SELECT COUNT(*) FROM pitch_requests) AS total_pitch_requests,
    json_agg(pitch_data) as pitch_requests
FROM pitch_data;


-- name: UpdatePitchRequest :execrows
UPDATE pitch_requests
    set admin_viewed = sqlc.arg('admin_viewed'), customer_request = sqlc.narg('customer_request')
WHERE pitch_requests.id = sqlc.arg('pitch_id');

-- name: DeletePitchRequest :execrows
DELETE FROM pitch_requests WHERE pitch_requests.id = @pitch_id;