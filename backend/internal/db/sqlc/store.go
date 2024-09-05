package db

import (
	// "context"
	"context"
	"database/sql"
	"fmt"
	"time"
	// "fmt"
)

// Interface that implements quiries and Transactions
type Store interface {
	Querier
	UpdateUserTxn(ctx context.Context, args UpdateUsrTxnArgs) error
}

// Extended single queries of Queries struct to enable transactions
type SqlStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SqlStore{
		Queries: New(db),
		db:      db,
	}
}

func (store *SqlStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	txQuery  := New(tx)
	txErr := fn(txQuery)
	if txErr != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("txErr: %v\nrbErr: %v", txErr, rbErr)
		}
		return txErr
	}
	return tx.Commit()
}

// SetNullPitchID converts user pitch_id into sql format
func SetNullPitchID(pitchID int64) sql.NullInt64 {
	if pitchID > 0 {
		return sql.NullInt64{
			Int64: pitchID,
			Valid: true,
		}
	}
	return sql.NullInt64{
		Int64: 0,
		Valid: false,
	}
}



type UpdateUsrTxnArgs struct {
	ID        int64  `json:"user_id" binding:"numeric,gt=0"`
	Fullname  string `json:"fullname" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Username  string `json:"username" binding:"required,alphanum"`
}

// UpdateUserTxn: takes a user details, id, username, fullname and email, 
// updates pitch request, deals. If a fail occur at any point the transaction
// is rolled back.
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