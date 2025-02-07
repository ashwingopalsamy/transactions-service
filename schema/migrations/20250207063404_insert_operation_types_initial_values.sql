-- +goose Up

-- +goose StatementBegin
INSERT INTO
    operation_types (description)
VALUES
    ('Normal Purchase'),
    ('Purchase with Installments'),
    ('Withdrawal'),
    ('Credit Voucher');
-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin
DELETE FROM operation_types
WHERE description IN ('Normal Purchase', 'Purchase with Installments', 'Withdrawal', 'Credit Voucher');
-- +goose StatementEnd

