-- +goose Up
-- +goose StatementBegin
CREATE TYPE StatusType AS ENUM (
    'New',
	'AwaitingPayment',
	'Failed',
	'Payed',
	'Cancelled'
);

CREATE TABLE IF NOT EXISTS orders (
    order_id    bigserial not null PRIMARY KEY,  -- int64
    user_id     bigint not null,                 -- int64
    status      StatusType not null,
    created_at  timestamp not null,
    updated_at  timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;

DROP TYPE StatusType;
-- +goose StatementEnd
