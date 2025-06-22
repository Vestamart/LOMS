-- +goose Up
-- +goose StatementBegin
CREATE TABLE items (
    id BIGSERIAL PRIMARY KEY,
    sku INTEGER NOT NULL,
    count INTEGER NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE items;
-- +goose StatementEnd
