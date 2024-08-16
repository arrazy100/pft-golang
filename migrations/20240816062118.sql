CREATE OR REPLACE FUNCTION update_balance_total()
RETURNS TRIGGER AS $$
BEGIN
    -- Extract month and year from the transaction date
    DECLARE
        trans_month INT := EXTRACT(MONTH FROM NEW.transaction_date);
        trans_year INT := EXTRACT(YEAR FROM NEW.transaction_date);
    BEGIN
        -- Update or insert the balance in the balance_total table
        IF NEW.type = 0 THEN
            INSERT INTO balance_totals (id, user_id, month, year, income_total, expense_total, created_at, created_by)
            VALUES (uuid_generate_v4(), NEW.user_id, trans_month, trans_year, NEW.amount, 0, CURRENT_TIMESTAMP, NEW.user_id)
            ON CONFLICT (user_id, month, year)
            DO UPDATE SET income_total = balance_totals.income_total + NEW.amount;
        ELSIF NEW.type = 1 THEN
            INSERT INTO balance_totals (id, user_id, month, year, income_total, expense_total, created_at, created_by)
            VALUES (uuid_generate_v4(), NEW.user_id, trans_month, trans_year, 0, NEW.amount, CURRENT_TIMESTAMP, NEW.user_id)
            ON CONFLICT (user_id, month, year)
            DO UPDATE SET expense_total = balance_totals.expense_total + NEW.amount;
        END IF;
        
        RETURN NEW;
    END;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER trigger_update_balance_total
AFTER INSERT ON transactions
FOR EACH ROW
EXECUTE FUNCTION update_balance_total();