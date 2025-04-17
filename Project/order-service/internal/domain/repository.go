package domain

import "context"

type OrderRepository interface {
	Create(ctx context.Context, o *Order) error
	Get(ctx context.Context, id string) (*Order, error)
	Update(ctx context.Context, o *Order) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*Order, error)
}
