-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS cart_items (
    user_id    bigint not null,     -- int64
    sku        bigint not null,     -- uint32
    count      int not null,        -- uint16
    created_at timestamp not null,
    updated_at timestamp,
    deleted_at timestamp,

    PRIMARY KEY (user_id, sku)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS cart_items;
-- +goose StatementEnd
