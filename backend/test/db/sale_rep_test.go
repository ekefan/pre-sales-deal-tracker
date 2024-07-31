package test

import (
	"context"
	"testing"
	"time"

	db "github.com/ekefan/deal-tracker/internal/db/sqlc"
	util "github.com/ekefan/deal-tracker/test/util"
	"github.com/stretchr/testify/require"
)

func createPitch(t *testing.T, salesRepId int64) db.PitchRequest {
	
	// args fields neeed to create a newPitch row
	args := db.CreatePitchRequestParams{
		SalesRepID:      salesRepId,
		Status:          util.RandomString(4),
		SalesRepName: util.GenFullname(),
		CustomerName:    util.GenFullname(),
		PitchTag:        util.RandomString(3),
		CustomerRequest: util.RandomString(6),
		RequestDeadline: time.Now().UTC().Add(42 * time.Hour),
	}
	newPitch, err := ts.CreatePitchRequest(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, newPitch)
	require.NotEmpty(t, newPitch.ID)
	require.NotEmpty(t, newPitch.CreatedAt)
	require.Empty(t, newPitch.UpdatedAt)
	require.False(t, newPitch.AdminViewed)
	require.Equal(t, args.SalesRepID, newPitch.SalesRepID)
	require.Equal(t, args.SalesRepName, newPitch.SalesRepName)
	require.WithinDuration(t, args.RequestDeadline, newPitch.RequestDeadline, time.Second)
	return newPitch
}

func TestCreatePitchRequest(t *testing.T) {
	userRole := []UsrRole{{role: "salesRep"}}
	salesRep := createNewUser(t, userRole)[0]
	require.NotEmpty(t, salesRep)
	createPitch(t, salesRep.ID)
}
