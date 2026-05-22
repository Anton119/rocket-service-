package errs

import "errors"

var (
	ErrOrderNotFound = errors.New("order not found")

	ErrPartNotFound       = errors.New("part not found")
	ErrPartOutOfStock     = errors.New("part out of stock")
	ErrOrderPayNotAllowed = errors.New("order pay not allowed")

	ErrOrderAlreadyPaid      = errors.New("заказ уже оплачен, отмена невозможна")
	ErrOrderAlreadyCancelled = errors.New("заказ уже отменён")
	ErrOrderCancelNotAllowed = errors.New("отмена заказа невозможна в текущем статусе")

	ErrInvalidPaymentMethod = errors.New("invalid payment method")
)

// InvalidArgumentError переносит сообщение внешнего сервиса (gRPC InvalidArgument).
type InvalidArgumentError struct {
	Message string
}

func (e *InvalidArgumentError) Error() string {
	if e == nil {
		return ""
	}
	return e.Message
}
