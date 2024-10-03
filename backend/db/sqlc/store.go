package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	DefaultUserPassword = "vasDeal45"
)

// Store holds function definitions for interacting with the database
// executes queries and transactions
type Store interface {
	Querier
	UpdateUserTx(ctx context.Context, args UpdateUserParams) error
}

// SqlStore holds fields required to interact with database
type SqlStore struct {
	connPool *pgxpool.Pool
	*Queries
}

// NewStore creates a new Store with dbpool
func NewStore(dbpool *pgxpool.Pool) Store {
	return &SqlStore{
		connPool: dbpool,
		Queries:  New(dbpool),
	}
}

// execTx starts a transactions, executes fn within the transaction and rollsback on err
// return nil if transaction is completed
func (store *SqlStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.connPool.Begin(ctx)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("rollback err: %v, transaction err %v", rbErr, err)
		}
		return err
	}

	return tx.Commit(ctx)
}
