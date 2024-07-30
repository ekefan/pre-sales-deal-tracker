package test

import (
	"context"
	"testing"

	db "github.com/ekefan/deal-tracker/internal/db/sqlc"
	util "github.com/ekefan/deal-tracker/test/util"
	"github.com/stretchr/testify/require"
)

// createRandomUserArg generates the params need for creating a user
func createRandomUserArg(role string) db.CreateNewUserParams {
	return db.CreateNewUserParams{
		Username: util.RandomString(5),
		Role:     role,
		FullName: util.GenFullname(),
		Email:    util.GenEmail(),
		Password: util.GenPassWord(),
	}
}

// createNewUser loops through UsrRole and creates a user for each role
func createNewUser(t *testing.T, createUsers []UsrRole) []db.User {
	users := []db.User{}
	for _, user := range createUsers {
		t.Run(user.role, func(t *testing.T) {
			arg := createRandomUserArg(user.role)
			require.NotEmpty(t, arg)
			newUsr, err := ts.CreateNewUser(context.Background(), arg)
			require.NoError(t, err)
			require.NotEmpty(t, newUsr)
			require.NotEmpty(t, newUsr.ID)
			require.NotEmpty(t, newUsr.CreatedAt)
			require.Empty(t, newUsr.UpdatedAt)
			require.False(t, newUsr.UpdatedAt.Valid)
			require.Equal(t, newUsr.Role, arg.Role)
			require.Equal(t, newUsr.Password, arg.Password)
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
		{role: "salesrep"},
	}

	// Create a new User for each role
	createNewUser(t, createUsers)
}

func TestCreateDeal(t *testing.T) {
	//create a randomUser with role sales_rep
	salesRep := createNewUser(t, []UsrRole{
		{role: "salesrep"},
	})[0]
	require.NotEmpty(t, salesRep)
	//create a new pitchrequest with the sales_rep id
	pitchReq := createPitch(t, salesRep.ID)
	require.NotEmpty(t, pitchReq)
	//Create a deal based on the pitch request and sales_rep
	args := db.CreateDealParams{
		PitchID:             pitchReq.ID,
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
	require.NotEmpty(t, deal.CreatedAt)
	require.Empty(t, deal.Profit)
	require.Empty(t, deal.UpdatedAt)
	require.Empty(t, deal.NetTotalCost)
	require.Empty(t, deal.ClosedAt)
	require.False(t, deal.Awarded)
}
