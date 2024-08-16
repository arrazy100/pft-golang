-- Create "transactions" table
CREATE TABLE "public"."transactions" (
  "description" character varying(255) NOT NULL,
  "amount" numeric(19,4) NOT NULL,
  "type" bigint NOT NULL,
  "transaction_date" timestamptz NULL,
  "category_id" uuid NOT NULL,
  "account_id" uuid NOT NULL,
  "attachment_id" uuid NOT NULL,
  "created_at" timestamptz NULL,
  "created_by" uuid NULL,
  "updated_at" timestamptz NULL,
  "updated_by" uuid NULL,
  "is_deleted" boolean NULL,
  "deleted_at" timestamptz NULL,
  "id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_transactions_account" FOREIGN KEY ("account_id") REFERENCES "public"."accounts" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_transactions_attachment" FOREIGN KEY ("attachment_id") REFERENCES "public"."attachments" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_transactions_category" FOREIGN KEY ("category_id") REFERENCES "public"."categories" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_transactions_deleted_at" to table: "transactions"
CREATE INDEX "idx_transactions_deleted_at" ON "public"."transactions" ("deleted_at");
