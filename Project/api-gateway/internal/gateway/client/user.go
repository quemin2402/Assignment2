package client

import (
	userpb "github.com/quemin2402/user-service/proto"
	"google.golang.org/grpc"
)

func NewUser(addr string) (userpb.UserServiceClient, *grpc.ClientConn, error) {
	cc, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}
	return userpb.NewUserServiceClient(cc), cc, nil
}
