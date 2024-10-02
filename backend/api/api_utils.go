package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

// UserLoginResp sub field in the login response body
type UserLoginResp struct {
	UserID          int64   `json:"user_id"`
	Username        string  `json:"username"`
	Fullname        string  `json:"fullname"`
	Role            string  `json:"role"`
	Email           string  `json:"email"`
	PasswordChanged bool    `json:"password_changed"`
	UpdatedAt       *string `json:"updated_at"`
}

// errorResponse sends custom error response to client
// with code: a sentinel error eg. NOT_FOUND and err the error
func errorResponse(err error, code string) gin.H {
	return gin.H{
		"code":  code,
		"error": err.Error(),
	}
}

// successMessage sends a custom success response to client
func successMessage(msg string) gin.H {
	return gin.H{
		"message": msg,
	}
}

const (
	jsonSource = iota
	uriSource
	querySource
)

// bindClientRequest takes a pointer to the request and a flag	representing
// the binding source
func bindClientRequest(ctx *gin.Context, req any, bindingSource int) error {
	switch bindingSource {
	case uriSource:
		if err := ctx.ShouldBindUri(req); err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err, "BAD_REQUEST"))
			return err
		}
	case querySource:
		if err := ctx.ShouldBindQuery(req); err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err, "BAD_REQUEST"))
			return err
		}
	case bindingSource:
		if err := ctx.ShouldBindJSON(req); err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err, "BAD_REQUEST"))
			return err
		}
	default:
		return errors.New("bindingSource not supported")
	}
	return nil
}

// HashPassword uses bcrypt to generate a hash from password
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// ValidatePassword compares a hash with a password to see if hash was created
// from the password
func ValidatePassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// Config holds environment variables needed to run server
type Config struct {
	SymmetricKey string `mapstructure:"SYMMETRIC_KEY"`
	ServerAddres string `mapstructure:"SERVER_ADDRESS"`
	DatabaseUrl  string `mapstructure:"DATABASE_URL"`
}

// ReadConfigFiles uses viper to read environment config or variables into Config
func ReadConfigFiles(configPath string) (*Config, error) {
	viper.SetConfigName("app")
	viper.AddConfigPath(configPath)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
