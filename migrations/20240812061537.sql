-- Create "transactions" table
CREATE TABLE "public"."transactions" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "is_deleted" boolean NULL,
  "user_id" uuid NOT NULL,
  "category_id" uuid NOT NULL,
  "description" character varying(255) NOT NULL,
  "account_id" uuid NOT NULL,
  "attachment_id" uuid NOT NULL,
  "amount" numeric(10,4) NOT NULL,
  "type" bigint NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_transactions_deleted_at" to table: "transactions"
CREATE INDEX "idx_transactions_deleted_at" ON "public"."transactions" ("deleted_at", "deleted_at");
-- Drop "user_data" table
DROP TABLE "public"."user_data";
