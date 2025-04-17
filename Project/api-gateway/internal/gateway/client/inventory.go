package client

import (
	invpb "github.com/quemin2402/inventory-service/proto"
	"google.golang.org/grpc"
)

func NewInventory(addr string) (invpb.InventoryServiceClient, *grpc.ClientConn, error) {
	cc, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}
	return invpb.NewInventoryServiceClient(cc), cc, nil
}
