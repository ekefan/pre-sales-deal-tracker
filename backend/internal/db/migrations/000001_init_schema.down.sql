-- Drop foreign key constraints 
ALTER TABLE "deals" DROP CONSTRAINT IF EXISTS "deals_pitch_id_fkey";
ALTER TABLE "pitch_requests" DROP CONSTRAINT IF EXISTS "pitch_requests_sales_rep_id_fkey";

-- Drop indexes
DROP INDEX IF EXISTS "users_role_idx";
DROP INDEX IF EXISTS "deals_service_to_render_idx";
DROP INDEX IF EXISTS "deals_customer_name_idx";
DROP INDEX IF EXISTS "deals_customer_name_service_to_render_idx";
DROP INDEX IF EXISTS "deals_sales_rep_name_idx";
DROP INDEX IF EXISTS "deals_pitch_id_idx";
DROP INDEX IF EXISTS "deals_status_idx";
DROP INDEX IF EXISTS "deals_awarded_idx";
DROP INDEX IF EXISTS "pitch_requests_admin_viewed_idx";
DROP INDEX IF EXISTS "pitch_requests_status_idx";
DROP INDEX IF EXISTS "pitch_requests_customer_name_idx";

-- Drop tables
DROP TABLE IF EXISTS "pitch_requests";
DROP TABLE IF EXISTS "deals";
DROP TABLE IF EXISTS "users";
