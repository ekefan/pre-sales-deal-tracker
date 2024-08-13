// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: admin.sql

package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

const adminDealExists = `-- name: AdminDealExists :one
SELECT EXISTS (
    SELECT 1
    FROM deals
    WHERE id = $1
)
`

func (q *Queries) AdminDealExists(ctx context.Context, id int64) (bool, error) {
	row := q.queryRow(ctx, q.adminDealExistsStmt, adminDealExists, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const adminDeleteDeal = `-- name: AdminDeleteDeal :exec
DELETE FROM deals
WHERE id = $1
`

func (q *Queries) AdminDeleteDeal(ctx context.Context, id int64) error {
	_, err := q.exec(ctx, q.adminDeleteDealStmt, adminDeleteDeal, id)
	return err
}

const adminDeleteUser = `-- name: AdminDeleteUser :exec
DELETE FROM users
WHERE id = $1
`

func (q *Queries) AdminDeleteUser(ctx context.Context, id int64) error {
	_, err := q.exec(ctx, q.adminDeleteUserStmt, adminDeleteUser, id)
	return err
}

const adminGetDealForUpdate = `-- name: AdminGetDealForUpdate :one
SELECT id, pitch_id, sales_rep_name, customer_name, service_to_render, status, status_tag, current_pitch_request, net_total_cost, profit, created_at, updated_at, closed_at, awarded FROM deals
WHERE id = $1
LIMIT 1
FOR UPDATE
`

func (q *Queries) AdminGetDealForUpdate(ctx context.Context, id int64) (Deal, error) {
	row := q.queryRow(ctx, q.adminGetDealForUpdateStmt, adminGetDealForUpdate, id)
	var i Deal
	err := row.Scan(
		&i.ID,
		&i.PitchID,
		&i.SalesRepName,
		&i.CustomerName,
		pq.Array(&i.ServiceToRender),
		&i.Status,
		&i.StatusTag,
		&i.CurrentPitchRequest,
		&i.NetTotalCost,
		&i.Profit,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ClosedAt,
		&i.Awarded,
	)
	return i, err
}

const adminUpdateDeal = `-- name: AdminUpdateDeal :one
UPDATE deals
    set service_to_render = $2, status = $3,
    status_tag = $4, current_pitch_request = $5, updated_at = $7,
    closed_at = $6
WHERE id = $1
RETURNING id, pitch_id, sales_rep_name, customer_name, service_to_render, status, status_tag, current_pitch_request, net_total_cost, profit, created_at, updated_at, closed_at, awarded
`

type AdminUpdateDealParams struct {
	ID                  int64
	ServiceToRender     []string
	Status              string
	StatusTag           string
	CurrentPitchRequest string
	ClosedAt            time.Time
	UpdatedAt           time.Time
}

func (q *Queries) AdminUpdateDeal(ctx context.Context, arg AdminUpdateDealParams) (Deal, error) {
	row := q.queryRow(ctx, q.adminUpdateDealStmt, adminUpdateDeal,
		arg.ID,
		pq.Array(arg.ServiceToRender),
		arg.Status,
		arg.StatusTag,
		arg.CurrentPitchRequest,
		arg.ClosedAt,
		arg.UpdatedAt,
	)
	var i Deal
	err := row.Scan(
		&i.ID,
		&i.PitchID,
		&i.SalesRepName,
		&i.CustomerName,
		pq.Array(&i.ServiceToRender),
		&i.Status,
		&i.StatusTag,
		&i.CurrentPitchRequest,
		&i.NetTotalCost,
		&i.Profit,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ClosedAt,
		&i.Awarded,
	)
	return i, err
}

const adminUpdateUser = `-- name: AdminUpdateUser :one
UPDATE users
    set full_name = $2, email = $3, username = $4, updated_at = $5
WHERE id = $1
RETURNING id, username, role, full_name, email, password, password_changed, updated_at, created_at
`

type AdminUpdateUserParams struct {
	ID        int64
	FullName  string
	Email     string
	Username  string
	UpdatedAt time.Time
}

func (q *Queries) AdminUpdateUser(ctx context.Context, arg AdminUpdateUserParams) (User, error) {
	row := q.queryRow(ctx, q.adminUpdateUserStmt, adminUpdateUser,
		arg.ID,
		arg.FullName,
		arg.Email,
		arg.Username,
		arg.UpdatedAt,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Role,
		&i.FullName,
		&i.Email,
		&i.Password,
		&i.PasswordChanged,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const adminUserExists = `-- name: AdminUserExists :one
SELECT EXISTS (
    SELECT 1 
    FROM users
    WHERE id = $1
)
`

func (q *Queries) AdminUserExists(ctx context.Context, id int64) (bool, error) {
	row := q.queryRow(ctx, q.adminUserExistsStmt, adminUserExists, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const adminViewDeals = `-- name: AdminViewDeals :many
SELECT id, pitch_id, sales_rep_name, customer_name, service_to_render, status, status_tag, current_pitch_request, net_total_cost, profit, created_at, updated_at, closed_at, awarded FROM deals
ORDER BY id
LIMIT $1
OFFSET $2
`

type AdminViewDealsParams struct {
	Limit  int32
	Offset int32
}

func (q *Queries) AdminViewDeals(ctx context.Context, arg AdminViewDealsParams) ([]Deal, error) {
	rows, err := q.query(ctx, q.adminViewDealsStmt, adminViewDeals, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Deal{}
	for rows.Next() {
		var i Deal
		if err := rows.Scan(
			&i.ID,
			&i.PitchID,
			&i.SalesRepName,
			&i.CustomerName,
			pq.Array(&i.ServiceToRender),
			&i.Status,
			&i.StatusTag,
			&i.CurrentPitchRequest,
			&i.NetTotalCost,
			&i.Profit,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ClosedAt,
			&i.Awarded,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const adminViewUsers = `-- name: AdminViewUsers :many
SELECT id, username, role, full_name, email, password, password_changed, updated_at, created_at FROM users
ORDER BY id
LIMIT $1
OFFSET $2
`

type AdminViewUsersParams struct {
	Limit  int32
	Offset int32
}

func (q *Queries) AdminViewUsers(ctx context.Context, arg AdminViewUsersParams) ([]User, error) {
	rows, err := q.query(ctx, q.adminViewUsersStmt, adminViewUsers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.Role,
			&i.FullName,
			&i.Email,
			&i.Password,
			&i.PasswordChanged,
			&i.UpdatedAt,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const createDeal = `-- name: CreateDeal :one
INSERT INTO deals (
    pitch_id, sales_rep_name, customer_name, service_to_render, status, status_tag, current_pitch_request
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING id, pitch_id, sales_rep_name, customer_name, service_to_render, status, status_tag, current_pitch_request, net_total_cost, profit, created_at, updated_at, closed_at, awarded
`

type CreateDealParams struct {
	PitchID             sql.NullInt64
	SalesRepName        string
	CustomerName        string
	ServiceToRender     []string
	Status              string
	StatusTag           string
	CurrentPitchRequest string
}

func (q *Queries) CreateDeal(ctx context.Context, arg CreateDealParams) (Deal, error) {
	row := q.queryRow(ctx, q.createDealStmt, createDeal,
		arg.PitchID,
		arg.SalesRepName,
		arg.CustomerName,
		pq.Array(arg.ServiceToRender),
		arg.Status,
		arg.StatusTag,
		arg.CurrentPitchRequest,
	)
	var i Deal
	err := row.Scan(
		&i.ID,
		&i.PitchID,
		&i.SalesRepName,
		&i.CustomerName,
		pq.Array(&i.ServiceToRender),
		&i.Status,
		&i.StatusTag,
		&i.CurrentPitchRequest,
		&i.NetTotalCost,
		&i.Profit,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ClosedAt,
		&i.Awarded,
	)
	return i, err
}

const createNewUser = `-- name: CreateNewUser :one
INSERT INTO users (
    username, role, full_name, email, password, password_changed
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING id, username, role, full_name, email, password, password_changed, updated_at, created_at
`

type CreateNewUserParams struct {
	Username        string
	Role            string
	FullName        string
	Email           string
	Password        string
	PasswordChanged bool
}

func (q *Queries) CreateNewUser(ctx context.Context, arg CreateNewUserParams) (User, error) {
	row := q.queryRow(ctx, q.createNewUserStmt, createNewUser,
		arg.Username,
		arg.Role,
		arg.FullName,
		arg.Email,
		arg.Password,
		arg.PasswordChanged,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Role,
		&i.FullName,
		&i.Email,
		&i.Password,
		&i.PasswordChanged,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const forgotPassword = `-- name: ForgotPassword :one
SELECT id, username, role, full_name, email, password, password_changed, updated_at, created_at FROM users
WHERE email = $1
LIMIT 1
FOR UPDATE
`

func (q *Queries) ForgotPassword(ctx context.Context, email string) (User, error) {
	row := q.queryRow(ctx, q.forgotPasswordStmt, forgotPassword, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Role,
		&i.FullName,
		&i.Email,
		&i.Password,
		&i.PasswordChanged,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getUserForUpdate = `-- name: GetUserForUpdate :one
SELECT id, username, role, full_name, email, password, password_changed, updated_at, created_at FROM users
WHERE id = $1
LIMIT 1
FOR UPDATE
`

func (q *Queries) GetUserForUpdate(ctx context.Context, id int64) (User, error) {
	row := q.queryRow(ctx, q.getUserForUpdateStmt, getUserForUpdate, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Role,
		&i.FullName,
		&i.Email,
		&i.Password,
		&i.PasswordChanged,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const updatePassWord = `-- name: UpdatePassWord :exec
UPDATE users
    set password = $2, password_changed = $3, updated_at = $4
WHERE id = $1
`

type UpdatePassWordParams struct {
	ID              int64
	Password        string
	PasswordChanged bool
	UpdatedAt       time.Time
}

func (q *Queries) UpdatePassWord(ctx context.Context, arg UpdatePassWordParams) error {
	_, err := q.exec(ctx, q.updatePassWordStmt, updatePassWord,
		arg.ID,
		arg.Password,
		arg.PasswordChanged,
		arg.UpdatedAt,
	)
	return err
}
