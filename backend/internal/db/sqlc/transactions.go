package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ekefan/deal-tracker/internal/utils"
)

//UpdateUsrTxnArgs: holds fields required to complete an UpdateUserTxn
type UpdateUsrTxnArgs struct {
	ID        int64  `json:"user_id" binding:"numeric,gt=0"`
	Fullname  string `json:"fullname" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Username  string `json:"username" binding:"required,alphanum"`
}

// UpdateUserTxn: takes a user id, username, fullname and email, 
// updates sales-rep name on pitch request and deals tables. 
// If an error occurs at any point the transaction is rolled back.
func (s *SqlStore) UpdateUserTxn(ctx context.Context, args UpdateUsrTxnArgs) error {
	return s.execTx(ctx, func(q *Queries) error {
		userToUpdate, err := q.GetUserForUpdate(ctx, args.ID)
		if err != nil {
			return fmt.Errorf("can't get user for update: %s", err)
		}

		// prepare args for updating sales-rep name in pitch request
		prArgs := UpdatePitchRequestUserNameParams {
			SalesRepID: userToUpdate.ID,
			SalesRepName: args.Fullname,
		}

		//update pitch request to reflect new name where sales_rep_id = user_id //UpdatePitchRequestUserName, 
		err = q.UpdatePitchRequestUserName(ctx, prArgs)
		if err != nil {
			return fmt.Errorf("can't update name in pitch requests: %s", err)
		}
		//update deals to reflext new name where sales_rep_name == fullname

		dealArg := UpdateDealUserNameParams {
			SalesRepName: userToUpdate.FullName,
			SalesRepName_2: args.Fullname,
		}
		err = q.UpdateDealUserName(ctx, dealArg)
		 if err != nil {
			return fmt.Errorf("can't update name in deals: %s", err)
		 }
		// update user to reflex new name and details...
		userArg := AdminUpdateUserParams {
			ID: args.ID,
			FullName: args.Fullname,
			Email: args.Email,
			Username: args.Username,
			UpdatedAt: time.Now(),
		}
		_, err = q.AdminUpdateUser(ctx, userArg)
		if err != nil {
			return fmt.Errorf("couldn't update user: %s", err)
		}
		return nil
	})
}

// ResetPasswordArgs holds the field user_id required to reset user password
type ResetPasswordArgs struct {
	UserID int64 `json:"user_id"`
}

// ResetPasswordTxn: gets a user with the user_id in ResetPasswordArgs
// Gets the user for update, hash the default password, then updates 
// the users password to the default password hash
func (s *SqlStore)ResetPasswordTxn(ctx context.Context, args ResetPasswordArgs) error {
	return s.execTx(ctx, func (q *Queries) error {
		userToUpdate, err := q.GetUserForUpdate(ctx, args.UserID)
		if err != nil {
			return fmt.Errorf("can't get user for update: %s", err)
		}
		hashedDefaultPassword, err := utils.HashPassword(utils.DefaultPassword)
		if err != nil {
			return fmt.Errorf("couldn't hash default password: %s", err)
		}
		args := UpdatePassWordParams{
			ID: userToUpdate.ID,
			Password: hashedDefaultPassword,
			PasswordChanged: !userToUpdate.PasswordChanged,
			UpdatedAt: time.Now(),
		}
		err = q.UpdatePassWord(ctx, args)
		if err != nil {
			return fmt.Errorf("can't reset userPassword: %s", err)
		}
		return nil
	})
}

// CreateDealTxnArgs holds the field for createDealTxn
type CreateDealTxnArgs struct {
	PitchID int64 `json:"pitch_id"`
}


// CreateDealTxn receives the pitch_id from args and creates a deal based on the request
// On failure deal is not created
func (s *SqlStore)CreateDealTxn(ctx context.Context, args CreateDealTxnArgs) error {
	return s.execTx(ctx, func (q *Queries) error {

		pitchReq, err := q.GetPitchRequestByID(ctx, args.PitchID)
		if err != nil {
			return fmt.Errorf("couldn't create deal from pitch request: %s", err)
		}

		requestID := sql.NullInt64{
			Int64: pitchReq.ID,
			Valid: true,
		}
		dealArgs := CreateDealParams{
			PitchID: requestID,
			SalesRepName: pitchReq.SalesRepName,
			CustomerName: pitchReq.CustomerName,
			ServiceToRender: pitchReq.CustomerRequest,
			Status: utils.DefaultStatus,
			StatusTag: utils.DefaultStatusTag,
			CurrentPitchRequest: pitchReq.PitchTag,
		}

		_, err = q.CreateDeal(ctx, dealArgs)
		if err != nil {
			return fmt.Errorf("couldn't create deal from pitch_reques %s", err)
		}
		return nil
	})
}
