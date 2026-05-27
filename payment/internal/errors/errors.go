package errs

import "errors"

var (
	ErrEmptyOrderUUID           = errors.New("uuid не может быть пустым")
	ErrInvalidOrderUUID         = errors.New("неверный формат uuid")
	ErrPaymentMethodUnspecified = errors.New("payment_method не указан")
	ErrInvalidPaymentMethod     = errors.New("неизвестный способ оплаты")
)
