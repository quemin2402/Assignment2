package transport

import (
	"context"
	"github.com/quemin2402/inventory-service/internal/domain"
	"github.com/quemin2402/inventory-service/internal/usecase"
	"github.com/quemin2402/inventory-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type srv struct {
	proto.UnimplementedInventoryServiceServer
	uc usecase.ProductUC
}

func NewServer(uc usecase.ProductUC) *grpc.Server {
	s := grpc.NewServer()
	proto.RegisterInventoryServiceServer(s, &srv{uc: uc})
	return s
}

func (s *srv) CreateProduct(ctx context.Context, req *proto.CreateProductRequest) (*proto.ProductResponse, error) {
	if err := s.uc.Create(ctx, toDomain(req.Product)); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &proto.ProductResponse{Product: req.Product}, nil
}

func (s *srv) GetProduct(ctx context.Context, id *proto.ProductID) (*proto.ProductResponse, error) {
	p, err := s.uc.Get(ctx, id.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &proto.ProductResponse{Product: toProto(p)}, nil
}

func (s *srv) UpdateProduct(ctx context.Context, req *proto.UpdateProductRequest) (*proto.ProductResponse, error) {
	if err := s.uc.Update(ctx, toDomain(req.Product)); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &proto.ProductResponse{Product: req.Product}, nil
}

func (s *srv) DeleteProduct(ctx context.Context, id *proto.ProductID) (*emptypb.Empty, error) {
	if err := s.uc.Delete(ctx, id.Id); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *srv) ListProducts(_ *proto.ListProductsRequest, stream proto.InventoryService_ListProductsServer) error {
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

func toDomain(p *proto.Product) *domain.Product {
	return &domain.Product{
		ID: p.Id, Name: p.Name, Category: p.Category, Price: p.Price, Stock: p.Stock,
	}
}
func toProto(p *domain.Product) *proto.Product {
	return &proto.Product{
		Id: p.ID, Name: p.Name, Category: p.Category, Price: p.Price, Stock: p.Stock,
	}
}
