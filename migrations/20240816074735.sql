CREATE OR REPLACE FUNCTION update_account_balance()
RETURNS TRIGGER AS $$
BEGIN
    BEGIN
        -- Update or insert the balance in the accounts table
        IF NEW.type = 0 THEN
            UPDATE accounts SET balance = accounts.balance + NEW.amount WHERE user_id = NEW.user_id;
        ELSIF NEW.type = 1 THEN
            UPDATE accounts SET balance = accounts.balance - NEW.amount WHERE user_id = NEW.user_id;
        END IF;
        
        RETURN NEW;
    END;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER trigger_update_account_balance
AFTER INSERT ON transactions
FOR EACH ROW
EXECUTE FUNCTION update_account_balance();