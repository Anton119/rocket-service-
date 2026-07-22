-- +goose Up
UPDATE parts SET stock_quantity = stock_quantity + 1000;

-- +goose Down
UPDATE parts SET stock_quantity = GREATEST(stock_quantity - 1000, 0);
