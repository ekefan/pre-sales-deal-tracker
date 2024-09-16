package db

import (
	"context"
	"testing"
	"time"
	"github.com/stretchr/testify/require"
)

func createPitch(t *testing.T, salesRepId int64, salesRepFullname string) PitchRequest {
	
	customerRequests := []string{RandomString(6), RandomString(5)}
	args := CreatePitchRequestParams{
		SalesRepID:      salesRepId,
		Status:          RandomString(4),
		SalesRepName: salesRepFullname,
		CustomerName:    GenFullname(),
		PitchTag:        RandomString(3),
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
	salesRep := createNewUsr(t, userRole)[0]
	require.NotEmpty(t, salesRep)
	createPitch(t, salesRep.ID, salesRep.FullName)
}
