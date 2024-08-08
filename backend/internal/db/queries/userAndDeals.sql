-- name: FilterDeals :many
WITH params AS (
    SELECT
        $1::TEXT AS customer_name,
        $2::TEXT AS service_to_render,
        $3::TEXT AS status,
        $4::NUMERIC AS min_profit,
        $5::NUMERIC AS max_profit,
        $6::BOOLEAN AS awarded,
        $7::TEXT AS sales_rep_name
)
SELECT * FROM deals
WHERE 
    (customer_name = (SELECT customer_name FROM params) OR (SELECT customer_name FROM params) IS NULL) AND
    (service_to_render = (SELECT service_to_render FROM params) OR (SELECT service_to_render FROM params) IS NULL) AND
    (status = (SELECT status FROM params) OR (SELECT status FROM params) IS NULL) AND
    (profit >= (SELECT min_profit FROM params) OR (SELECT min_profit FROM params) IS NULL) AND
    (profit <= (SELECT max_profit FROM params) OR (SELECT max_profit FROM params) IS NULL) AND
    (awarded = (SELECT awarded FROM params) OR (SELECT awarded FROM params) IS NULL) AND
    (sales_rep_name = (SELECT sales_rep_name FROM params) OR (SELECT sales_rep_name FROM params) IS NULL)
ORDER BY id
LIMIT $8
OFFSET $9;


-- name: CountFilteredDeals :one
WITH params AS (
    SELECT
        $1::TEXT AS customer_name,
        $2::TEXT AS service_to_render,
        $3::TEXT AS status,
        $4::NUMERIC AS min_profit,
        $5::NUMERIC AS max_profit,
        $6::BOOLEAN AS awarded,
        $7::TEXT AS sales_rep_name
)
SELECT COUNT(*) FROM deals
WHERE 
    (customer_name = (SELECT customer_name FROM params) OR (SELECT customer_name FROM params) IS NULL) AND
    (service_to_render = (SELECT service_to_render FROM params) OR (SELECT service_to_render FROM params) IS NULL) AND
    (status = (SELECT status FROM params) OR (SELECT status FROM params) IS NULL) AND
    (profit >= (SELECT min_profit FROM params) OR (SELECT min_profit FROM params) IS NULL) AND
    (profit <= (SELECT max_profit FROM params) OR (SELECT max_profit FROM params) IS NULL) AND
    (awarded = (SELECT awarded FROM params) OR (SELECT awarded FROM params) IS NULL) AND
    (sales_rep_name = (SELECT sales_rep_name FROM params) OR (SELECT sales_rep_name FROM params) IS NULL);


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