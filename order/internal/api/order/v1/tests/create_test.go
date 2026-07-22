package v1_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	orderapi "github.com/Anton119/rocket-service-/order/internal/api/order/v1"
	apimocks "github.com/Anton119/rocket-service-/order/internal/api/order/v1/mocks"
	errs "github.com/Anton119/rocket-service-/order/internal/errors"
	"github.com/Anton119/rocket-service-/order/internal/service/input"
	orderv1 "github.com/Anton119/rocket-service-/shared/pkg/openapi/order/v1"
)

func TestCreateOrder(t *testing.T) {
	type expected struct {
		notFound bool
		conflict bool
	}

	var (
		ctx = context.Background()

		hullUUID   = uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
		engineUUID = uuid.MustParse("550e8400-e29b-41d4-a716-446655440003")
		userUUID   = uuid.MustParse("550e8400-e29b-41d4-a716-446655440010")
		orderUUID  = uuid.MustParse("550e8400-e29b-41d4-a716-446655440099")
		totalPrice = int64(800_000)

		req = &orderv1.CreateOrderRequest{}
	)

	req.SetUserUUID(userUUID)
	req.SetHullUUID(hullUUID)
	req.SetEngineUUID(engineUUID)

	tests := []struct {
		name      string
		setupMock func(svc *apimocks.OrderService)
		expected  expected
	}{
		{
			name: "успешное создание",
			setupMock: func(svc *apimocks.OrderService) {
				svc.EXPECT().
					CreateOrder(mock.Anything, input.CreateOrderInput{
						UserUUID:   userUUID,
						HullUUID:   hullUUID,
						EngineUUID: engineUUID,
					}).
					Return(&input.CreateOrderResult{
						OrderUUID:  orderUUID,
						TotalPrice: totalPrice,
					}, nil)
			},
			expected: expected{},
		},
		{
			name: "деталь не найдена",
			setupMock: func(svc *apimocks.OrderService) {
				svc.EXPECT().
					CreateOrder(mock.Anything, mock.Anything).
					Return(nil, errs.ErrPartNotFound)
			},
			expected: expected{notFound: true},
		},
		{
			name: "деталь отсутствует на складе",
			setupMock: func(svc *apimocks.OrderService) {
				svc.EXPECT().
					CreateOrder(mock.Anything, mock.Anything).
					Return(nil, errs.ErrPartOutOfStock)
			},
			expected: expected{conflict: true},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := apimocks.NewOrderService(t)
			tc.setupMock(svc)

			api := orderapi.NewAPI(svc)
			res, err := api.CreateOrder(ctx, req)

			require.NoError(t, err)
			require.NotNil(t, res)

			switch {
			case tc.expected.notFound:
				resp, ok := res.(*orderv1.CreateOrderNotFound)
				require.True(t, ok)
				assert.Equal(t, http.StatusNotFound, resp.Code)
			case tc.expected.conflict:
				resp, ok := res.(*orderv1.CreateOrderConflict)
				require.True(t, ok)
				assert.Equal(t, http.StatusConflict, resp.Code)
			default:
				resp, ok := res.(*orderv1.CreateOrderResponse)
				require.True(t, ok)
				assert.Equal(t, orderUUID, resp.OrderUUID)
				assert.Equal(t, totalPrice, resp.TotalPrice)
			}
		})
	}
}
