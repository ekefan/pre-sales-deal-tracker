-- name: FilterDeals :many
SELECT * FROM deals
WHERE 
    (customer_name ILIKE $1 OR $1 IS NULL) AND
    (service_to_render ILIKE $2 OR $2 IS NULL) AND
    (status = $3 OR $3 IS NULL) AND
    (profit >= $4 OR $4 IS NULL) AND
    (profit <= $5 OR $5 IS NULL) AND
    (awarded = $6 OR $6 IS NULL) AND
    (sales_rep_name ILIKE $7 OR $7 IS NULL)
ORDER BY id 
LIMIT $8
OFFSET $9;


-- name: CountFilteredDeals :one
SELECT COUNT(*)
FROM deals
WHERE 
    (customer_name ILIKE $1 OR $1 IS NULL) AND
    (service_to_render ILIKE $2 OR $2 IS NULL) AND
    (status = $3 OR $3 IS NULL) AND
    (profit >= $4 OR $4 IS NULL) AND
    (profit <= $5 OR $5 IS NULL) AND
    (awarded = $6 OR $6 IS NULL) AND
    (sales_rep_name ILIKE $7 OR $7 IS NULL);


-- name: GetDealsByStatus :many
SELECT * FROM deals
WHERE status = $1
ORDER BY id;


-- name: GetDealsBySalesRep :many
SELECT * FROM deals
WHERE sales_rep_name = $1
ORDER BY id 
LIMIT $2
OFFSET $3;


-- name: GetUser :one
SELECT * FROM users
WHERE username = $1
LIMIT 1;