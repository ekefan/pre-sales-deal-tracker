-- name: UpdateDealSalesName :exec
UPDATE deals SET sales_rep_name = @new_sales_name
WHERE sales_rep_name = @old_sales_name;


-- name: GetDealToUpdateSalesName :one
SELECT * FROM deals
WHERE sales_rep_name = $1
FOR NO KEY UPDATE;