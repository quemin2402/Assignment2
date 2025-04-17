package usecase

import (
	"context"
	"errors"
	"github.com/quemin2402/order-service/internal/domain"
)

type OrderUC interface {
	Create(context.Context, *domain.Order) error
	Get(context.Context, string) (*domain.Order, error)
	Update(context.Context, *domain.Order) error
	Delete(context.Context, string) error
	List(context.Context) ([]*domain.Order, error)
}

type uc struct{ repo domain.OrderRepository }

func New(r domain.OrderRepository) OrderUC { return &uc{r} }

func (u *uc) validate(o *domain.Order) error {
	if len(o.Items) == 0 {
		return errors.New("empty order")
	}
	return nil
}

func (u *uc) Create(ctx context.Context, o *domain.Order) error {
	if err := u.validate(o); err != nil {
		return err
	}
	o.Status = "pending"
	return u.repo.Create(ctx, o)
}

func (u *uc) Get(ctx context.Context, id string) (*domain.Order, error) { return u.repo.Get(ctx, id) }

func (u *uc) Update(ctx context.Context, o *domain.Order) error { return u.repo.Update(ctx, o) }

func (u *uc) Delete(ctx context.Context, id string) error { return u.repo.Delete(ctx, id) }

func (u *uc) List(ctx context.Context) ([]*domain.Order, error) { return u.repo.List(ctx) }
