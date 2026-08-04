package v1

import (
	"errors"
	"net/http"

	errs "github.com/Anton119/rocket-service-/order/internal/errors"
	orderv1 "github.com/Anton119/rocket-service-/shared/pkg/openapi/order/v1"
)

func mapCreateOrderError(err error) (orderv1.CreateOrderRes, error) {
	switch {
	case errors.Is(err, errs.ErrPartNotFound):
		return &orderv1.CreateOrderNotFound{
			Code:    http.StatusNotFound,
			Message: "деталь не найдена",
		}, nil
	case errors.Is(err, errs.ErrPartOutOfStock):
		return &orderv1.CreateOrderConflict{
			Code:    http.StatusConflict,
			Message: "деталь отсутствует на складе",
		}, nil
	case errors.Is(err, errs.ErrUnauthorized):
		return &orderv1.CreateOrderBadRequest{
			Code:    http.StatusUnauthorized,
			Message: errs.ErrUnauthorized.Error(),
		}, nil
	default:
		var ia *errs.InvalidArgumentError
		if errors.As(err, &ia) {
			return &orderv1.CreateOrderBadRequest{
				Code:    http.StatusBadRequest,
				Message: ia.Message,
			}, nil
		}

		return nil, err
	}
}

func mapGetOrderError(err error) (orderv1.GetOrderRes, error) {
	switch {
	case errors.Is(err, errs.ErrOrderNotFound):
		return &orderv1.GetOrderNotFound{
			Code:    http.StatusNotFound,
			Message: "заказ не найден",
		}, nil
	default:
		return nil, err
	}
}

func mapPayOrderError(err error) (orderv1.PayOrderRes, error) {
	switch {
	case errors.Is(err, errs.ErrOrderNotFound):
		return &orderv1.PayOrderNotFound{
			Code:    http.StatusNotFound,
			Message: "заказ не найден",
		}, nil
	case errors.Is(err, errs.ErrOrderPayNotAllowed):
		return &orderv1.PayOrderConflict{
			Code:    http.StatusConflict,
			Message: "оплата невозможна в текущем статусе",
		}, nil
	case errors.Is(err, errs.ErrInvalidPaymentMethod):
		return &orderv1.PayOrderBadRequest{
			Code:    http.StatusBadRequest,
			Message: "неизвестный способ оплаты",
		}, nil
	default:
		var ia *errs.InvalidArgumentError
		if errors.As(err, &ia) {
			return &orderv1.PayOrderBadRequest{
				Code:    http.StatusBadRequest,
				Message: ia.Message,
			}, nil
		}

		return nil, err
	}
}

func mapCancelOrderError(err error) (orderv1.CancelOrderRes, error) {
	switch {
	case errors.Is(err, errs.ErrOrderNotFound):
		return &orderv1.CancelOrderNotFound{
			Code:    http.StatusNotFound,
			Message: "заказ не найден",
		}, nil
	case errors.Is(err, errs.ErrOrderAlreadyPaid):
		return &orderv1.CancelOrderConflict{
			Code:    http.StatusConflict,
			Message: errs.ErrOrderAlreadyPaid.Error(),
		}, nil
	case errors.Is(err, errs.ErrOrderAlreadyCancelled):
		return &orderv1.CancelOrderConflict{
			Code:    http.StatusConflict,
			Message: errs.ErrOrderAlreadyCancelled.Error(),
		}, nil
	case errors.Is(err, errs.ErrOrderCancelNotAllowed):
		return &orderv1.CancelOrderConflict{
			Code:    http.StatusConflict,
			Message: errs.ErrOrderCancelNotAllowed.Error(),
		}, nil
	default:
		return nil, err
	}
}
