-- +goose Up
-- +goose StatementBegin
CREATE TYPE StatusType AS ENUM (
    'New',
	'AwaitingPayment',
	'Failed',
	'Payed',
	'Cancelled'
);

CREATE TABLE IF NOT EXISTS orders_status (
    order_id    bigserial not null PRIMARY KEY,
    user_id     bigint not null,
    status      StatusType not null,
    created_at  timestamp not null,
    updated_at  timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders_status;

DROP TYPE StatusType;
-- +goose StatementEnd
