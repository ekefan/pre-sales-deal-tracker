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