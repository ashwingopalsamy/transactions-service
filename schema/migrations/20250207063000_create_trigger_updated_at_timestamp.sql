-- +goose Up

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION updatedat_timestamp()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin
DROP FUNCTION IF EXISTS updatedat_timestamp;
-- +goose StatementEnd
