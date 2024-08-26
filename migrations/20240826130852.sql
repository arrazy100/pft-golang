CREATE OR REPLACE FUNCTION update_account_balance()
RETURNS TRIGGER AS $$
BEGIN
    -- Handle INSERT case
    IF TG_OP = 'INSERT' THEN
        IF NEW.type = 0 THEN
            UPDATE accounts SET balance = accounts.balance + NEW.amount WHERE user_id = NEW.user_id;
        ELSIF NEW.type = 1 THEN
            UPDATE accounts SET balance = accounts.balance - NEW.amount WHERE user_id = NEW.user_id;
        END IF;
    -- Handle UPDATE case
    ELSIF TG_OP = 'UPDATE' THEN
        -- Revert the old transaction amount
        IF OLD.type = 0 THEN
            UPDATE accounts SET balance = accounts.balance - OLD.amount WHERE user_id = OLD.user_id;
        ELSIF OLD.type = 1 THEN
            UPDATE accounts SET balance = accounts.balance + OLD.amount WHERE user_id = OLD.user_id;
        END IF;
        
        -- Apply the new transaction amount
        IF NEW.type = 0 THEN
            UPDATE accounts SET balance = accounts.balance + NEW.amount WHERE user_id = NEW.user_id;
        ELSIF NEW.type = 1 THEN
            UPDATE accounts SET balance = accounts.balance - NEW.amount WHERE user_id = NEW.user_id;
        END IF;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER trigger_update_account_balance
AFTER INSERT OR UPDATE ON transactions
FOR EACH ROW
EXECUTE FUNCTION update_account_balance();