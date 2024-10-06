package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

// successMessage sends a custom success response to client
func successMessage() gin.H {
	return gin.H{
		"message": "success",
	}
}

const (
	jsonSource = iota
	uriSource
	querySource
)

// bindClientRequest binds client request to req, which must always be a pointer to the request struct
// the binding source represents the kind of binding to perform on the request
// Fixme: here you're doing too many things: bind the client request and setting the response payload.
// Done: binding the request only
func bindClientRequest(ctx *gin.Context, req any, bindingSource int) error {
	var err error
	switch bindingSource {
	case uriSource:
		err = ctx.ShouldBindUri(req)
	case querySource:
		err = ctx.ShouldBindQuery(req)
	case bindingSource:
		err = ctx.ShouldBindJSON(req)
	}
	return err
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
	SymmetricKey    string        `mapstructure:"SYMMETRIC_KEY"`
	ServerAddres    string        `mapstructure:"SERVER_ADDRESS"`
	DatabaseSource  string        `mapstructure:"DATABASE_SOURCE"`
	MigrationSource string        `mapstructure:"MIGRATION_SOURCE"`
	TokenDuration   time.Duration `mapstructure:"TOKEN_DURATION"`
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

// Pagination holds pagination data for users resource
type Pagination struct {
	TotalRecords int32 `json:"total_records"`
	CurrentPage  int32 `json:"current_page"`
	TotalPages   int32 `json:"total_pages"`
	HasNext      bool  `json:"has_next"`
	HasPrevious  bool  `json:"has_previous"`
}

// generatePagination returns pagination data associated with
// totalRecords, pageID, and pageSize
func generatePagination(totalRecords, pageID, pageSize int32) Pagination {
	var totalPages int32
	if totalRecords/pageSize < 1 {
		totalPages = 1
	} else {
		totalPages = totalRecords / pageSize
	}
	return Pagination{
		TotalRecords: totalRecords,
		CurrentPage:  pageID,
		TotalPages:   totalPages,
		HasNext:      totalPages > pageID,
		HasPrevious:  pageID > 1,
	}
}
