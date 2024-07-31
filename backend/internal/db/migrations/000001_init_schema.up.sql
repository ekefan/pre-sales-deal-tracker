CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar UNIQUE NOT NULL,
  "role" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password" varchar NOT NULL,
  "updated_at" timestamp DEFAULT null,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "deals" (
  "id" bigserial PRIMARY KEY,
  "pitch_id" bigint UNIQUE NOT NULL,
  "sales_rep_name" varchar NOT NULL,
  "customer_name" varchar NOT NULL,
  "service_to_render" varchar NOT NULL DEFAULT 'none',
  "status" varchar NOT NULL DEFAULT 'ongoing',
  "status_tag" varchar NOT NULL,
  "current_pitch_request" varchar NOT NULL,
  "net_total_cost" numeric(11,2),
  "profit" numeric(11,2),
  "created_at" timestamp DEFAULT (now()) NOT NULL,
  "updated_at" timestamp DEFAULT NULL,
  "closed_at" timestamp DEFAULT NULL,
  "awarded" bool DEFAULT false NOT NULL
);

CREATE TABLE "pitch_requests" (
  "id" bigserial PRIMARY KEY,
  "sales_rep_id" bigint NOT NULL,
  "sales_rep_name" varchar NOT NULL,
  "status" varchar NOT NULL DEFAULT 'ongoing',
  "customer_name" varchar NOT NULL,
  "pitch_tag" varchar NOT NULL,
  "customer_request" varchar NOT NULL DEFAULT 'proposal',
  "request_deadline" timestamp NOT NULL DEFAULT(CURRENT_TIMESTAMP + INTERVAL '3 days'),
  "admin_viewed" bool DEFAULT false NOT NULL,
  "created_at" timestamp DEFAULT (now()) NOT NULL,
  "updated_at" timestamp DEFAULT NULL
);

CREATE INDEX ON "users" ("role");

CREATE INDEX ON "deals" ("service_to_render");

CREATE INDEX ON "deals" ("customer_name");

CREATE INDEX ON "deals" ("customer_name", "service_to_render");

CREATE INDEX ON "deals" ("sales_rep_name");

CREATE INDEX ON "deals" ("pitch_id");

CREATE INDEX ON "deals" ("status");

CREATE INDEX ON "deals" ("awarded");

CREATE INDEX ON "pitch_requests" ("admin_viewed");

CREATE INDEX ON "pitch_requests" ("status");

CREATE INDEX ON "pitch_requests" ("customer_name");

COMMENT ON COLUMN "users"."role" IS 'Role of the user, e.g., sales rep, admin, manager';

COMMENT ON COLUMN "deals"."status" IS 'it is either ongoing or closed';

COMMENT ON COLUMN "deals"."status_tag" IS 'tag for the status can be, working on: survey, proposal, costing depending';

ALTER TABLE "deals" ADD FOREIGN KEY ("pitch_id") REFERENCES "pitch_requests" ("id") ON DELETE SET NULL;

ALTER TABLE "pitch_requests" ADD FOREIGN KEY ("sales_rep_id") REFERENCES "users" ("id") ON DELETE CASCADE;
