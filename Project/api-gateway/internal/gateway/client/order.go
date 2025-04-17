package client

import (
	orderpb "github.com/quemin2402/order-service/proto"
	"google.golang.org/grpc"
)

func NewOrder(addr string) (orderpb.OrderServiceClient, *grpc.ClientConn, error) {
	cc, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}
	return orderpb.NewOrderServiceClient(cc), cc, nil
}
