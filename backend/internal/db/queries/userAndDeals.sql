-- name: GetDealsByCustomerName :many
SELECT * FROM deals
WHERE customer_name = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: GetDealsByServicesRendered :many
SELECT * FROM deals 
WHERE service_to_render = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: GetDealsByCustomerAndService :many
SELECT * FROM deals
WHERE service_to_render = $1 
AND customer_name = $2
ORDER BY id
LIMIT $3
OFFSET $4;

-- name: GetDealsByProfit :many
SELECT * FROM deals
WHERE profit >= $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: GetDealsByStatus :many
SELECT * FROM deals
WHERE status = $1
ORDER BY id;


-- name: GetDealsByAward :many
SELECT * FROM deals
WHERE awarded = $1
ORDER BY closed_at
LIMIT $2
OFFSET $3;

-- name: GetDealsBySalesRep :many
SELECT * FROM deals
WHERE sales_rep_name = $1
ORDER BY id 
LIMIT $2
OFFSET $3;

-- name: GetDealsBySalesRepAndAwarded :many
SELECT * FROM deals
WHERE sales_rep_name = $1
AND awarded = $2
ORDER BY id 
LIMIT $3
OFFSET $4;


-- name: GetUser :one
SELECT * FROM users
WHERE username = $1
LIMIT 1;