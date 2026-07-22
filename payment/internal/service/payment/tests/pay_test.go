package payment_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Anton119/rocket-service-/payment/internal/model"
	"github.com/Anton119/rocket-service-/payment/internal/service/input"
	paymentsvc "github.com/Anton119/rocket-service-/payment/internal/service/payment"
)

func TestPayOrder(t *testing.T) {
	type args struct {
		in input.PayOrderInput
	}

	var (
		ctx       = context.Background()
		orderUUID = uuid.MustParse("550e8400-e29b-41d4-a716-446655440099")
	)

	tests := []struct {
		name string
		args args
	}{
		{
			name: "успешная оплата",
			args: args{
				in: input.PayOrderInput{
					OrderUUID: orderUUID,
					Method:    model.PaymentMethodCard,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := paymentsvc.NewService()
			txID, err := svc.PayOrder(ctx, tc.args.in)

			require.NoError(t, err)
			assert.NotEqual(t, uuid.Nil, txID)
		})
	}
}
