package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// FixME: create a struct exposed to hold the error. Then, create a "constructor" function to instantiate it.
// FixME: it's redundant. We already have this information in the response Status Code Header. //DONE
type ErrResp struct {
	Code    string            `json:"code"`
	Error   string            `json:"error"`
	Details map[string]string `json:"details"`
}

func NewErrResp(code, err string) ErrResp {
	return ErrResp{
		Code:  code,
		Error: err,
	}
}

// handleServerError snipet for returning server error to the client
func handleServerError(ctx *gin.Context, err error) {
	slog.Error("server error", "the error", err.Error())
	er := NewErrResp("SERVER_ERROR", "error from the server")
	er.Details = map[string]string{
		"server_error": "server failed to process request",
	}
	ctx.JSON(http.StatusInternalServerError, er)
}

var (
	customNotFound = errors.New("not found")
	deleteMasterAdminErr = errors.New("a master user must exist in the system")
)
// handleDbError handles expected db errors,
// takes the context, the error and a possible error detail
// returns true if an error was handled and false if no predicted db error is handled
//
// Note: If err can not be associated with pre-defined db errors, err.Error() will contain custom err message
func handleDbError(ctx *gin.Context, err error, detail string) bool {
	slog.Error("from the database", "error", err.Error())
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		er := NewErrResp("STATUS_CONFLICT", "a similar resource already exists")
		er.Details = map[string]string{"data conflict": detail}
		ctx.JSON(http.StatusConflict, er)
		return true
	}
	if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, customNotFound){
		er := NewErrResp("NOT_FOUND", "requested resource doesn't exist")
		er.Details = map[string]string{"not found": detail}
		ctx.JSON(http.StatusNotFound, er)
		return true
	}

	if errors.Is(err, deleteMasterAdminErr) {
		er := NewErrResp("FORBIDDEN", "can not delete master admin")
		er.Details = map[string]string{"deleting admin": detail}
		ctx.JSON(http.StatusForbidden, er)
	}
	return false
}


// handlePasswordValidationError snipet for returning error details 
// to client after a failed password validation
func handlePasswordValidationError(ctx *gin.Context, err error) {
	slog.Error("error validating user password", "the error", err.Error())
	er := NewErrResp("UNAUTHORIZED", "invalid password")
	er.Details = map[string]string{
		"password": "user password is wrong, input the correct passowd",
	}
	ctx.JSON(http.StatusUnauthorized, er)
}

// validatioTagsMap holds every binding tag used in pre-sales-deal-tracker end point handlers
var validationTagsMap = map[string]string{
	"required": "required",
	"email":    "is invalid",
	"min":      "min should be",
	"max":      "max should be",
	"alphanum": "must contain alphanum only",
	"oneof":    "must be one of:",
	"eq":       "must be equal to",
}

// composeMsg returns an error message for a every tag in e, the validator field error
func composeMsg(e validator.FieldError) string {
	if errMsg, ok := validationTagsMap[e.Tag()]; ok {
		if e.Param() != "" {
			return fmt.Sprintf("%v: %v", errMsg, e.Param())
		}
		return fmt.Sprintf("%v", errMsg)
	}
	return "oops, we couldn't determine the field error this time, please re-check the request params"
}

// translateError takes a slice of validation errors, transalates them to a map of field to error message
func translateError(validationErr validator.ValidationErrors) map[string]string {
	vErr := map[string]string{}
	for _, ve := range validationErr {
		errMsg := composeMsg(ve)
		vErr[ve.Field()] = errMsg
	}
	return vErr
}

// handleClientReqError handles bad requests from client requests
func handleClientReqError(ctx *gin.Context, err error) {
	// possible pre-defined error types
	var (
		validationErr validator.ValidationErrors
		intErr        *strconv.NumError
		jsonSyntaxErr *json.SyntaxError
	)
	er := NewErrResp("BAD_REQUEST", "invalid request parameters")
	if errors.As(err, &validationErr) {
		er.Details = translateError(validationErr)
		ctx.JSON(http.StatusBadRequest, er)
		return
	}
	if errors.As(err, &intErr) {
		er.Details = map[string]string{
			"numeric inputs": fmt.Sprintf("only numeric inputs, \"%v\" is not valid", intErr.Num),
		}
		ctx.JSON(http.StatusBadRequest, er)
		return
	}
	if errors.As(err, &jsonSyntaxErr) {
		er.Details = map[string]string{
			"json syntax": fmt.Sprintf("json syntax error - %v", jsonSyntaxErr.Error()),
		}
		ctx.JSON(http.StatusBadRequest, er)
		return
	}
	slog.Error("undetermined client error", "the error", fmt.Sprintf("%T", err))
	er.Error = "ops!!!! couldn't determine the exact error this time, you might relate to the details"
	er.Details = map[string]string{
		"unknown details": err.Error(),
	}
	ctx.JSON(http.StatusBadRequest, er)

}
