-- +goose Up
-- +goose StatementBegin
CREATE TABLE stocks (
    id SERIAL PRIMARY KEY,
    total_count INTEGER NOT NULL,
    reserved INTEGER NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE stocks;
-- +goose StatementEnd
