-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS notifications (
    record_id  bigint not null PRIMARY KEY,  -- int64
    user_id    bigint not null,              -- int64
    message    text not null,                -- string
    created_at timestamp not null            -- time
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS notifications;
-- +goose StatementEnd
