-- Create "accounts" table
CREATE TABLE "public"."accounts" (
  "type" bigint NOT NULL,
  "balance" numeric(19,4) NOT NULL,
  "name" character varying(255) NOT NULL,
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "created_by" uuid NULL,
  "updated_at" timestamptz NULL,
  "updated_by" uuid NULL,
  "is_deleted" boolean NULL,
  "deleted_at" timestamptz NULL,
  "user_id" uuid NOT NULL,
  PRIMARY KEY ("id", "created_at")
);
-- Create index "idx_accounts_deleted_at" to table: "accounts"
CREATE INDEX "idx_accounts_deleted_at" ON "public"."accounts" ("deleted_at");
-- Create "attachments" table
CREATE TABLE "public"."attachments" (
  "type" bigint NOT NULL,
  "content_url" character varying(255) NOT NULL,
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "created_by" uuid NULL,
  "updated_at" timestamptz NULL,
  "updated_by" uuid NULL,
  "is_deleted" boolean NULL,
  "deleted_at" timestamptz NULL,
  "user_id" uuid NOT NULL,
  PRIMARY KEY ("id", "created_at")
);
-- Create index "idx_attachments_deleted_at" to table: "attachments"
CREATE INDEX "idx_attachments_deleted_at" ON "public"."attachments" ("deleted_at");
-- Create "balance_totals" table
CREATE TABLE "public"."balance_totals" (
  "income_total" numeric(19,4) NOT NULL,
  "expense_total" numeric(19,4) NOT NULL,
  "month" character varying(2) NOT NULL,
  "year" character varying(4) NOT NULL,
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "created_by" uuid NULL,
  "updated_at" timestamptz NULL,
  "updated_by" uuid NULL,
  "is_deleted" boolean NULL,
  "deleted_at" timestamptz NULL,
  "user_id" uuid NOT NULL,
  PRIMARY KEY ("id", "created_at")
);
-- Create index "idx_balance_totals_deleted_at" to table: "balance_totals"
CREATE INDEX "idx_balance_totals_deleted_at" ON "public"."balance_totals" ("deleted_at");
-- Create "categories" table
CREATE TABLE "public"."categories" (
  "name" character varying(255) NOT NULL,
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "created_by" uuid NULL,
  "updated_at" timestamptz NULL,
  "updated_by" uuid NULL,
  "is_deleted" boolean NULL,
  "deleted_at" timestamptz NULL,
  "user_id" uuid NOT NULL,
  PRIMARY KEY ("id", "created_at")
);
-- Create index "idx_categories_deleted_at" to table: "categories"
CREATE INDEX "idx_categories_deleted_at" ON "public"."categories" ("deleted_at");
-- Create index "idx_category_name" to table: "categories"
CREATE UNIQUE INDEX "idx_category_name" ON "public"."categories" ("name");
