-- Create "user_data" table
CREATE TABLE "user_data" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "name" text NULL,
  "currency" text NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_user_data_deleted_at" to table: "user_data"
CREATE INDEX "idx_user_data_deleted_at" ON "user_data" ("deleted_at");
