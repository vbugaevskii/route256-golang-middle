-- +goose Up
-- +goose StatementBegin
CREATE TYPE NotificationState AS ENUM (
    'Waiting',
    'Delivered'
);

CREATE TABLE notifications_outbox (
    record_id  bigserial not null PRIMARY KEY,
    key        text not null,
    value      text not null,
    state      NotificationState not null,
    created_at timestamp not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE notifications_outbox;

DROP TYPE NotificationState;
-- +goose StatementEnd
