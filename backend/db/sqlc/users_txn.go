package db

import (
	"context"
	"fmt"
)

// UpdateUserTx a database transaction that updates users
func (s *SqlStore) UpdateUserTx(ctx context.Context, args UpdateUserParams) error {
	return s.execTx(ctx, func(q *Queries) error {
		err := q.UpdateUser(ctx, args)
		if err != nil {
			return fmt.Errorf("failed to update user: %v", err)
		}
		err = q.UpdateDealSalesName(ctx, args.FullName)
		if err != nil {
			return fmt.Errorf("failed to update deals sales rep name: %v", err)
		}
		return nil
	})
}
