-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS cart_items (
    user_id    bigint not null,     -- int64
    sku        bigint not null,     -- uint32
    count      int not null,        -- uint16

    PRIMARY KEY (user_id, sku)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS cart_items;
-- +goose StatementEnd
