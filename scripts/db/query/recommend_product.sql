-- name: CreateRecommendProduct :execresult
INSERT INTO recommend_products (product_id)
VALUES (?);

-- name: ListRecommendProducts :many
SELECT products.id, products.name, products.price
FROM recommend_products
LEFT JOIN products
ON recommend_products.product_id = products.id ORDER BY products.id ASC;

-- name: DeleteRecommendProduct :exec
DELETE
FROM recommend_products
WHERE id = ?;
