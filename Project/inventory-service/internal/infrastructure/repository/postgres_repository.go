package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"inventory-service/internal/domain"
)

type repo struct{ db *pgxpool.Pool }

func New(db *pgxpool.Pool) domain.ProductRepository { return &repo{db} }

func (r *repo) Create(ctx context.Context, p *domain.Product) error {
	_, err := r.db.Exec(ctx, `INSERT INTO products(id,name,category,price,stock)
		VALUES($1,$2,$3,$4,$5)`, p.ID, p.Name, p.Category, p.Price, p.Stock)
	return err
}

func (r *repo) Get(ctx context.Context, id string) (*domain.Product, error) {
	row := r.db.QueryRow(ctx, `SELECT id,name,category,price,stock FROM products WHERE id=$1`, id)
	var p domain.Product
	err := row.Scan(&p.ID, &p.Name, &p.Category, &p.Price, &p.Stock)
	return &p, err
}

func (r *repo) Update(ctx context.Context, p *domain.Product) error {
	_, err := r.db.Exec(ctx, `UPDATE products SET name=$2,category=$3,price=$4,stock=$5 WHERE id=$1`,
		p.ID, p.Name, p.Category, p.Price, p.Stock)
	return err
}

func (r *repo) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM products WHERE id=$1`, id)
	return err
}

func (r *repo) List(ctx context.Context) ([]*domain.Product, error) {
	rows, err := r.db.Query(ctx, `SELECT id,name,category,price,stock FROM products`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*domain.Product
	for rows.Next() {
		var p domain.Product
		_ = rows.Scan(&p.ID, &p.Name, &p.Category, &p.Price, &p.Stock)
		out = append(out, &p)
	}
	return out, rows.Err()
}
