-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders_reservations (
    order_id     bigint not null,
    warehouse_id bigint not null,
    sku          bigint not null,
    count        int not null,

    PRIMARY KEY (order_id, warehouse_id, sku)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders_reservations;
-- +goose StatementEnd
