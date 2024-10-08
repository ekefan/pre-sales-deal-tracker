// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Deal struct {
	ID               int64            `json:"id"`
	PitchID          *int64           `json:"pitch_id"`
	SalesRepName     string           `json:"sales_rep_name"`
	CustomerName     string           `json:"customer_name"`
	ServicesToRender []string         `json:"services_to_render"`
	Status           string           `json:"status"`
	Department       string           `json:"department"`
	NetTotalCost     pgtype.Numeric   `json:"net_total_cost"`
	Profit           pgtype.Numeric   `json:"profit"`
	CreatedAt        pgtype.Timestamp `json:"created_at"`
	UpdatedAt        pgtype.Timestamp `json:"updated_at"`
	ClosedAt         pgtype.Timestamp `json:"closed_at"`
	Awarded          bool             `json:"awarded"`
}

type PitchRequest struct {
	ID              int64            `json:"id"`
	UserID          int64            `json:"user_id"`
	CustomerName    string           `json:"customer_name"`
	CustomerRequest []string         `json:"customer_request"`
	AdminTask       string           `json:"admin_task"`
	AdminDeadline   pgtype.Timestamp `json:"admin_deadline"`
	AdminViewed     bool             `json:"admin_viewed"`
	CreatedAt       pgtype.Timestamp `json:"created_at"`
	UpdatedAt       pgtype.Timestamp `json:"updated_at"`
}

type User struct {
	ID              int64            `json:"id"`
	Username        string           `json:"username"`
	Role            string           `json:"role"`
	FullName        string           `json:"full_name"`
	Email           string           `json:"email"`
	Password        string           `json:"password"`
	IsMaster        bool             `json:"is_master"`
	PasswordChanged bool             `json:"password_changed"`
	UpdatedAt       pgtype.Timestamp `json:"updated_at"`
	CreatedAt       pgtype.Timestamp `json:"created_at"`
}
