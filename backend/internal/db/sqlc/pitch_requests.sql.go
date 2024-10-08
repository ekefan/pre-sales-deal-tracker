// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: pitch_requests.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createPitchRequest = `-- name: CreatePitchRequest :execrows
INSERT INTO pitch_requests (user_id, customer_name, customer_request, admin_task, admin_deadline)
VALUES ($1, $2, $3, $4, $5)
`

type CreatePitchRequestParams struct {
	UserID          int64            `json:"user_id"`
	CustomerName    string           `json:"customer_name"`
	CustomerRequest []string         `json:"customer_request"`
	AdminTask       string           `json:"admin_task"`
	AdminDeadline   pgtype.Timestamp `json:"admin_deadline"`
}

func (q *Queries) CreatePitchRequest(ctx context.Context, arg CreatePitchRequestParams) (int64, error) {
	result, err := q.db.Exec(ctx, createPitchRequest,
		arg.UserID,
		arg.CustomerName,
		arg.CustomerRequest,
		arg.AdminTask,
		arg.AdminDeadline,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const deletePitchRequest = `-- name: DeletePitchRequest :execrows
DELETE FROM pitch_requests WHERE pitch_requests.id = $1
`

func (q *Queries) DeletePitchRequest(ctx context.Context, pitchID int64) (int64, error) {
	result, err := q.db.Exec(ctx, deletePitchRequest, pitchID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const getPitchRequestsPaginated = `-- name: GetPitchRequestsPaginated :one
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
FROM pitch_data
`

type GetPitchRequestsPaginatedParams struct {
	UserID int64 `json:"user_id"`
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type GetPitchRequestsPaginatedRow struct {
	TotalPitchRequests int64  `json:"total_pitch_requests"`
	PitchRequests      []byte `json:"pitch_requests"`
}

func (q *Queries) GetPitchRequestsPaginated(ctx context.Context, arg GetPitchRequestsPaginatedParams) (GetPitchRequestsPaginatedRow, error) {
	row := q.db.QueryRow(ctx, getPitchRequestsPaginated, arg.UserID, arg.Limit, arg.Offset)
	var i GetPitchRequestsPaginatedRow
	err := row.Scan(&i.TotalPitchRequests, &i.PitchRequests)
	return i, err
}

const updatePitchRequest = `-- name: UpdatePitchRequest :execrows
UPDATE pitch_requests
    set admin_viewed = $1, customer_request = $2
WHERE pitch_requests.id = $3
`

type UpdatePitchRequestParams struct {
	AdminViewed     bool     `json:"admin_viewed"`
	CustomerRequest []string `json:"customer_request"`
	PitchID         int64    `json:"pitch_id"`
}

func (q *Queries) UpdatePitchRequest(ctx context.Context, arg UpdatePitchRequestParams) (int64, error) {
	result, err := q.db.Exec(ctx, updatePitchRequest, arg.AdminViewed, arg.CustomerRequest, arg.PitchID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}