package transport

import (
	"context"
	"github.com/google/uuid"
	userpb "github.com/quemin2402/user-service"
	"github.com/quemin2402/user-service/internal/domain"
	"github.com/quemin2402/user-service/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type srv struct {
	userpb.UnimplementedUserServiceServer
	uc usecase.UserUC
}

func NewServer(uc usecase.UserUC) *grpc.Server {
	s := grpc.NewServer()
	userpb.RegisterUserServiceServer(s, &srv{uc: uc})
	return s
}

func (s *srv) RegisterUser(ctx context.Context, r *userpb.UserRequest) (*userpb.UserResponse, error) {
	usr := &domain.User{ID: uuid.New().String(), Username: r.Username, Email: r.Email}
	u, err := s.uc.Register(ctx, usr, r.Password)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &userpb.UserResponse{User: toProto(u)}, nil
}

func (s *srv) AuthenticateUser(ctx context.Context, r *userpb.AuthRequest) (*userpb.AuthResponse, error) {
	token, err := s.uc.Auth(ctx, r.Username, r.Password)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	return &userpb.AuthResponse{Token: token}, nil
}

func (s *srv) GetUserProfile(ctx context.Context, id *userpb.UserID) (*userpb.UserResponse, error) {
	u, err := s.uc.GetProfile(ctx, id.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &userpb.UserResponse{User: toProto(u)}, nil
}

func toProto(u *domain.User) *userpb.User {
	return &userpb.User{Id: u.ID, Username: u.Username, Email: u.Email}
}
