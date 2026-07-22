package v1_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	orderapi "github.com/Anton119/rocket-service-/order/internal/api/order/v1"
	apimocks "github.com/Anton119/rocket-service-/order/internal/api/order/v1/mocks"
	errs "github.com/Anton119/rocket-service-/order/internal/errors"
	orderv1 "github.com/Anton119/rocket-service-/shared/pkg/openapi/order/v1"
)

func TestCancelOrder(t *testing.T) {
	type expected struct {
		notFound        bool
		conflict        bool
		conflictMessage string
	}

	var (
		ctx       = context.Background()
		orderUUID = uuid.MustParse("550e8400-e29b-41d4-a716-446655440099")
		params    = orderv1.CancelOrderParams{OrderUUID: orderUUID}
	)

	tests := []struct {
		name      string
		setupMock func(svc *apimocks.OrderService)
		expected  expected
	}{
		{
			name: "успешная отмена",
			setupMock: func(svc *apimocks.OrderService) {
				svc.EXPECT().
					CancelOrder(ctx, orderUUID).
					Return(nil)
			},
		},
		{
			name: "заказ не найден",
			setupMock: func(svc *apimocks.OrderService) {
				svc.EXPECT().
					CancelOrder(ctx, orderUUID).
					Return(errs.ErrOrderNotFound)
			},
			expected: expected{notFound: true},
		},
		{
			name: "отмена оплаченного заказа",
			setupMock: func(svc *apimocks.OrderService) {
				svc.EXPECT().
					CancelOrder(ctx, orderUUID).
					Return(errs.ErrOrderAlreadyPaid)
			},
			expected: expected{
				conflict:        true,
				conflictMessage: errs.ErrOrderAlreadyPaid.Error(),
			},
		},
		{
			name: "повторная отмена заказа",
			setupMock: func(svc *apimocks.OrderService) {
				svc.EXPECT().
					CancelOrder(ctx, orderUUID).
					Return(errs.ErrOrderAlreadyCancelled)
			},
			expected: expected{
				conflict:        true,
				conflictMessage: errs.ErrOrderAlreadyCancelled.Error(),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := apimocks.NewOrderService(t)
			tc.setupMock(svc)

			res, err := orderapi.NewAPI(svc).CancelOrder(ctx, params)

			require.NoError(t, err)
			require.NotNil(t, res)

			switch {
			case tc.expected.notFound:
				resp, ok := res.(*orderv1.CancelOrderNotFound)
				require.True(t, ok)
				assert.Equal(t, http.StatusNotFound, resp.Code)
			case tc.expected.conflict:
				resp, ok := res.(*orderv1.CancelOrderConflict)
				require.True(t, ok)
				assert.Equal(t, http.StatusConflict, resp.Code)
				assert.Equal(t, tc.expected.conflictMessage, resp.Message)
			default:
				_, ok := res.(*orderv1.CancelOrderResponse)
				require.True(t, ok)
			}
		})
	}
}
