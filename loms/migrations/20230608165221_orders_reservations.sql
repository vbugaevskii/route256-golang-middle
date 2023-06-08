-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders_reservations (
    order_id     bigint not null,  -- int64
    warehouse_id bigint not null,  -- int64
    sku          bigint not null,  -- uint32
    count        int not null,

    PRIMARY KEY (order_id, warehouse_id, sku)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders_reservations;
-- +goose StatementEnd
