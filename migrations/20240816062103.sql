-- Create index "idx_balance_total_month_year" to table: "balance_totals"
CREATE UNIQUE INDEX "idx_balance_total_month_year" ON "public"."balance_totals" ("month", "year", "user_id");
