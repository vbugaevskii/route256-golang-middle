-- +goose Up
-- +goose StatementBegin
INSERT INTO stocks (warehouse_id, sku, count) VALUES
	(1, 773587830, 3),
	(2, 773587830, 2),
	(3, 773587830, 1),
    (1, 773596051, 5);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM stocks;
-- +goose StatementEnd
