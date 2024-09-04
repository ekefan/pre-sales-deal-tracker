-- name: FilterDeals :many
SELECT id, pitch_id, sales_rep_name, customer_name, service_to_render, status, status_tag, current_pitch_request, net_total_cost, profit, created_at, updated_at, closed_at, awarded
FROM deals
WHERE 
    (customer_name = $1 OR $1 IS NULL) AND
    (service_to_render && $2::TEXT[] OR $2 IS NULL) AND
    (status = $3 OR $3 IS NULL) AND
    (profit >= $4 OR $4 IS NULL) AND
    (profit <= $5 OR $5 IS NULL) AND
    (awarded = $6 OR $6 IS NULL) AND
    (sales_rep_name = $7 OR $7 IS NULL)
ORDER BY id
LIMIT $8
OFFSET $9;


-- name: CountFilteredDeals :one
SELECT COUNT(*)
FROM deals
WHERE 
    (customer_name = $1 OR $1 IS NULL) AND
    (service_to_render && $2::TEXT[] OR $2 IS NULL) AND
    (status = $3 OR $3 IS NULL) AND
    (profit >= $4 OR $4 IS NULL) AND
    (profit <= $5 OR $5 IS NULL) AND
    (awarded = $6 OR $6 IS NULL) AND
    (sales_rep_name = $7 OR $7 IS NULL);

-- name: GetDealsByStatus :many
SELECT * FROM deals
WHERE status = $1
ORDER BY id;

-- name: GetDealsById :one
SELECT * FROM deals
WHERE id = $1
LIMIT 1;

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