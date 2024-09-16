package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

// createRandomUserArg generates the params need for creating a user
func createRandomUserArg(role string) CreateNewUserParams {
	return CreateNewUserParams{
		Username: RandomString(5),
		Role:     role,
		FullName: GenFullname(),
		Email:    GenEmail(),
		Password: GenPassWord(),
	}
}

// createNewUser loops through UsrRole and creates a user for each role
func createNewUsr(t *testing.T, createUsers []UsrRole) []User {
	users := []User{}
	for _, user := range createUsers {
		t.Run(user.role, func(t *testing.T) {
			arg := createRandomUserArg(user.role)
			require.NotEmpty(t, arg)
			newUsr, err := ts.CreateNewUser(context.Background(), arg)
			require.NoError(t, err)
			require.NotEmpty(t, newUsr)
			require.NotEmpty(t, newUsr.ID)
			require.NotEmpty(t, newUsr.CreatedAt)
			require.False(t, newUsr.CreatedAt.IsZero())
			require.NotEmpty(t, newUsr.UpdatedAt)
			require.True(t, newUsr.UpdatedAt.IsZero())
			require.Equal(t, newUsr.Role, arg.Role)
			require.Equal(t, newUsr.Password, arg.Password)
			require.Equal(t, newUsr.PasswordChanged, false)
			require.Equal(t, newUsr.Email, arg.Email)
			require.Equal(t, newUsr.FullName, arg.FullName)
			users = append(users, newUsr)
		})
	}
	return users
}

// UsrRole holds the role field for creating users based on role
type UsrRole struct {
	role string
}

func TestCreateNewUser(t *testing.T) {
	createUsers := []UsrRole{
		{role: "admin"},
		{role: "manager"},
		{role: "sales"},
	}

	// Create a new User for each role
	createNewUsr(t, createUsers)
}

func TestCreateDeal(t *testing.T) {
	// create a randomUser with role sales_rep
	salesRep := createNewUsr(t, []UsrRole{
		{role: "sales"},
	})[0]
	require.NotEmpty(t, salesRep)
	// create a new pitchrequest with the sales_rep id
	pitchReq := createPitch(t, salesRep.ID, salesRep.FullName)
	require.NotEmpty(t, pitchReq)
	// Create a deal based on the pitch request and sales_rep
	args := CreateDealParams{
		PitchID:             SetNullPitchID(pitchReq.ID),
		SalesRepName:        salesRep.FullName,
		CustomerName:        pitchReq.CustomerName,
		ServiceToRender:     pitchReq.CustomerRequest,
		Status:              "ongoing",
		StatusTag:           "presales",
		CurrentPitchRequest: pitchReq.PitchTag,
	}

	deal, err := ts.CreateDeal(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, deal)
	require.NotEmpty(t, deal.ID)
	require.Equal(t, deal.SalesRepName, args.SalesRepName)
	require.Equal(t, deal.CustomerName, args.CustomerName)
	require.NotEmpty(t, deal.ServiceToRender, args.ServiceToRender)
	require.Equal(t, len(deal.ServiceToRender), len(args.ServiceToRender))
	require.NotEmpty(t, deal.Profit)
	require.Equal(t, deal.Profit, "0.00")
	require.NotEmpty(t, deal.NetTotalCost)
	require.Equal(t, deal.NetTotalCost, "0.00")
	require.True(t, deal.UpdatedAt.IsZero())
	require.True(t, deal.ClosedAt.IsZero())
	require.NotEmpty(t, deal.CreatedAt)
	require.False(t, deal.CreatedAt.IsZero())
	require.False(t, deal.Awarded)
}


// NOTE TO MYSELF: only system tests should live in a dedicated test folder