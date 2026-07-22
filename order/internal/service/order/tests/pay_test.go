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
	"github.com/Anton119/rocket-service-/order/internal/service/input"
	ordersvc "github.com/Anton119/rocket-service-/order/internal/service/order"
	"github.com/Anton119/rocket-service-/order/internal/service/order/mocks"
)

func TestPayOrder(t *testing.T) {
	type args struct {
		in input.PayOrderInput
	}

	type expected struct {
		err error
	}

	var (
		ctx = context.Background()

		orderUUID = uuid.MustParse("550e8400-e29b-41d4-a716-446655440099")
		txUUID    = uuid.MustParse("550e8400-e29b-41d4-a716-4466554400aa")

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
	)

	tests := []struct {
		name      string
		args      args
		setupMock func(repo *mocks.OrderRepository, pay *mocks.PaymentClient, paid *mocks.OrderPaidProducer)
		expected  expected
	}{
		{
			name: "успешная оплата",
			args: args{
				in: input.PayOrderInput{
					OrderUUID:           orderUUID,
					Method:              model.PaymentMethodCard,
					PaymentMethodStored: "CARD",
				},
			},
			setupMock: func(repo *mocks.OrderRepository, pay *mocks.PaymentClient, paid *mocks.OrderPaidProducer) {
				repo.EXPECT().GetForUpdate(ctx, orderUUID).Return(pendingOrder, nil)
				pay.EXPECT().
					PayOrder(ctx, orderUUID.String(), model.PaymentMethodCard).
					Return(txUUID, nil)
				repo.EXPECT().
					Save(ctx, mock.MatchedBy(func(o model.Order) bool {
						return o.OrderUUID == orderUUID &&
							o.Status == model.OrderStatusPaid &&
							o.TransactionUUID != nil &&
							*o.TransactionUUID == txUUID
					})).
					Return(nil)
				paid.EXPECT().
					Produce(ctx, mock.MatchedBy(func(e model.OrderPaidEvent) bool {
						return e.OrderUUID == orderUUID.String() &&
							e.TransactionUUID == txUUID.String() &&
							e.PaymentMethod == "CARD"
					})).
					Return(nil)
			},
			expected: expected{err: nil},
		},
		{
			name: "заказ не найден",
			args: args{
				in: input.PayOrderInput{
					OrderUUID:           orderUUID,
					Method:              model.PaymentMethodCard,
					PaymentMethodStored: "CARD",
				},
			},
			setupMock: func(repo *mocks.OrderRepository, _ *mocks.PaymentClient, _ *mocks.OrderPaidProducer) {
				repo.EXPECT().
					GetForUpdate(ctx, orderUUID).
					Return(model.Order{}, errs.ErrOrderNotFound)
			},
			expected: expected{err: errs.ErrOrderNotFound},
		},
		{
			name: "оплата уже оплаченного заказа",
			args: args{
				in: input.PayOrderInput{
					OrderUUID:           orderUUID,
					Method:              model.PaymentMethodCard,
					PaymentMethodStored: "CARD",
				},
			},
			setupMock: func(repo *mocks.OrderRepository, _ *mocks.PaymentClient, _ *mocks.OrderPaidProducer) {
				repo.EXPECT().
					GetForUpdate(ctx, orderUUID).
					Return(paidOrder, nil)
			},
			expected: expected{err: errs.ErrOrderPayNotAllowed},
		},
		{
			name: "невалидный способ оплаты",
			args: args{
				in: input.PayOrderInput{
					OrderUUID:           orderUUID,
					Method:              model.PaymentMethodInvalid,
					PaymentMethodStored: "",
				},
			},
			setupMock: func(_ *mocks.OrderRepository, _ *mocks.PaymentClient, _ *mocks.OrderPaidProducer) {},
			expected:  expected{err: errs.ErrInvalidPaymentMethod},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewOrderRepository(t)
			inv := mocks.NewInventoryClient(t)
			pay := mocks.NewPaymentClient(t)
			paid := mocks.NewOrderPaidProducer(t)

			tc.setupMock(repo, pay, paid)

			svc := ordersvc.NewService(repo, inv, pay, paid, passthroughTx{})
			txID, err := svc.PayOrder(ctx, tc.args.in)

			if tc.expected.err != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.err)
				assert.Equal(t, uuid.Nil, txID)
			} else {
				require.NoError(t, err)
				assert.Equal(t, txUUID, txID)
			}
		})
	}
}
