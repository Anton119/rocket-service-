package v1

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	apimocks "github.com/Anton119/rocket-service-/order/internal/api/order/v1/mocks"
	errs "github.com/Anton119/rocket-service-/order/internal/errors"
	"github.com/Anton119/rocket-service-/order/internal/service/input"
	orderv1 "github.com/Anton119/rocket-service-/shared/pkg/openapi/order/v1"
)

func TestPayOrder(t *testing.T) {
	type expected struct {
		badRequest bool
		notFound   bool
		conflict   bool
	}

	var (
		ctx       = context.Background()
		orderUUID = uuid.MustParse("550e8400-e29b-41d4-a716-446655440099")
		txUUID    = uuid.MustParse("550e8400-e29b-41d4-a716-4466554400aa")

		params = orderv1.PayOrderParams{OrderUUID: orderUUID}

		validReq = &orderv1.PayOrderRequest{}
	)

	validReq.SetPaymentMethod(orderv1.PaymentMethodCARD)

	tests := []struct {
		name      string
		req       *orderv1.PayOrderRequest
		setupMock func(svc *apimocks.OrderService)
		expected  expected
	}{
		{
			name: "успешная оплата",
			req:  validReq,
			setupMock: func(svc *apimocks.OrderService) {
				svc.EXPECT().
					PayOrder(mock.Anything, mock.MatchedBy(func(in input.PayOrderInput) bool {
						return in.OrderUUID == orderUUID
					})).
					Return(txUUID, nil)
			},
		},
		{
			name: "неизвестный способ оплаты",
			req: func() *orderv1.PayOrderRequest {
				r := &orderv1.PayOrderRequest{}
				r.SetPaymentMethod(orderv1.PaymentMethod("UNKNOWN"))
				return r
			}(),
			expected: expected{badRequest: true},
		},
		{
			name: "заказ не найден",
			req:  validReq,
			setupMock: func(svc *apimocks.OrderService) {
				svc.EXPECT().
					PayOrder(mock.Anything, mock.Anything).
					Return(uuid.Nil, errs.ErrOrderNotFound)
			},
			expected: expected{notFound: true},
		},
		{
			name: "оплата невозможна в текущем статусе",
			req:  validReq,
			setupMock: func(svc *apimocks.OrderService) {
				svc.EXPECT().
					PayOrder(mock.Anything, mock.Anything).
					Return(uuid.Nil, errs.ErrOrderPayNotAllowed)
			},
			expected: expected{conflict: true},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := apimocks.NewOrderService(t)
			if tc.setupMock != nil {
				tc.setupMock(svc)
			}

			res, err := NewAPI(svc).PayOrder(ctx, tc.req, params)

			require.NoError(t, err)
			require.NotNil(t, res)

			switch {
			case tc.expected.badRequest:
				resp, ok := res.(*orderv1.PayOrderBadRequest)
				require.True(t, ok)
				assert.Equal(t, http.StatusBadRequest, resp.Code)
			case tc.expected.notFound:
				resp, ok := res.(*orderv1.PayOrderNotFound)
				require.True(t, ok)
				assert.Equal(t, http.StatusNotFound, resp.Code)
			case tc.expected.conflict:
				resp, ok := res.(*orderv1.PayOrderConflict)
				require.True(t, ok)
				assert.Equal(t, http.StatusConflict, resp.Code)
			default:
				resp, ok := res.(*orderv1.PayOrderResponse)
				require.True(t, ok)
				assert.Equal(t, txUUID, resp.TransactionUUID)
			}
		})
	}
}
