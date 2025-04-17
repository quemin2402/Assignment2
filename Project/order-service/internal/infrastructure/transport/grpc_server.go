package transport

import (
	"context"
	"github.com/quemin2402/order-service/internal/domain"
	"github.com/quemin2402/order-service/internal/usecase"
	orderpb "github.com/quemin2402/order-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type srv struct {
	orderpb.UnimplementedOrderServiceServer
	uc usecase.OrderUC
}

func NewServer(uc usecase.OrderUC) *grpc.Server {
	s := grpc.NewServer()
	orderpb.RegisterOrderServiceServer(s, &srv{uc: uc})
	return s
}

func (s *srv) CreateOrder(ctx context.Context, r *orderpb.CreateOrderRequest) (*orderpb.OrderResponse, error) {
	if err := s.uc.Create(ctx, toDomain(r.Order)); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &orderpb.OrderResponse{Order: r.Order}, nil
}

func (s *srv) GetOrder(ctx context.Context, id *orderpb.OrderID) (*orderpb.OrderResponse, error) {
	o, err := s.uc.Get(ctx, id.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &orderpb.OrderResponse{Order: toProto(o)}, nil
}

func (s *srv) UpdateOrder(ctx context.Context, r *orderpb.UpdateOrderRequest) (*orderpb.OrderResponse, error) {
	if err := s.uc.Update(ctx, toDomain(r.Order)); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &orderpb.OrderResponse{Order: r.Order}, nil
}

func (s *srv) ListOrders(ctx context.Context, _ *orderpb.ListOrdersRequest) (*orderpb.ListOrdersResponse, error) {
	list, err := s.uc.List(ctx)
	if err != nil {
		return nil, err
	}
	var resp orderpb.ListOrdersResponse
	for _, o := range list {
		resp.Orders = append(resp.Orders, toProto(o))
	}
	return &resp, nil
}

func (s *srv) DeleteOrder(ctx context.Context, id *orderpb.OrderID) (*emptypb.Empty, error) {
	if err := s.uc.Delete(ctx, id.Id); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func toDomain(o *orderpb.Order) *domain.Order {
	var items []domain.OrderItem
	for _, it := range o.Items {
		items = append(items, domain.OrderItem{ProductID: it.ProductId, Quantity: it.Quantity})
	}
	return &domain.Order{ID: o.Id, Status: o.Status, Items: items}
}

func toProto(o *domain.Order) *orderpb.Order {
	var items []*orderpb.OrderItem
	for _, it := range o.Items {
		items = append(items, &orderpb.OrderItem{ProductId: it.ProductID, Quantity: it.Quantity})
	}
	return &orderpb.Order{Id: o.ID, Status: o.Status, Items: items}
}
