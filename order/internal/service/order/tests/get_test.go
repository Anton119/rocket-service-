package order_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errs "github.com/Anton119/rocket-service-/order/internal/errors"
	"github.com/Anton119/rocket-service-/order/internal/model"
	orderproducer "github.com/Anton119/rocket-service-/order/internal/producer/order_producer"
	ordersvc "github.com/Anton119/rocket-service-/order/internal/service/order"
	"github.com/Anton119/rocket-service-/order/internal/service/order/mocks"
)

func TestGetOrder(t *testing.T) {
	type expected struct {
		err   error
		order model.Order
	}

	var (
		ctx        = context.Background()
		orderUUID  = uuid.MustParse("550e8400-e29b-41d4-a716-446655440099")
		hullUUID   = uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
		engineUUID = uuid.MustParse("550e8400-e29b-41d4-a716-446655440003")

		storedOrder = model.Order{
			OrderUUID:  orderUUID,
			HullUUID:   hullUUID,
			EngineUUID: engineUUID,
			TotalPrice: 800_000,
			Status:     model.OrderStatusPendingPayment,
			CreatedAt:  time.Now(),
		}
	)

	tests := []struct {
		name      string
		setupMock func(repo *mocks.OrderRepository)
		expected  expected
	}{
		{
			name: "успешное получение заказа",
			setupMock: func(repo *mocks.OrderRepository) {
				repo.EXPECT().
					Get(ctx, orderUUID).
					Return(storedOrder, nil)
			},
			expected: expected{err: nil, order: storedOrder},
		},
		{
			name: "заказ не найден",
			setupMock: func(repo *mocks.OrderRepository) {
				repo.EXPECT().
					Get(ctx, orderUUID).
					Return(model.Order{}, errs.ErrOrderNotFound)
			},
			expected: expected{err: errs.ErrOrderNotFound},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewOrderRepository(t)
			inv := mocks.NewInventoryClient(t)
			pay := mocks.NewPaymentClient(t)

			tc.setupMock(repo)

			svc := ordersvc.NewService(repo, inv, pay, orderproducer.NoopProducer{}, passthroughTx{})
			order, err := svc.GetOrder(ctx, orderUUID)

			if tc.expected.err != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.err)
				assert.Equal(t, uuid.Nil, order.OrderUUID)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected.order, order)
			}
		})
	}
}
