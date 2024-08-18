-- Modify "transactions" table
ALTER TABLE "public"."transactions" ALTER COLUMN "attachment_id" DROP NOT NULL;
-- Modify "attachments" table
ALTER TABLE "public"."attachments" ADD COLUMN "transaction_id" uuid NULL, ADD
 CONSTRAINT "fk_attachments_transaction" FOREIGN KEY ("transaction_id") REFERENCES "public"."transactions" ("id") ON UPDATE CASCADE ON DELETE CASCADE;
