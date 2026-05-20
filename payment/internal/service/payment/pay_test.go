package payment

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Anton119/rocket-service-/payment/internal/model"
	"github.com/Anton119/rocket-service-/payment/internal/service/input"
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
			svc := NewService()
			out, err := svc.PayOrder(ctx, tc.args.in)

			require.NoError(t, err)
			require.NotNil(t, out)
			assert.NotEqual(t, uuid.Nil, out.TransactionUUID)
		})
	}
}
