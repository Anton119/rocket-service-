package v1

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	apimocks "github.com/Anton119/rocket-service-/payment/internal/api/payment/v1/mocks"
	"github.com/Anton119/rocket-service-/payment/internal/model"
	"github.com/Anton119/rocket-service-/payment/internal/service/input"
	paymentv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/payment/v1"
)

func TestPayOrder(t *testing.T) {
	type expected struct {
		code codes.Code
	}

	var (
		ctx       = context.Background()
		orderUUID = uuid.MustParse("550e8400-e29b-41d4-a716-446655440099")
		txUUID    = uuid.MustParse("550e8400-e29b-41d4-a716-4466554400aa")
	)

	tests := []struct {
		name      string
		req       *paymentv1.PayOrderRequest
		setupMock func(svc *apimocks.PaymentService)
		expected  expected
	}{
		{
			name: "успешная оплата",
			req: &paymentv1.PayOrderRequest{
				OrderUuid:     orderUUID.String(),
				PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_CARD,
			},
			setupMock: func(svc *apimocks.PaymentService) {
				svc.EXPECT().
					PayOrder(ctx, mock.MatchedBy(func(in input.PayOrderInput) bool {
						return in.OrderUUID == orderUUID && in.Method == model.PaymentMethodCard
					})).
					Return(&input.PayOrderResult{TransactionUUID: txUUID}, nil)
			},
		},
		{
			name: "пустой order_uuid",
			req: &paymentv1.PayOrderRequest{
				PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_CARD,
			},
			expected: expected{code: codes.InvalidArgument},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := apimocks.NewPaymentService(t)
			if tc.setupMock != nil {
				tc.setupMock(svc)
			}

			res, err := NewAPI(svc).PayOrder(ctx, tc.req)

			if tc.expected.code != codes.OK {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, tc.expected.code, st.Code())
				assert.Nil(t, res)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, res)
			assert.Equal(t, txUUID.String(), res.GetTransactionUuid())
		})
	}
}
