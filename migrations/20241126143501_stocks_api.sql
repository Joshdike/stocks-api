-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS stocks (
    stockid SERIAL PRIMARY KEY,
    name TEXT,
    price INT,
    company TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS stocks;
-- +goose StatementEnd
