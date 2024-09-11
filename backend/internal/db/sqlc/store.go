package db

import (
	// "context"
	"context"
	"database/sql"
	"fmt"
)

// Interface that implements quiries and Transactions
type Store interface {
	Querier
	UpdateUserTxn(ctx context.Context, args UpdateUsrTxnArgs) error
	ResetPasswordTxn(ctx context.Context, args ResetPasswordArgs) error
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