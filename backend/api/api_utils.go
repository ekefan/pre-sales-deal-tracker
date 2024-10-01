package api

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// UserLoginResp sub field in the login response body
type UserLoginResp struct {
	UserID int64 `json:"user_id"`
	Username string `json:"username"`
	Fullname string `json:"fullname"`
	Role string `json:"role"`
	Email string `json:"email"`
	PasswordChanged bool `json:"password_changed"`
	UpdatedAt *string `json:"updated_at"`

}

// errorResponse sends custom error response to client
// with code: a sentinel error eg. NOT_FOUND and err the error
func errorResponse(err error, code string) gin.H {
	return gin.H{
		"code": code,
		"error": err.Error(),
	}
}

func successMessage(msg string) gin.H {
	return gin.H{
		"message": msg,
	}
}



func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func ValidatePassword(hash, password string) (error) {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
