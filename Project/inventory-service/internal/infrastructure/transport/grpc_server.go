package transport

import (
	"context"
	pb "github.com/quemin2402/Assignment2/Project/inventory-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"inventory-service/internal/domain"
	"inventory-service/internal/usecase"
)

type srv struct {
	pb.UnimplementedInventoryServiceServer
	uc usecase.ProductUC
}

func NewServer(uc usecase.ProductUC) *grpc.Server {
	s := grpc.NewServer()
	pb.RegisterInventoryServiceServer(s, &srv{uc: uc})
	return s
}

func (s *srv) CreateProduct(ctx context.Context, in *pb.Product) (*pb.Product, error) {
	if err := s.uc.Create(ctx, toDomain(in)); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return in, nil
}

func (s *srv) GetProduct(ctx context.Context, id *pb.ProductID) (*pb.Product, error) {
	p, err := s.uc.Get(ctx, id.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return toProto(p), nil
}

func (s *srv) UpdateProduct(ctx context.Context, in *pb.Product) (*pb.Product, error) {
	if err := s.uc.Update(ctx, toDomain(in)); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return in, nil
}

func (s *srv) DeleteProduct(ctx context.Context, id *pb.ProductID) (*pb.Empty, error) {
	if err := s.uc.Delete(ctx, id.Id); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &pb.Empty{}, nil
}

func (s *srv) ListProducts(_ *pb.Empty, stream pb.InventoryService_ListProductsServer) error {
	list, err := s.uc.List(stream.Context())
	if err != nil {
		return err
	}
	for _, p := range list {
		if err := stream.Send(toProto(p)); err != nil {
			return err
		}
	}
	return nil
}

func toDomain(p *pb.Product) *domain.Product {
	return &domain.Product{
		ID: p.Id, Name: p.Name, Category: p.Category, Price: p.Price, Stock: p.Stock}
}

func toProto(p *domain.Product) *pb.Product {
	return &pb.Product{
		Id: p.ID, Name: p.Name, Category: p.Category, Price: p.Price, Stock: p.Stock}
}
