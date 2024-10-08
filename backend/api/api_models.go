package api

type User struct {
	UserID          int64  `json:"user_id"`
	Username        string `json:"username"`
	Role            string `json:"role"`
	Email           string `json:"email"`
	FullName        string `json:"full_name"`
	PasswordChanged bool   `json:"password_changed"`
	UpdatedAt       string `json:"updated_at"`
	CreatedAt       string `json:"created_at"`
}

type Deal struct {
	DealID           int64    `json:"deal_id"`
	PitchId          int64    `json:"pitch_id"`
	SalesRepName     string   `json:"sales_rep_name"`
	CustomerName     string   `json:"customer_name"`
	ServicesToRender []string `json:"services_to_render"`
	Status           string   `json:"status"`
	Department       string   `json:"department"`
	NetTotalCOst     float64  `json:"net_total_cost"`
	Profit           float64  `json:"profit"`
	Awarded          bool     `json:"awarded"`
	UpdatedAt        string   `json:"updated_at"`
	CreatedAt        string   `json:"created_at"`
}

type PitchRequest struct {
	PitchID         int64    `json:"pitch_id"`
	UserID          int64    `json:"user_id"`
	CustomerName    string   `json:"customer_name"`
	CustomerRequest []string `json:"customer_request"`
	AdminTask       string   `json:"admin_task"`
	AdminDeadline   string   `json:"admin_deadline"`
	AdminViewed     bool     `json:"admin_viewed"`
	UpdatedAt       string   `json:"updated_at"`
	CreatedAt       string   `json:"created_at"`
}
