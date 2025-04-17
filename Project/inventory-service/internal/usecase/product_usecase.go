package usecase

import (
	"context"
	"errors"
	"github.com/quemin2402/inventory-service/internal/domain"
)

type ProductUC interface {
	Create(context.Context, *domain.Product) error
	Get(context.Context, string) (*domain.Product, error)
	Update(context.Context, *domain.Product) error
	Delete(context.Context, string) error
	List(context.Context) ([]*domain.Product, error)
}

type productUC struct{ repo domain.ProductRepository }

func New(repo domain.ProductRepository) ProductUC { return &productUC{repo} }

func (u *productUC) validate(p *domain.Product) error {
	if p.ID == "" || p.Name == "" || p.Price < 0 {
		return errors.New("invalid product")
	}
	return nil
}

func (u *productUC) Create(ctx context.Context, p *domain.Product) error {
	if err := u.validate(p); err != nil {
		return err
	}
	return u.repo.Create(ctx, p)
}

func (u *productUC) Get(ctx context.Context, id string) (*domain.Product, error) {
	return u.repo.Get(ctx, id)
}

func (u *productUC) Update(ctx context.Context, p *domain.Product) error {
	if err := u.validate(p); err != nil {
		return err
	}
	return u.repo.Update(ctx, p)
}

func (u *productUC) Delete(ctx context.Context, id string) error { return u.repo.Delete(ctx, id) }

func (u *productUC) List(ctx context.Context) ([]*domain.Product, error) { return u.repo.List(ctx) }
