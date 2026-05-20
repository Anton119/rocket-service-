package converter

import (
	"errors"

	"github.com/google/uuid"

	errs "github.com/Anton119/rocket-service-/payment/internal/errors"
	"github.com/Anton119/rocket-service-/payment/internal/model"
	"github.com/Anton119/rocket-service-/payment/internal/service/input"
	paymentv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/payment/v1"
)

// PayOrderInputFromRequest собирает вход use case из rpc PayOrder.
func PayOrderInputFromRequest(req *paymentv1.PayOrderRequest) (input.PayOrderInput, error) {
	if req.GetOrderUuid() == "" {
		return input.PayOrderInput{}, errs.ErrEmptyOrderUUID
	}

	orderUUID, err := uuid.Parse(req.GetOrderUuid())
	if err != nil {
		return input.PayOrderInput{}, errs.ErrInvalidOrderUUID
	}

	if req.GetPaymentMethod() == paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED {
		return input.PayOrderInput{}, errs.ErrPaymentMethodUnspecified
	}

	method, ok := PaymentMethodFromProto(req.GetPaymentMethod())
	if !ok {
		return input.PayOrderInput{}, errs.ErrInvalidPaymentMethod
	}

	return input.PayOrderInput{
		OrderUUID: orderUUID,
		Method:    method,
	}, nil
}

// PayOrderResponseFromResult собирает protobuf-ответ из результата use case.
func PayOrderResponseFromResult(out input.PayOrderResult) *paymentv1.PayOrderResponse {
	return &paymentv1.PayOrderResponse{
		TransactionUuid: out.TransactionUUID.String(),
	}
}

// PaymentMethodFromProto переводит protobuf enum в доменный тип.
func PaymentMethodFromProto(m paymentv1.PaymentMethod) (model.PaymentMethod, bool) {
	switch m {
	case paymentv1.PaymentMethod_PAYMENT_METHOD_CARD:
		return model.PaymentMethodCard, true
	case paymentv1.PaymentMethod_PAYMENT_METHOD_SBP:
		return model.PaymentMethodSBP, true
	case paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:
		return model.PaymentMethodCreditCard, true
	case paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
		return model.PaymentMethodInvestorMoney, true
	default:
		return model.PaymentMethodUnspecified, false
	}
}

// IsValidationError возвращает true для ошибок разбора/валидации запроса на API-слое.
func IsValidationError(err error) bool {
	return errors.Is(err, errs.ErrEmptyOrderUUID) ||
		errors.Is(err, errs.ErrInvalidOrderUUID) ||
		errors.Is(err, errs.ErrPaymentMethodUnspecified) ||
		errors.Is(err, errs.ErrInvalidPaymentMethod)
}

// ValidationErrorMessage возвращает текст для gRPC InvalidArgument.
func ValidationErrorMessage(err error) string {
	if errors.Is(err, errs.ErrEmptyOrderUUID) {
		return errs.ErrEmptyOrderUUID.Error()
	}
	if errors.Is(err, errs.ErrInvalidOrderUUID) {
		return errs.ErrInvalidOrderUUID.Error()
	}
	if errors.Is(err, errs.ErrPaymentMethodUnspecified) {
		return errs.ErrPaymentMethodUnspecified.Error()
	}
	if errors.Is(err, errs.ErrInvalidPaymentMethod) {
		return "неизвестный способ оплаты"
	}

	return err.Error()
}
