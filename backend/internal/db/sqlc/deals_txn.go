package db

import (
	"context"
	"errors"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	preSalesDepartment = "pre-sales"
)

func (s *SqlStore) CreateDealFromPitchId(ctx context.Context, pitch_id int64) error {
	return s.execTx(ctx, func(q *Queries) error {
		pitchReq, err := q.GetPitchRequestById(ctx, pitch_id)
		if err != nil {
			return pgx.ErrNoRows
		}
		salesRespName, err := q.GetUserFullName(ctx, pitchReq.UserID)
		if err != nil {
			return pgx.ErrNoRows
		}
		netTotal := pgtype.Numeric{}
		profit := pgtype.Numeric{}
		netTotalStr := strconv.FormatFloat(0.0, 'f', -1, 64)
		profitStr := strconv.FormatFloat(0.0, 'f', -1, 64)
		netTotal.Scan(netTotalStr)
		if err := netTotal.Scan(netTotalStr); err != nil {
			return err
		}
		if err := profit.Scan(profitStr); err != nil {
			return err
		}
		dealRow, err := q.CreateDeal(ctx, CreateDealParams{
			PitchID:          &pitchReq.ID,
			SalesRepName:     salesRespName,
			CustomerName:     pitchReq.CustomerName,
			ServicesToRender: pitchReq.CustomerRequest,
			Department:       preSalesDepartment,
			NetTotalCost:     netTotal,
			Profit:           profit,
		})
		if err != nil {
			return err
		}
		if dealRow != 1 {
			return errors.New("couldn't create deal")
		}
		return nil
	})
}
