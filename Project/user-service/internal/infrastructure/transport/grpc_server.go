package transport

import (
	"context"
	"github.com/google/uuid"
	"github.com/quemin2402/user-service/internal/domain"
	"github.com/quemin2402/user-service/internal/usecase"
	"github.com/quemin2402/user-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type srv struct {
	proto.UnimplementedUserServiceServer
	uc usecase.UserUC
}

func NewServer(uc usecase.UserUC) *grpc.Server {
	s := grpc.NewServer()
	proto.RegisterUserServiceServer(s, &srv{uc: uc})
	return s
}

func (s *srv) RegisterUser(ctx context.Context, r *proto.UserRequest) (*proto.UserResponse, error) {
	usr := &domain.User{ID: uuid.New().String(), Username: r.Username, Email: r.Email}
	u, err := s.uc.Register(ctx, usr, r.Password)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &proto.UserResponse{User: toProto(u)}, nil
}

func (s *srv) AuthenticateUser(ctx context.Context, r *proto.AuthRequest) (*proto.AuthResponse, error) {
	token, err := s.uc.Auth(ctx, r.Username, r.Password)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	return &proto.AuthResponse{Token: token}, nil
}

func (s *srv) GetUserProfile(ctx context.Context, id *proto.UserID) (*proto.UserResponse, error) {
	u, err := s.uc.GetProfile(ctx, id.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &proto.UserResponse{User: toProto(u)}, nil
}

func toProto(u *domain.User) *proto.User {
	return &proto.User{Id: u.ID, Username: u.Username, Email: u.Email}
}
