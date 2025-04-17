package repository

import (
	"Assignment2/Project/order-service/internal/domain"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repo struct{ db *pgxpool.Pool }

func New(db *pgxpool.Pool) domain.OrderRepository { return &repo{db} }

func (r *repo) Create(ctx context.Context, o *domain.Order) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, "insert into orders(id,status) values($1,$2)", o.ID, o.Status)
	if err != nil {
		return err
	}
	for _, it := range o.Items {
		_, err = tx.Exec(ctx, "insert into order_items(order_id,product_id,quantity) values($1,$2,$3)",
			o.ID, it.ProductID, it.Quantity)
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (r *repo) Get(ctx context.Context, id string) (*domain.Order, error) {
	var o domain.Order
	o.ID = id
	row := r.db.QueryRow(ctx, "select status from orders where id=$1", id)
	if err := row.Scan(&o.Status); err != nil {
		return nil, err
	}
	rows, err := r.db.Query(ctx, "select product_id,quantity from order_items where order_id=$1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var it domain.OrderItem
		_ = rows.Scan(&it.ProductID, &it.Quantity)
		o.Items = append(o.Items, it)
	}
	return &o, rows.Err()
}

func (r *repo) Update(ctx context.Context, o *domain.Order) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, "update orders set status=$2 where id=$1", o.ID, o.Status)
	if err != nil {
		return err
	}
	_, _ = tx.Exec(ctx, "delete from order_items where order_id=$1", o.ID)
	for _, it := range o.Items {
		_, err = tx.Exec(ctx, "insert into order_items(order_id,product_id,quantity) values($1,$2,$3)",
			o.ID, it.ProductID, it.Quantity)
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (r *repo) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, "delete from orders where id=$1", id)
	return err
}

func (r *repo) List(ctx context.Context) ([]*domain.Order, error) {
	rows, err := r.db.Query(ctx, "select id,status from orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*domain.Order
	for rows.Next() {
		var o domain.Order
		if err := rows.Scan(&o.ID, &o.Status); err != nil {
			return nil, err
		}
		itRows, err := r.db.Query(ctx, "select product_id,quantity from order_items where order_id=$1", o.ID)
		if err != nil {
			return nil, err
		}
		for itRows.Next() {
			var it domain.OrderItem
			_ = itRows.Scan(&it.ProductID, &it.Quantity)
			o.Items = append(o.Items, it)
		}
		itRows.Close()
		out = append(out, &o)
	}
	return out, rows.Err()
}
