-- +goose Up

-- +goose StatementBegin
CREATE TABLE accounts (
    id BIGSERIAL PRIMARY KEY,
    document_number TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER updatedat_timestamp_trigger_accounts
    BEFORE UPDATE ON accounts
    FOR EACH ROW EXECUTE FUNCTION updatedat_timestamp();
-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin
DROP TRIGGER IF EXISTS updatedat_timestamp_trigger_accounts ON accounts;
-- +goose StatementEnd

-- +goose StatementBegin
DROP FUNCTION IF EXISTS updatedat_timestamp;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS accounts;
-- +goose StatementEnd
