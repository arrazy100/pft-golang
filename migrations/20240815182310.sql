-- Modify "accounts" table
ALTER TABLE "public"."accounts" DROP CONSTRAINT "accounts_pkey", ALTER COLUMN "created_at" DROP NOT NULL, ADD PRIMARY KEY ("id");
-- Modify "attachments" table
ALTER TABLE "public"."attachments" DROP CONSTRAINT "attachments_pkey", ALTER COLUMN "created_at" DROP NOT NULL, ADD PRIMARY KEY ("id");
-- Modify "balance_totals" table
ALTER TABLE "public"."balance_totals" DROP CONSTRAINT "balance_totals_pkey", ALTER COLUMN "created_at" DROP NOT NULL, ADD PRIMARY KEY ("id");
-- Modify "categories" table
ALTER TABLE "public"."categories" DROP CONSTRAINT "categories_pkey", ALTER COLUMN "created_at" DROP NOT NULL, ADD PRIMARY KEY ("id");
