-- +goose Up

-- +goose StatementBegin
    CREATE TABLE transactions (
    id BIGSERIAL PRIMARY KEY,
    account_id BIGINT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    operation_type_id BIGINT NOT NULL REFERENCES operation_types(id),
    amount NUMERIC(15,2) NOT NULL,
    event_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER updatedat_timestamp_trigger_transactions
    BEFORE UPDATE ON transactions
    FOR EACH ROW
EXECUTE FUNCTION updatedat_timestamp();
-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin
DROP TRIGGER IF EXISTS updatedat_timestamp_trigger_transactions ON transactions;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS transactions;
-- +goose StatementEnd

