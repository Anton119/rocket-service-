package order_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	errs "github.com/Anton119/rocket-service-/order/internal/errors"
	"github.com/Anton119/rocket-service-/order/internal/model"
	orderproducer "github.com/Anton119/rocket-service-/order/internal/producer/order_producer"
	ordersvc "github.com/Anton119/rocket-service-/order/internal/service/order"
	"github.com/Anton119/rocket-service-/order/internal/service/order/mocks"
)

func TestCancelOrder(t *testing.T) {
	type expected struct {
		err error
	}

	var (
		ctx       = context.Background()
		orderUUID = uuid.MustParse("550e8400-e29b-41d4-a716-446655440099")

		pendingOrder = model.Order{
			OrderUUID:  orderUUID,
			HullUUID:   uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"),
			EngineUUID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440003"),
			TotalPrice: 800_000,
			Status:     model.OrderStatusPendingPayment,
			CreatedAt:  time.Now(),
		}

		paidOrder = model.Order{
			OrderUUID:  orderUUID,
			HullUUID:   pendingOrder.HullUUID,
			EngineUUID: pendingOrder.EngineUUID,
			TotalPrice: 800_000,
			Status:     model.OrderStatusPaid,
			CreatedAt:  pendingOrder.CreatedAt,
		}

		cancelledOrder = model.Order{
			OrderUUID:  orderUUID,
			HullUUID:   pendingOrder.HullUUID,
			EngineUUID: pendingOrder.EngineUUID,
			TotalPrice: 800_000,
			Status:     model.OrderStatusCancelled,
			CreatedAt:  pendingOrder.CreatedAt,
		}
	)

	tests := []struct {
		name      string
		setupMock func(repo *mocks.OrderRepository, inv *mocks.InventoryClient)
		expected  expected
	}{
		{
			name: "успешная отмена",
			setupMock: func(repo *mocks.OrderRepository, inv *mocks.InventoryClient) {
				repo.EXPECT().GetForUpdate(ctx, orderUUID).Return(pendingOrder, nil)
				repo.EXPECT().
					Save(ctx, mock.MatchedBy(func(o model.Order) bool {
						return o.OrderUUID == orderUUID && o.Status == model.OrderStatusCancelled
					})).
					Return(nil)
				inv.EXPECT().
					ReleaseParts(ctx, []string{pendingOrder.HullUUID.String(), pendingOrder.EngineUUID.String()}).
					Return(nil)
			},
			expected: expected{err: nil},
		},
		{
			name: "заказ не найден",
			setupMock: func(repo *mocks.OrderRepository, inv *mocks.InventoryClient) {
				repo.EXPECT().
					GetForUpdate(ctx, orderUUID).
					Return(model.Order{}, errs.ErrOrderNotFound)
			},
			expected: expected{err: errs.ErrOrderNotFound},
		},
		{
			name: "отмена оплаченного заказа запрещена",
			setupMock: func(repo *mocks.OrderRepository, inv *mocks.InventoryClient) {
				repo.EXPECT().
					GetForUpdate(ctx, orderUUID).
					Return(paidOrder, nil)
			},
			expected: expected{err: errs.ErrOrderAlreadyPaid},
		},
		{
			name: "повторная отмена заказа",
			setupMock: func(repo *mocks.OrderRepository, inv *mocks.InventoryClient) {
				repo.EXPECT().
					GetForUpdate(ctx, orderUUID).
					Return(cancelledOrder, nil)
			},
			expected: expected{err: errs.ErrOrderAlreadyCancelled},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewOrderRepository(t)
			inv := mocks.NewInventoryClient(t)
			pay := mocks.NewPaymentClient(t)

			tc.setupMock(repo, inv)

			svc := ordersvc.NewService(repo, inv, pay, orderproducer.NoopProducer{}, passthroughTx{})
			err := svc.CancelOrder(ctx, orderUUID)

			if tc.expected.err != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
