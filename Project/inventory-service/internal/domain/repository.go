package domain

import "context"

type ProductRepository interface {
	Create(ctx context.Context, p *Product) error
	Get(ctx context.Context, id string) (*Product, error)
	Update(ctx context.Context, p *Product) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*Product, error)
}
