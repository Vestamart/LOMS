package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/vestamart/loms/internal/domain"
	"log"
)

type OrderRepositoryPostgres struct {
	conn *pgx.Conn
}

func NewOrderRepositoryPostgres(conn *pgx.Conn) *OrderRepositoryPostgres {
	return &OrderRepositoryPostgres{conn: conn}
}

func (r OrderRepositoryPostgres) Create(ctx context.Context, userID int64, items *[]domain.Item) (int64, error) {
	internalRepository := New(r.conn)
	orderID, err := internalRepository.InsertOrder(ctx, &InsertOrderParams{
		UserID: userID,
		Status: 0,
	})
	if err != nil {
		return 0, fmt.Errorf("create order failed: %w", err)
	}

	for _, item := range *items {
		itemID, err := internalRepository.InsertItems(ctx, &InsertItemsParams{
			Sku:   int32(item.Sku),
			Count: int32(item.Count),
		})
		if err != nil {
			return 0, fmt.Errorf("insert items failed: %w", err)
		}

		err = internalRepository.InsertOrderItems(ctx, &InsertOrderItemsParams{
			OrderID: orderID,
			ItemID:  itemID,
		})
		if err != nil {
			return 0, fmt.Errorf("insert order items failed : %w", err)
		}
	}

	return orderID, nil

}

func (r OrderRepositoryPostgres) SetStatus(ctx context.Context, orderID int64, status domain.OrderStatus) error {
	internalRepository := New(r.conn)
	err := internalRepository.UpdateStatusOrders(ctx, &UpdateStatusOrdersParams{
		Status:  int16(status),
		OrderID: orderID,
	})

	if err != nil {
		return fmt.Errorf("update status failed: %w", err)
	}

	return nil
}

func (r OrderRepositoryPostgres) GetByID(ctx context.Context, orderID int64) (*domain.Order, error) {
	internalRepository := New(r.conn)
	resp, err := internalRepository.GetInfoFromOrders(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("get info from order failed: %w", err)
	}
	log.Printf("%T - %v\n", resp.Items, string(resp.Items))
	var items []domain.Item
	if err = json.Unmarshal(resp.Items, &items); err != nil {
		return nil, fmt.Errorf("unmarshal items failed: %w", err)
	}

	response := domain.Order{
		UserID: resp.UserID,
		Status: domain.OrderStatus(resp.Status),
		Items:  items,
	}

	return &response, nil
}
