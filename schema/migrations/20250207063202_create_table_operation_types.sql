-- +goose Up

-- +goose StatementBegin
CREATE TABLE operation_types (
     id BIGSERIAL PRIMARY KEY,
     description TEXT NOT NULL,
     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER updatedat_timestamp_trigger_operation_types
    BEFORE UPDATE ON operation_types
    FOR EACH ROW
EXECUTE FUNCTION updatedat_timestamp();
-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin
DROP TRIGGER IF EXISTS updatedat_timestamp_trigger_operation_types ON operation_types;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS operation_types;
-- +goose StatementEnd
