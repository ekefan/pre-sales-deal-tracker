package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ekefan/pre-sales-deal-tracker/backend/middleware"
	"github.com/ekefan/pre-sales-deal-tracker/backend/token"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

// variables for designing error responses
var (
	statusCode               int
	errCode, errMsg, details string
)

// errorResponse sends custom error response to client
// with code: a sentinel error eg. NOT_FOUND and err the error
func errorResponse(statusCode int, code, message, details string) gin.H {
	errResp := struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Details string `json:"details"`
	}{
		Code:    code,
		Message: message,
		Details: details,
	}
	return gin.H{
		"status_code": statusCode,
		"error":       errResp,
	}
}

// handleServerError snipet for returning server error to the client
func handleServerError(ctx *gin.Context, err error) {
	statusCode = http.StatusInternalServerError
	errCode = "SERVER_ERROR"
	errMsg = "a server error"
	details = err.Error()
	ctx.JSON(statusCode, errorResponse(statusCode, errCode, errMsg, details))
}

// // successMessage sends a custom success response to client
// func successMessage(msg string) gin.H {
// 	return gin.H{
// 		"message": msg,
// 	}
// }

const (
	jsonSource = iota
	uriSource
	querySource
)

// bindClientRequest binds client request to req, which must always be a pointer to the request
// the binding source represents which kind of binding to perform on the request
func bindClientRequest(ctx *gin.Context, req any, bindingSource int) error {
	statusCode := http.StatusBadRequest
	errCode := "BAD_REQUEST"
	errMsg := "failed to bind client request"
	details := "invalid request params have been sent, make sure request params are valid"
	switch bindingSource {
	case uriSource:
		if err := ctx.ShouldBindUri(req); err != nil {
			ctx.JSON(statusCode, errorResponse(statusCode, errCode, errMsg, details))
			return err
		}
	case querySource:
		if err := ctx.ShouldBindQuery(req); err != nil {
			ctx.JSON(statusCode, errorResponse(statusCode, errCode, errMsg, details))
			return err
		}
	case bindingSource:
		if err := ctx.ShouldBindJSON(req); err != nil {
			ctx.JSON(statusCode, errorResponse(statusCode, errCode, errMsg, details))
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

const (
	adminRole   = "admin"
	salesRole   = "sales"
	managerRole = "manager"
)

// authAccess ensures the request is from a user with role in roles
func authAccess(ctx *gin.Context, roles []string) bool {
	payload := ctx.MustGet(middleware.AuthPayloadKey).(*token.Payload)
	for _, role := range roles {
		if role == payload.Role {
			return true
		}
	}
	errMsg := "user not authorized to access resource"
	errCode := "FORBIDDEN"
	errDetail := fmt.Sprintf("resource can only be accessed by user with such role: %v", strings.Join(roles, ", "))
	statusCode := http.StatusForbidden
	ctx.JSON(statusCode,
		errorResponse(statusCode, errCode, errMsg, errDetail))
	return false
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

// pgxError handles database error, returns true if error is a pgxError
// else returns false
func pgxError(ctx *gin.Context, err error, msg, errDetails string) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		statusCode = http.StatusConflict
		errCode = "STATUS_CONFLICT"
		errMsg = msg
		details = errDetails
		ctx.JSON(statusCode, errorResponse(statusCode, errCode, errMsg, details))
		return true
	}
	if errors.Is(err, pgx.ErrNoRows) {
		statusCode = http.StatusNotFound
		errCode = "NOT_FOUND"
		errMsg = msg
		details = errDetails
		ctx.JSON(statusCode, errorResponse(statusCode, errCode, errMsg, details))
		return true
	}
	return false
}
