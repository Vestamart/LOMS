package loms_test

//
//import (
//	"context"
//	"errors"
//	"github.com/vestamart/loms/internal/app/loms"
//	mock2 "github.com/vestamart/loms/internal/app/loms/mock"
//	"github.com/vestamart/loms/internal/domain"
//	"github.com/vestamart/loms/internal/repository"
//	desc "github.com/vestamart/loms/pkg/api/loms/v1"
//	"google.golang.org/grpc/codes"
//	"google.golang.org/grpc/status"
//	"testing"
//
//	"github.com/gojuno/minimock/v3"
//	"github.com/stretchr/testify/assert"
//)
//
//func TestOrderCreate(t *testing.T) {
//	mc := minimock.NewController(t)
//
//	ordersRepoMock := mock2.NewOrdersRepositoryMock(mc)
//	stocksRepoMock := mock2.NewStocksStorageMock(mc)
//
//	svc := loms.NewServer(ordersRepoMock, stocksRepoMock)
//
//	tests := []struct {
//		name    string
//		request *desc.OrderCreateRequest
//		setup   func()
//		expect  codes.Code
//	}{
//		{
//			name:    "invalid request",
//			request: &desc.OrderCreateRequest{},
//			expect:  codes.InvalidArgument,
//		},
//		{
//			name: "success",
//			request: &desc.OrderCreateRequest{
//				User:  1,
//				Items: []*desc.Item{{Sku: 1, Count: 2}},
//			},
//			setup: func() {
//				ordersRepoMock.CreateMock.Return(1, nil)
//				stocksRepoMock.ReserveMock.Return(nil)
//				ordersRepoMock.SetStatusMock.Return(nil)
//			},
//			expect: codes.OK,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if tt.setup != nil {
//				tt.setup()
//			}
//			_, err := svc.OrderCreate(context.Background(), tt.request)
//			assert.Equal(t, tt.expect, status.Code(err))
//		})
//	}
//}
//
//func TestOrderInfo(t *testing.T) {
//	mc := minimock.NewController(t)
//
//	ordersRepoMock := mock2.NewOrdersRepositoryMock(mc)
//	svc := loms.NewServer(ordersRepoMock, nil)
//
//	tests := []struct {
//		name    string
//		request *desc.OrderInfoRequest
//		setup   func()
//		expect  codes.Code
//	}{
//		{
//			name:    "order not found",
//			request: &desc.OrderInfoRequest{OrderId: 1},
//			setup: func() {
//				ordersRepoMock.GetByIDMock.Return(nil, repository.OrderNotExistErr)
//			},
//			expect: codes.NotFound,
//		},
//		{
//			name:    "success",
//			request: &desc.OrderInfoRequest{OrderId: 1},
//			setup: func() {
//				ordersRepoMock.GetByIDMock.Return(&domain.Order{
//					UserID: 1,
//					Status: 1,
//					Items:  []domain.Item{{Sku: 1, Count: 2}},
//				}, nil)
//			},
//			expect: codes.OK,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if tt.setup != nil {
//				tt.setup()
//			}
//			_, err := svc.OrderInfo(context.Background(), tt.request)
//			assert.Equal(t, tt.expect, status.Code(err))
//		})
//	}
//}
//
//func TestOrderPay(t *testing.T) {
//	mc := minimock.NewController(t)
//
//	ordersRepoMock := mock2.NewOrdersRepositoryMock(mc)
//	stocksRepoMock := mock2.NewStocksStorageMock(mc)
//	svc := loms.NewServer(ordersRepoMock, stocksRepoMock)
//
//	tests := []struct {
//		name    string
//		request *desc.OrderPayRequest
//		setup   func()
//		expect  codes.Code
//	}{
//		{
//			name:    "order not found",
//			request: &desc.OrderPayRequest{OrderID: 1},
//			setup: func() {
//				ordersRepoMock.GetByIDMock.Return(nil, repository.OrderNotExistErr)
//			},
//			expect: codes.NotFound,
//		},
//		{
//			name:    "reserve remove error",
//			request: &desc.OrderPayRequest{OrderID: 1},
//			setup: func() {
//				ordersRepoMock.GetByIDMock.Return(&domain.Order{
//					Items: []domain.Item{{Sku: 1, Count: 2}},
//				}, nil)
//				stocksRepoMock.ReserveRemoveMock.Return(errors.New("reserve remove error"))
//			},
//			expect: codes.Internal,
//		},
//		{
//			name:    "set status error",
//			request: &desc.OrderPayRequest{OrderID: 1},
//			setup: func() {
//				ordersRepoMock.GetByIDMock.Return(&domain.Order{
//					Items: []domain.Item{{Sku: 1, Count: 2}},
//				}, nil)
//				stocksRepoMock.ReserveRemoveMock.Return(nil)
//				ordersRepoMock.SetStatusMock.Return(errors.New("set status error"))
//			},
//			expect: codes.Internal,
//		},
//		{
//			name:    "success",
//			request: &desc.OrderPayRequest{OrderID: 1},
//			setup: func() {
//				ordersRepoMock.GetByIDMock.Return(&domain.Order{
//					Items: []domain.Item{{Sku: 1, Count: 2}},
//				}, nil)
//				stocksRepoMock.ReserveRemoveMock.Return(nil)
//				ordersRepoMock.SetStatusMock.Return(nil)
//			},
//			expect: codes.OK,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if tt.setup != nil {
//				tt.setup()
//			}
//			_, err := svc.OrderPay(context.Background(), tt.request)
//			assert.Equal(t, tt.expect, status.Code(err))
//		})
//	}
//}
//
//func TestOrderCancel(t *testing.T) {
//	mc := minimock.NewController(t)
//
//	ordersRepoMock := mock2.NewOrdersRepositoryMock(mc)
//	stocksRepoMock := mock2.NewStocksStorageMock(mc)
//	svc := loms.NewServer(ordersRepoMock, stocksRepoMock)
//
//	tests := []struct {
//		name    string
//		request *desc.OrderCancelRequest
//		setup   func()
//		expect  codes.Code
//	}{
//		{
//			name:    "order not found",
//			request: &desc.OrderCancelRequest{OrderID: 1},
//			setup: func() {
//				ordersRepoMock.GetByIDMock.Return(nil, repository.OrderNotExistErr)
//			},
//			expect: codes.NotFound,
//		},
//		{
//			name:    "reserve cancel error",
//			request: &desc.OrderCancelRequest{OrderID: 1},
//			setup: func() {
//				ordersRepoMock.GetByIDMock.Return(&domain.Order{
//					Items: []domain.Item{{Sku: 1, Count: 2}},
//				}, nil)
//				stocksRepoMock.ReserveCancelMock.Return(errors.New("reserve cancel error"))
//			},
//			expect: codes.Internal,
//		},
//		{
//			name:    "success (ignoring set status error)",
//			request: &desc.OrderCancelRequest{OrderID: 1},
//			setup: func() {
//				ordersRepoMock.GetByIDMock.Return(&domain.Order{
//					Items: []domain.Item{{Sku: 1, Count: 2}},
//				}, nil)
//				stocksRepoMock.ReserveCancelMock.Return(nil)
//				ordersRepoMock.SetStatusMock.Return(errors.New("set status error"))
//			},
//			expect: codes.OK,
//		},
//		{
//			name:    "success",
//			request: &desc.OrderCancelRequest{OrderID: 1},
//			setup: func() {
//				ordersRepoMock.GetByIDMock.Return(&domain.Order{
//					Items: []domain.Item{{Sku: 1, Count: 2}},
//				}, nil)
//				stocksRepoMock.ReserveCancelMock.Return(nil)
//				ordersRepoMock.SetStatusMock.Return(nil)
//			},
//			expect: codes.OK,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if tt.setup != nil {
//				tt.setup()
//			}
//			_, err := svc.OrderCancel(context.Background(), tt.request)
//			assert.Equal(t, tt.expect, status.Code(err))
//		})
//	}
//}
//
//func TestStocksInfo(t *testing.T) {
//	mc := minimock.NewController(t)
//	defer mc.Finish()
//
//	stocksRepoMock := mock2.NewStocksStorageMock(mc)
//	svc := loms.NewServer(nil, stocksRepoMock)
//
//	tests := []struct {
//		name    string
//		request *desc.StocksInfoRequest
//		setup   func()
//		expect  codes.Code
//		count   uint64
//	}{{
//		name:    "error on get by SKU",
//		request: &desc.StocksInfoRequest{Sku: 1},
//		setup: func() {
//			stocksRepoMock.GetBySKUMock.Return(uint32(0), errors.New("get error"))
//		},
//		expect: codes.NotFound,
//	}, {
//		name:    "success",
//		request: &desc.StocksInfoRequest{Sku: 1},
//		setup: func() {
//			stocksRepoMock.GetBySKUMock.Return(10, nil)
//		},
//		expect: codes.OK,
//		count:  10,
//	}}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if tt.setup != nil {
//				tt.setup()
//			}
//			resp, err := svc.StocksInfo(context.Background(), tt.request)
//			assert.Equal(t, tt.expect, status.Code(err))
//			if err == nil {
//				assert.Equal(t, tt.count, resp.Count)
//			}
//		})
//	}
//}
