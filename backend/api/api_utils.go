package api

import (
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

func ValidatePassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// Config holds environment variables needed to run server
type Config struct {
	SymmetricKey string `mapstructure:"SYMMETRIC_KEY"`
	ServerAddres string `mapstructure:"SERVER_ADDRESS"`
	DatabaseUrl  string `mapstructure:"DATABASE_URL"`
}

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