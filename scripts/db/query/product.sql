-- name: CreateProduct :execresult
INSERT INTO products (name, price)
VALUES (?, ?);

-- name: GetProduct :one
SELECT *
FROM products
WHERE id = ? LIMIT 1;

-- name: ListProducts :many
SELECT *
FROM products
ORDER BY id LIMIT ?
OFFSET ?;

-- name: UpdateProduct :execresult
UPDATE products
SET name     = ?,
    price    = ?,
    updated_at = ?
WHERE id = ?;

-- name: DeleteProduct :exec
DELETE
FROM products
WHERE id = ?;