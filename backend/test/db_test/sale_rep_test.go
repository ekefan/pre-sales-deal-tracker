package db_test

import (
	"context"
	"testing"
	"time"

	db "github.com/ekefan/deal-tracker/internal/db/sqlc"
	util "github.com/ekefan/deal-tracker/test/util"
	"github.com/stretchr/testify/require"
)

func createPitch(t *testing.T, salesRepId int64, salesRepFullname string) db.PitchRequest {
	
	customerRequests := []string{util.RandomString(6), util.RandomString(5)}
	args := db.CreatePitchRequestParams{
		SalesRepID:      salesRepId,
		Status:          util.RandomString(4),
		SalesRepName: salesRepFullname,
		CustomerName:    util.GenFullname(),
		PitchTag:        util.RandomString(3),
		CustomerRequest: customerRequests,
		RequestDeadline: time.Now().UTC().Add(42 * time.Hour),
	}
	newPitch, err := ts.CreatePitchRequest(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, newPitch)
	require.NotEmpty(t, newPitch.ID)
	require.NotEmpty(t, newPitch.CreatedAt)
	require.NotEmpty(t, newPitch.UpdatedAt)
	require.True(t, newPitch.UpdatedAt.IsZero())
	require.Equal(t, newPitch.AdminViewed, false)
	require.Equal(t, args.SalesRepID, newPitch.SalesRepID)
	require.Equal(t, args.SalesRepName, newPitch.SalesRepName)
	require.WithinDuration(t, args.RequestDeadline, newPitch.RequestDeadline, time.Second)
	require.NotEmpty(t, newPitch.CustomerRequest)
	require.Equal(t, len(newPitch.CustomerRequest), len(customerRequests))
	return newPitch
}

func TestCreatePitchRequest(t *testing.T) {
	userRole := []UsrRole{{role: "sales"}}
	salesRep := createNewUser(t, userRole)[0]
	require.NotEmpty(t, salesRep)
	createPitch(t, salesRep.ID, salesRep.FullName)
}
