package order

import "context"

// TxManager определяет контракт для управления транзакциями.
// Реализация — *manager.Manager из go-transaction-manager.
type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}
