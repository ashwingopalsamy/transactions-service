-- +goose Up

-- +goose StatementBegin
ALTER TABLE transactions
    ADD COLUMN balance NUMERIC(15,2) NOT NULL DEFAULT 0.00;
-- +goose StatementEnd

-- +goose StatementBegin
UPDATE transactions SET balance = amount;
-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin
ALTER TABLE transactions
    DROP COLUMN balance;
-- +goose StatementEnd
