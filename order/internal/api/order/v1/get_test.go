package v1

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apimocks "github.com/Anton119/rocket-service-/order/internal/api/order/v1/mocks"
	errs "github.com/Anton119/rocket-service-/order/internal/errors"
	"github.com/Anton119/rocket-service-/order/internal/model"
	orderv1 "github.com/Anton119/rocket-service-/shared/pkg/openapi/order/v1"
)

func TestGetOrder(t *testing.T) {
	type expected struct {
		notFound bool
	}

	var (
		ctx       = context.Background()
		orderUUID  = uuid.MustParse("550e8400-e29b-41d4-a716-446655440099")
		totalPrice = int64(800_000)

		storedOrder = model.Order{
			OrderUUID:  orderUUID,
			HullUUID:   uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"),
			EngineUUID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440003"),
			TotalPrice: totalPrice,
			Status:     model.OrderStatusPendingPayment,
			CreatedAt:  time.Now(),
		}

		params = orderv1.GetOrderParams{OrderUUID: orderUUID}
	)

	tests := []struct {
		name      string
		setupMock func(svc *apimocks.OrderService)
		expected  expected
	}{
		{
			name: "успешное получение",
			setupMock: func(svc *apimocks.OrderService) {
				svc.EXPECT().
					GetOrder(ctx, orderUUID).
					Return(storedOrder, nil)
			},
		},
		{
			name: "заказ не найден",
			setupMock: func(svc *apimocks.OrderService) {
				svc.EXPECT().
					GetOrder(ctx, orderUUID).
					Return(model.Order{}, errs.ErrOrderNotFound)
			},
			expected: expected{notFound: true},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := apimocks.NewOrderService(t)
			tc.setupMock(svc)

			res, err := NewAPI(svc).GetOrder(ctx, params)

			require.NoError(t, err)
			require.NotNil(t, res)

			if tc.expected.notFound {
				resp, ok := res.(*orderv1.GetOrderNotFound)
				require.True(t, ok)
				assert.Equal(t, http.StatusNotFound, resp.Code)
				return
			}

			dto, ok := res.(*orderv1.OrderDto)
			require.True(t, ok)
			assert.Equal(t, orderUUID, dto.OrderUUID)
			assert.Equal(t, totalPrice, dto.TotalPrice)
		})
	}
}
