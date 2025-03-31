-- +goose Up
CREATE TABLE IF NOT EXISTS endpoints
(
    id        TEXT NOT NULL UNIQUE,
    name  TEXT NOT NULL,
    url     TEXT NOT NULL UNIQUE,
    protocol BLOB NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_endpoint_name ON users (name);

-- +goose Down
DROP TABLE IF EXISTS endpoints;