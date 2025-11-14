-- name: ListProducts :many
SELECT * FROM products; 

-- name: FindProductByID :one
SELECT * FROM products WHERE id = ? LIMIT 1;

-- name: CreateProduct :one
INSERT INTO products (id, name, price_in_cents, quantity, created_at) VALUES (?, ?, ?, ?, ?) RETURNING *;
