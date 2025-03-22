package loms

import (
	"context"
	"errors"
	"fmt"
	"github.com/vestamart/loms/internal/domain"
	"github.com/vestamart/loms/internal/localErr"
	desc "github.com/vestamart/loms/pkg/api/loms/v1"
)

// OrdersRepository и StocksStorage интерфейсы для взаимодействия с репозиториями
//
//go:generate minimock -i github.com/vestamart/loms/internal/app/loms.OrdersRepository -o ./mock/orders_repository_mock.go -n OrdersRepositoryMock -p mock
type OrdersRepository interface {
	Create(_ context.Context, userID int64, items *[]domain.Item) (int64, error)
	SetStatus(_ context.Context, orderID int64, status domain.OrderStatus) error
	GetByID(_ context.Context, orderID int64) (*domain.Order, error)
}

//go:generate minimock -i github.com/vestamart/loms/internal/app/loms.StocksStorage -o ./mock/stock_repository_mock.go -n StocksStorageMock -p mock
type StocksStorage interface {
	Reserve(_ context.Context, sku uint32, count uint32) error
	ReserveRemove(_ context.Context, sku uint32, count uint32) error
	ReserveCancel(_ context.Context, skus map[uint32]uint32) error
	GetBySKU(_ context.Context, sku uint32) (uint32, error)
	RollbackReserve(_ context.Context, skus map[uint32]uint32) error
}

type Service struct {
	ordersRepository OrdersRepository
	stocksRepository StocksStorage
}

func NewService(ordersRepository OrdersRepository, stocksRepository StocksStorage) *Service {
	return &Service{ordersRepository: ordersRepository, stocksRepository: stocksRepository}
}

func (s Service) OrderCreate(ctx context.Context, request *desc.OrderCreateRequest) (*desc.OrderCreateResponse, error) {
	items := make([]domain.Item, 0, len(request.Items))
	for _, v := range request.Items {
		items = append(items, domain.Item{
			Sku:   v.Sku,
			Count: v.Count,
		})
	}

	orderId, err := s.ordersRepository.Create(ctx, request.User, &items)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	var reservedSKUs = make(map[uint32]uint32)

	for _, v := range items {
		err = s.stocksRepository.Reserve(ctx, v.Sku, v.Count)
		if err != nil {
			if errRollback := s.stocksRepository.RollbackReserve(ctx, reservedSKUs); errRollback != nil {
				return nil, fmt.Errorf("failed to rollback reserved items: %w", errRollback)
			}
			if errStatus := s.ordersRepository.SetStatus(ctx, orderId, domain.Failed); errStatus != nil {
				return nil, fmt.Errorf("failed to set status: %w", errStatus)
			}
			return nil, fmt.Errorf("failed to reverse: %w", err)
		}
		reservedSKUs[v.Sku] = v.Count
	}
	if err = s.ordersRepository.SetStatus(ctx, orderId, domain.AwwaitingPayment); err != nil {
		return nil, fmt.Errorf("failed to set status: %w", err)
	}

	return &desc.OrderCreateResponse{OrderId: orderId}, nil
}

func (s Service) OrderInfo(ctx context.Context, request *desc.OrderInfoRequest) (*desc.OrderInfoResponse, error) {
	rawResponse, err := s.ordersRepository.GetByID(ctx, request.OrderId)
	if err != nil {
		if errors.Is(err, localErr.OrderNotFoundErr) {
			return nil, fmt.Errorf("failed to get order %w", err)
		}
		return nil, err
	}
	items := make([]*desc.Item, 0, len(rawResponse.Items))
	for _, v := range rawResponse.Items {
		items = append(items, &desc.Item{
			Sku:   v.Sku,
			Count: v.Count,
		})
	}

	response := &desc.OrderInfoResponse{
		Status: desc.OrderStatus(rawResponse.Status),
		User:   rawResponse.UserID,
		Items:  items,
	}

	return response, nil
}

func (s Service) OrderPay(ctx context.Context, request *desc.OrderPayRequest) (*desc.OrderPayResponse, error) {
	getByID, err := s.ordersRepository.GetByID(ctx, request.OrderID)
	if err != nil {
		if errors.Is(err, localErr.OrderNotFoundErr) {
			return nil, fmt.Errorf("failed to get order %w", err)
		}
		return nil, fmt.Errorf("failed to get order %w", err)
	}

	for _, v := range getByID.Items {
		if err = s.stocksRepository.ReserveRemove(ctx, v.Sku, v.Count); err != nil {
			return nil, fmt.Errorf("failed to reserve remove item: %w", err)
		}
	}

	err = s.ordersRepository.SetStatus(ctx, request.OrderID, domain.Payed)
	if err != nil {
		return nil, fmt.Errorf("failed to set status: %w", err)
	}
	return &desc.OrderPayResponse{}, nil
}

func (s Service) OrderCancel(ctx context.Context, request *desc.OrderCancelRequest) (*desc.OrderCancelResponse, error) {
	rawResponse, err := s.ordersRepository.GetByID(ctx, request.OrderID)
	if err != nil {
		if errors.Is(err, localErr.OrderNotFoundErr) {
			return nil, fmt.Errorf("failed to get order %w", err)
		}
		return nil, fmt.Errorf("failed to get order %w", err)
	}

	items := make(map[uint32]uint32)
	for _, v := range rawResponse.Items {
		items[v.Sku] = v.Count
	}

	if err = s.stocksRepository.ReserveCancel(ctx, items); err != nil {
		if errors.Is(err, localErr.OrderNotFoundErr) {
			return nil, fmt.Errorf("failed to reserve cancel %w", err)
		}
		return nil, fmt.Errorf("failed to reserve cancel %w", err)
	}

	if err = s.ordersRepository.SetStatus(ctx, request.OrderID, domain.Cancelled); err != nil {
	}
	return &desc.OrderCancelResponse{}, nil
}

func (s Service) StocksInfo(ctx context.Context, request *desc.StocksInfoRequest) (*desc.StocksInfoResponse, error) {
	v, err := s.stocksRepository.GetBySKU(ctx, request.Sku)
	if err != nil {
		return nil, fmt.Errorf("failed to get stocks %w", err)
	}

	return &desc.StocksInfoResponse{Count: uint64(v)}, nil
}
