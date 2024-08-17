-- Modify "accounts" table
ALTER TABLE "public"."accounts" DROP COLUMN "is_deleted", DROP COLUMN "deleted_at";
-- Modify "attachments" table
ALTER TABLE "public"."attachments" DROP COLUMN "is_deleted", DROP COLUMN "deleted_at";
-- Modify "balance_totals" table
ALTER TABLE "public"."balance_totals" DROP COLUMN "is_deleted", DROP COLUMN "deleted_at";
-- Modify "categories" table
ALTER TABLE "public"."categories" DROP COLUMN "is_deleted", DROP COLUMN "deleted_at";
-- Modify "transactions" table
ALTER TABLE "public"."transactions" DROP COLUMN "is_deleted", DROP COLUMN "deleted_at";
