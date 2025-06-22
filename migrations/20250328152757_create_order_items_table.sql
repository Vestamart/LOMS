-- +goose Up
-- +goose StatementBegin
CREATE TABLE order_items (
    order_id BIGINT NOT NULL,
    item_id BIGINT NOT NULL,
    PRIMARY KEY (order_id, item_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE order_items;
-- +goose StatementEnd
