package errs

import "errors"

var (
	ErrOrderNotFound = errors.New("заказ не найден")

	ErrPartNotFound       = errors.New("деталь не найдена")
	ErrPartOutOfStock     = errors.New("деталь отсутствует на складе")
	ErrOrderPayNotAllowed = errors.New("оплата невозможна в текущем статусе")

	ErrOrderAlreadyPaid      = errors.New("заказ уже оплачен, отмена невозможна")
	ErrOrderAlreadyCancelled = errors.New("заказ уже отменён")
	ErrOrderCancelNotAllowed = errors.New("отмена заказа невозможна в текущем статусе")

	ErrInvalidPaymentMethod = errors.New("неизвестный способ оплаты")
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
