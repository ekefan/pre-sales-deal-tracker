-- name: UpdateDealSalesName :exec
UPDATE deals SET sales_rep_name = @new_sales_name
WHERE sales_rep_name = @old_sales_name;


-- name: GetDealToUpdateSalesName :one
SELECT * FROM deals
WHERE sales_rep_name = $1
FOR NO KEY UPDATE;


-- name: CreateDeal :execrows
INSERT INTO deals (pitch_id, sales_rep_name, customer_name, services_to_render, department, net_total_cost, profit)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: GetDealPaginated :one
WITH deal_data AS (
    SELECT
        deals.id as deal_id,
        deals.pitch_id,
        deals.sales_rep_name,
        deals.customer_name,
        deals.services_to_render,
        deals.status,
        deals.department,
        deals.net_total_cost,
        deals.profit,
        deals.awarded,
        to_char(deals.updated_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"') AS updated_at,
        to_char(deals.created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"') AS created_at
    FROM deals
    ORDER BY id
    LIMIT $1 OFFSET $2
)
SELECT
    (SELECT COUNT(*) FROM deals) AS total_deals,
    json_agg(deal_data) AS deals
FROM deal_data;

-- name: UpdateDeals :execrows
UPDATE deals
    set services_to_render = $2, status = $3, department = $4, net_total_cost = $5, profit = $6, awarded = $7
WHERE deals.id = $1;

-- name: DeleteDeals :execrows
DELETE FROM deals WHERE id = $1;