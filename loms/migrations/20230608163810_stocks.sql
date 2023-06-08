-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS stocks (
    warehouse_id bigint not null,
    sku          bigint not null,
    count        int not null,
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
