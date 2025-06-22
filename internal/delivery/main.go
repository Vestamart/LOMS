package delivery

import (
	"context"
	"errors"

	"github.com/vestamart/loms/internal/app/loms"
	"github.com/vestamart/loms/internal/localErr"
	desc "github.com/vestamart/loms/pkg/api/loms/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	desc.UnimplementedLomsServer
	Service loms.Service
}

func NewServer(service loms.Service) *Server {
	return &Server{Service: service}
}

func validateOrderCreateRequest(req *desc.OrderCreateRequest) error {
	if req.User <= 0 {
		return errors.New("user id must be positive")
	}
	if len(req.Items) == 0 {
		return errors.New("items must not be empty")
	}
	for _, item := range req.Items {
		if item.Sku == 0 {
			return errors.New("item sku must be positive")
		}
		if item.Count == 0 {
			return errors.New("item count must be positive")
		}
	}
	return nil
}

func validateOrderId(id int64) error {
	if id <= 0 {
		return errors.New("order id must be positive")
	}
	return nil
}

func validateSku(sku uint32) error {
	if sku == 0 {
		return errors.New("sku must be positive")
	}
	return nil
}

func (s Server) OrderCreate(ctx context.Context, request *desc.OrderCreateRequest) (*desc.OrderCreateResponse, error) {
	ops := "Server OrderCreate"
	if request == nil || (request.User == 0 || len(request.Items) == 0) {
		return nil, status.Errorf(codes.InvalidArgument, "%s empty order create request", ops)
	}

	err := validateOrderCreateRequest(request)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s: %v", ops, err)
	}

	resp, err := s.Service.OrderCreate(ctx, request)
	if err != nil {
		if errors.Is(err, localErr.SKUNotExistErr) {
			return nil, status.Errorf(codes.NotFound, "%s: %v", ops, err)
		}
		if errors.Is(err, localErr.ItemNotEnoughErr) {
			return nil, status.Errorf(codes.ResourceExhausted, "%s: %v", ops, err)
		}
		return nil, status.Errorf(codes.Internal, "%s: %v", ops, err)
	}

	return resp, status.Error(codes.OK, "")
}

func (s Server) OrderInfo(ctx context.Context, request *desc.OrderInfoRequest) (*desc.OrderInfoResponse, error) {
	ops := "Server OrderInfo"

	if err := validateOrderId(request.OrderId); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s: %v", ops, err)
	}

	resp, err := s.Service.OrderInfo(ctx, request)
	if err != nil {
		if errors.Is(err, localErr.OrderNotFoundErr) {
			return nil, status.Errorf(codes.NotFound, "%s: %v", ops, err)
		}
		return nil, status.Errorf(codes.Internal, "%s: %v ", ops, err)
	}

	return resp, status.Error(codes.OK, "")
}

func (s Server) OrderPay(ctx context.Context, request *desc.OrderPayRequest) (*desc.OrderPayResponse, error) {
	ops := "Server OrderPay"

	if err := validateOrderId(request.OrderID); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s: %v", ops, err)
	}

	resp, err := s.Service.OrderPay(ctx, request)
	if err != nil {
		if errors.Is(err, localErr.OrderNotFoundErr) {
			return nil, status.Errorf(codes.NotFound, "%s: %v", ops, err)
		}
		return nil, status.Errorf(codes.Internal, "%s: %v", ops, err)
	}

	return resp, nil
}

func (s Server) OrderCancel(ctx context.Context, request *desc.OrderCancelRequest) (*desc.OrderCancelResponse, error) {
	ops := "Server OrderCancel"

	if err := validateOrderId(request.OrderID); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s: %v", ops, err)
	}

	resp, err := s.Service.OrderCancel(ctx, request)
	if err != nil {
		if errors.Is(err, localErr.OrderNotFoundErr) {
			return nil, status.Errorf(codes.NotFound, "%s: %v ", ops, err)
		}
		return nil, status.Errorf(codes.Internal, "%s: %w", ops, err)
	}

	return resp, status.Error(codes.OK, "")
}

func (s Server) StocksInfo(ctx context.Context, request *desc.StocksInfoRequest) (*desc.StocksInfoResponse, error) {
	ops := "Server StocksInfo"

	if err := validateSku(request.Sku); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s: %v", ops, err)
	}

	resp, err := s.Service.StocksInfo(ctx, request)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "%s: %v", ops, err)
	}

	return resp, status.Error(codes.OK, "")
}
