package errs

import "errors"

var (
	// ErrEmptyOrderUUID — пустой order_uuid в запросе.
	ErrEmptyOrderUUID = errors.New("uuid не может быть пустым")
	// ErrInvalidOrderUUID — order_uuid не является UUID.
	ErrInvalidOrderUUID = errors.New("неверный формат uuid")
	// ErrPaymentMethodUnspecified — payment_method не указан.
	ErrPaymentMethodUnspecified = errors.New("payment_method не указан")
	// ErrInvalidPaymentMethod — неизвестный способ оплаты.
	ErrInvalidPaymentMethod = errors.New("invalid payment method")
)
