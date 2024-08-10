package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/spf13/viper"
)

const (
	MinNumPass = 123456
	MaxNumPass = 999999
)

// Custom type for Unix timestamp
type UnixTime struct {
	Time  time.Time
	Valid bool
}

// UnmarshalJSON for UnixTime
func (ut *UnixTime) UnmarshalJSON(data []byte) error {
	var timestamp int64
	if err := json.Unmarshal(data, &timestamp); err != nil {
		return err
	}

	if timestamp == 0 {
		ut.Valid = false
		return nil
	}

	ut.Time = time.Unix(timestamp, 0).UTC()
	ut.Valid = true
	return nil
}

// pqErrHandler checks if pq error exist in err,
// sends a http response depending on the useCase: user, pitchrequest or deal
// with http status forbidden
func pqErrHandler(ctx *gin.Context, use string, err error) (pqErrExist bool) {
	pqErr, ok := err.(*pq.Error)
	if !ok {
		return false
	}
	var errMsg string

	switch pqErr.Code.Name() {
	case "foreign_key_violation":
		errMsg = fmt.Sprintf("foreign key violated for %v", use)
	case "unique_violation":
		errMsg = fmt.Sprintf("%v already exists, unique violation", use)
	case "not_null_violation":
		errMsg = fmt.Sprintf("%v cannot have null values, not null violation", use)
	default:
		errMsg = fmt.Sprintf("database error: %s", pqErr.Code.Name())
	}
	ctx.JSON(http.StatusConflict, errorResponse(fmt.Errorf(errMsg)))
	return true
}

// sqlNoRowsHandler return true if no rows exist in database
func sqlNoRowsHandler(ctx *gin.Context, err error) (sqlErrrNoRowsExist bool) {
	if err == sql.ErrNoRows {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return true
	}
	return false
}

// Set config for getting env variables
type Config struct {
	DBSource     string `mapstructure:"DB_SOURCE"`
	DBDriver     string `mapstructure:"DB_DRIVER"`
	SymmetricKey string `mapstructure:"SYMMETRIC_KEY"`
}

func LoadConfig(envPath string) (config Config, err error) {
	viper.SetConfigFile(envPath)
	// viper.AddConfigPath(".")

	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return 
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return 
	}
	return
}

func randomPasswordCode() string {
	code := MinNumPass + rand.Intn(MaxNumPass-MinNumPass)
	return strconv.Itoa(code)
}
