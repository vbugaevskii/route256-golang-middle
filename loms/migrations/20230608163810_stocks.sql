-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS stocks (
    warehouse_id bigint not null,     -- int64
    sku          bigint not null,     -- uint32
    count        int not null,        -- uint16
    created_at   timestamp not null,
    updated_at   timestamp,
    deleted_at   timestamp,

    PRIMARY KEY (warehouse_id, sku)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS stocks;
-- +goose StatementEnd
