package order_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	errs "github.com/Anton119/rocket-service-/order/internal/errors"
	"github.com/Anton119/rocket-service-/order/internal/model"
	orderproducer "github.com/Anton119/rocket-service-/order/internal/producer/order_producer"
	"github.com/Anton119/rocket-service-/order/internal/service/input"
	ordersvc "github.com/Anton119/rocket-service-/order/internal/service/order"
	"github.com/Anton119/rocket-service-/order/internal/service/order/mocks"
)

func TestCreateOrder(t *testing.T) {
	type args struct {
		in input.CreateOrderInput
	}

	type expected struct {
		err error
	}

	var (
		ctx = context.Background()

		hullUUID   = uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
		engineUUID = uuid.MustParse("550e8400-e29b-41d4-a716-446655440003")
		userUUID   = uuid.MustParse("550e8400-e29b-41d4-a716-446655440010")

		partsInStock = []model.Part{
			{UUID: hullUUID.String(), Price: 500_000, StockQuantity: 10},
			{UUID: engineUUID.String(), Price: 300_000, StockQuantity: 5},
		}

		partsOutOfStock = []model.Part{
			{UUID: hullUUID.String(), Price: 500_000, StockQuantity: 10},
			{UUID: engineUUID.String(), Price: 300_000, StockQuantity: 0},
		}
	)

	tests := []struct {
		name      string
		args      args
		setupMock func(repo *mocks.OrderRepository, inv *mocks.InventoryClient)
		expected  expected
	}{
		{
			name: "успешное создание заказа",
			args: args{
				in: input.CreateOrderInput{
					UserUUID:   userUUID,
					HullUUID:   hullUUID,
					EngineUUID: engineUUID,
				},
			},
			setupMock: func(repo *mocks.OrderRepository, inv *mocks.InventoryClient) {
				inv.EXPECT().
					ListParts(ctx, []string{hullUUID.String(), engineUUID.String()}).
					Return(partsInStock, nil)

				inv.EXPECT().
					ReserveParts(ctx, []string{hullUUID.String(), engineUUID.String()}).
					Return(nil)

				repo.EXPECT().
					Create(ctx, mock.MatchedBy(func(o model.Order) bool {
						return o.UserUUID == userUUID &&
							o.HullUUID == hullUUID &&
							o.EngineUUID == engineUUID &&
							o.TotalPrice == 800_000 &&
							o.Status == model.OrderStatusPendingPayment &&
							o.OrderUUID != uuid.Nil &&
							len(o.Items) == 2
					})).
					Return(nil)
			},
			expected: expected{err: nil},
		},
		{
			name: "деталь не найдена в inventory",
			args: args{
				in: input.CreateOrderInput{
					UserUUID:   userUUID,
					HullUUID:   hullUUID,
					EngineUUID: engineUUID,
				},
			},
			setupMock: func(_ *mocks.OrderRepository, inv *mocks.InventoryClient) {
				inv.EXPECT().
					ListParts(ctx, []string{hullUUID.String(), engineUUID.String()}).
					Return(nil, errs.ErrPartNotFound)
			},
			expected: expected{err: errs.ErrPartNotFound},
		},
		{
			name: "деталь отсутствует на складе",
			args: args{
				in: input.CreateOrderInput{
					UserUUID:   userUUID,
					HullUUID:   hullUUID,
					EngineUUID: engineUUID,
				},
			},
			setupMock: func(_ *mocks.OrderRepository, inv *mocks.InventoryClient) {
				inv.EXPECT().
					ListParts(ctx, []string{hullUUID.String(), engineUUID.String()}).
					Return(partsOutOfStock, nil)
			},
			expected: expected{err: errs.ErrPartOutOfStock},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewOrderRepository(t)
			inv := mocks.NewInventoryClient(t)
			pay := mocks.NewPaymentClient(t)

			tc.setupMock(repo, inv)

			svc := ordersvc.NewService(repo, inv, pay, orderproducer.NoopProducer{}, passthroughTx{})
			out, err := svc.CreateOrder(ctx, tc.args.in)

			if tc.expected.err != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.err)
				assert.Nil(t, out)
			} else {
				require.NoError(t, err)
				require.NotNil(t, out)
				assert.NotEqual(t, uuid.Nil, out.OrderUUID)
				assert.Equal(t, int64(800_000), out.TotalPrice)
			}
		})
	}
}
