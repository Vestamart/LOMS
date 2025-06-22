-- +goose Up
-- +goose StatementBegin
INSERT INTO stocks (id, total_count, reserved)
VALUES  (773297411, 150, 10),
        (1002, 200, 20),
        (1003, 250, 30),
        (1004, 300, 40),
        (1005, 350, 50);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM stocks
WHERE id IN (773297411, 1002, 1003, 1004, 1005);
-- +goose StatementEnd
