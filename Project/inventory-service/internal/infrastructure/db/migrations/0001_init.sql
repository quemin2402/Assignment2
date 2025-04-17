-- +goose Up
CREATE TABLE IF NOT EXISTS products (
                                        id       TEXT PRIMARY KEY,
                                        name     TEXT NOT NULL,
                                        category TEXT NOT NULL,
                                        price    NUMERIC(12,2) NOT NULL,
    stock    INTEGER NOT NULL
    );
-- +goose Down
DROP TABLE products;
