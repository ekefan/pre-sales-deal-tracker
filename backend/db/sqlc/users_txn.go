package db

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

type UpdateUserTxParams struct {
	UpdateUserParams
	OldFullName string `json:"old_full_name"`
}

// UpdateUserTx a database transaction that updates users
func (s *SqlStore) UpdateUserTx(ctx context.Context, args UpdateUserTxParams) error {
	return s.execTx(ctx, func(q *Queries) error {
		newUsr, err := q.UpdateUser(ctx, args.UpdateUserParams)
		if err != nil {
			return fmt.Errorf("failed to update user: %v", err)
		}
		if newUsr < 1 {
			return fmt.Errorf("no rows affected after update")
		}
		if deal, err := q.GetDealToUpdateSalesName(ctx, args.OldFullName); err == nil {
			fmt.Println(deal)
			err = q.UpdateDealSalesName(ctx, UpdateDealSalesNameParams{
				NewSalesName: args.FullName,
				OldSalesName: args.OldFullName,
			})
			if err != nil {
				return fmt.Errorf("failed to update deals sales rep name: %v", err)
			}
		} else if errors.Is(err, pgx.ErrNoRows) {
			slog.Info("did not update any deal for", "user", args.OldFullName,)
			return nil
		} else {
			return fmt.Errorf("failed to get deals to update sales rep name")
		}
		return nil
	})
}
