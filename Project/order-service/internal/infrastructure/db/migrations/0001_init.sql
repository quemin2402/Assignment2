CREATE TABLE orders (
                        id TEXT PRIMARY KEY,
                        status TEXT NOT NULL
);
CREATE TABLE order_items (
                             order_id TEXT REFERENCES orders(id) ON DELETE CASCADE,
                             product_id TEXT NOT NULL,
                             quantity INT NOT NULL,
                             PRIMARY KEY (order_id, product_id)
);
-- +goose Down
DROP TABLE order_items;
DROP TABLE orders;