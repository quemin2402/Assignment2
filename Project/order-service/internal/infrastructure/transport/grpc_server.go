package transport

import (
	"Assignment2/Project/order-service/internal/domain"
	"Assignment2/Project/order-service/internal/usecase"
	"context"
	pb "github.com/quemin2402/order-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type srv struct {
	pb.UnimplementedOrderServiceServer
	uc usecase.OrderUC
}

func NewServer(uc usecase.OrderUC) *grpc.Server {
	s := grpc.NewServer()
	pb.RegisterOrderServiceServer(s, &srv{uc})
	return s
}

func (s *srv) CreateOrder(ctx context.Context, r *pb.CreateOrderRequest) (*pb.OrderResponse, error) {
	if err := s.uc.Create(ctx, toDomain(r.Order)); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &pb.OrderResponse{Order: r.Order}, nil
}
func (s *srv) GetOrder(ctx context.Context, id *pb.OrderID) (*pb.OrderResponse, error) {
	o, err := s.uc.Get(ctx, id.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &pb.OrderResponse{Order: toProto(o)}, nil
}
func (s *srv) UpdateOrder(ctx context.Context, r *pb.UpdateOrderRequest) (*pb.OrderResponse, error) {
	if err := s.uc.Update(ctx, toDomain(r.Order)); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &pb.OrderResponse{Order: r.Order}, nil
}
func (s *srv) ListOrders(ctx context.Context, _ *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	list, err := s.uc.List(ctx)
	if err != nil {
		return nil, err
	}
	var resp pb.ListOrdersResponse
	for _, o := range list {
		resp.Orders = append(resp.Orders, toProto(o))
	}
	return &resp, nil
}
func (s *srv) DeleteOrder(ctx context.Context, id *pb.OrderID) (*emptypb.Empty, error) {
	if err := s.uc.Delete(ctx, id.Id); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &emptypb.Empty{}, nil
}

/* mapping */
func toDomain(o *pb.Order) *domain.Order {
	var items []domain.OrderItem
	for _, it := range o.Items {
		items = append(items, domain.OrderItem{ProductID: it.ProductId, Quantity: it.Quantity})
	}
	return &domain.Order{ID: o.Id, Status: o.Status, Items: items}
}
func toProto(o *domain.Order) *pb.Order {
	var items []*pb.OrderItem
	for _, it := range o.Items {
		items = append(items, &pb.OrderItem{ProductId: it.ProductID, Quantity: it.Quantity})
	}
	return &pb.Order{Id: o.ID, Status: o.Status, Items: items}
}
