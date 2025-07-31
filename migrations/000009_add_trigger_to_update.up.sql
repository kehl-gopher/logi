CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at := NOW();

    IF TG_TABLE_NAME = 'sender' THEN
        UPDATE auth
        SET updated_at = NOW()
        WHERE id = NEW.id;
    ELSIF TG_TABLE_NAME = 'auth' THEN
        UPDATE sender
        SET updated_at = NOW()
        WHERE id = NEW.id;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER trigger_update_sender
BEFORE UPDATE ON sender
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER trigger_update_auth
BEFORE UPDATE ON auth
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();
