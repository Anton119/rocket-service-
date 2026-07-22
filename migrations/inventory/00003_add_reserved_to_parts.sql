-- +goose Up
ALTER TABLE parts ADD COLUMN reserved INT NOT NULL DEFAULT 0;

-- +goose Down
ALTER TABLE parts DROP COLUMN IF EXISTS reserved;
