package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"user-service/internal/domain"
)

const createTable = `
CREATE TABLE IF NOT EXISTS users(
  id text primary key,
  username text unique not null,
  email text not null,
  password text not null
);`

type repo struct{ db *pgxpool.Pool }

func New(db *pgxpool.Pool) (domain.UserRepository, error) {
	if _, err := db.Exec(context.Background(), createTable); err != nil {
		return nil, err
	}
	return &repo{db}, nil
}

func (r *repo) Create(ctx context.Context, u *domain.User) error {
	_, err := r.db.Exec(ctx,
		`insert into users(id,username,email,password) values($1,$2,$3,$4)`,
		u.ID, u.Username, u.Email, u.Password)
	return err
}

func (r *repo) GetByUsername(ctx context.Context, un string) (*domain.User, error) {
	row := r.db.QueryRow(ctx, `select id,username,email,password from users where username=$1`, un)
	var u domain.User
	if err := row.Scan(&u.ID, &u.Username, &u.Email, &u.Password); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	row := r.db.QueryRow(ctx, `select id,username,email,password from users where id=$1`, id)
	var u domain.User
	if err := row.Scan(&u.ID, &u.Username, &u.Email, &u.Password); err != nil {
		return nil, err
	}
	return &u, nil
}
