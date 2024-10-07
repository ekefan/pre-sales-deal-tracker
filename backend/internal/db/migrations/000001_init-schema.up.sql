CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar UNIQUE NOT NULL,
  "role" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password" varchar NOT NULL,
  "is_master" boolean NOT NULL DEFAULT false,
  "password_changed" bool NOT NULL DEFAULT false,
  "updated_at" timestamp NOT NULL DEFAULT (now()),
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "deals" (
  "id" bigserial PRIMARY KEY,
  "pitch_id" bigint UNIQUE,
  "sales_rep_name" varchar NOT NULL,
  "customer_name" varchar NOT NULL,
  "services_to_render" text[] NOT NULL DEFAULT '{ }',
  "status" varchar NOT NULL DEFAULT 'ongoing',
  "department" varchar NOT NULL,
  "net_total_cost" numeric(11,2) NOT NULL,
  "profit" numeric(11,2) NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now()),
  "closed_at" timestamp NOT NULL DEFAULT ('0001-01-01 00:00:00'),
  "awarded" bool NOT NULL DEFAULT false
);

CREATE TABLE "pitch_requests" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "customer_name" varchar NOT NULL,
  "customer_request" varchar NOT NULL DEFAULT 'proposal',
  "admin_task" varchar NOT NULL,
  "admin_deadline" timestamp NOT NULL,
  "admin_viewed" bool NOT NULL DEFAULT false,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);

ALTER TABLE "pitch_requests" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "deals" ADD FOREIGN KEY ("pitch_id") REFERENCES "pitch_requests" ("id") ON DELETE SET NULL;

CREATE UNIQUE INDEX unique_master_user ON "users" (is_master) WHERE is_master = true;