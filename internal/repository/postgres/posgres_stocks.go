package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/vestamart/loms/internal/localErr"
)

func NewStocksRepositoryPostgres(conn *pgx.Conn) *StocksRepositoryPostgres {
	return &StocksRepositoryPostgres{conn: conn}
}

type StocksRepositoryPostgres struct {
	conn *pgx.Conn
}

func (s StocksRepositoryPostgres) Reserve(ctx context.Context, sku uint32, count uint32) error {
	err := pgx.BeginFunc(ctx, s.conn, func(tx pgx.Tx) (err error) {
		internalRepository := New(tx)
		resp, err := internalRepository.GetBySKIStocks(ctx, int32(sku))
		if err != nil {
			return fmt.Errorf("failed to get reserved stocks: %w", err)
		}
		if resp.TotalCount-resp.Reserved < int32(count) {
			return localErr.ItemNotEnoughErr
		}
		err = internalRepository.ReserveStocks(ctx, &ReserveStocksParams{
			Reserved: resp.Reserved + int32(count),
			Sku:      int32(sku),
		})
		if err != nil {
			return fmt.Errorf("failed to reserve stocks: %w", err)
		}

		return nil
	})
	return err
}

func (s StocksRepositoryPostgres) ReserveRemove(ctx context.Context, skus map[uint32]uint32) error {
	err := pgx.BeginFunc(ctx, s.conn, func(tx pgx.Tx) (err error) {
		repository := New(tx)
		for k, v := range skus {
			resp, err := repository.GetBySKIStocks(ctx, int32(k))
			if err != nil {
				return fmt.Errorf("failed to get stocks: %w", err)
			}

			if resp.TotalCount-resp.Reserved < int32(v) {
				return localErr.ItemNotEnoughErr
			}

			err = repository.ReserveRemoveStocks(ctx, &ReserveRemoveStocksParams{
				Reserved:   resp.Reserved + int32(v),
				Sku:        int32(k),
				TotalCount: resp.TotalCount - int32(v),
			})
			if err != nil {
				return fmt.Errorf("failed to reserve stocks: %w", err)
			}
		}
		return nil
	})
	return err
}

func (s StocksRepositoryPostgres) ReserveCancel(ctx context.Context, skus map[uint32]uint32) error {
	err := pgx.BeginFunc(ctx, s.conn, func(tx pgx.Tx) (err error) {
		repository := New(tx)
		for k, v := range skus {
			resp, err := repository.GetBySKIStocks(ctx, int32(k))
			if err != nil {
				return fmt.Errorf("failed to get stocks: %w", err)
			}
			err = repository.ReserveCancelStocks(ctx, &ReserveCancelStocksParams{
				Reserved: resp.Reserved + int32(v),
				Sku:      int32(k),
			})
			if err != nil {
				return fmt.Errorf("failed to reserve stocks: %w", err)
			}
		}
		return nil
	})
	return err
}

func (s StocksRepositoryPostgres) GetBySKU(ctx context.Context, sku uint32) (uint32, uint32, error) {

	internalRepository := New(s.conn)
	resp, err := internalRepository.GetBySKIStocks(ctx, int32(sku))
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get stocks: %w", err)
	}

	return uint32(resp.TotalCount), uint32(resp.Reserved), nil
}

func (s StocksRepositoryPostgres) RollbackReserve(ctx context.Context, skus map[uint32]uint32) error {
	return nil
}
