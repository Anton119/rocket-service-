package v1_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	v1 "github.com/Anton119/rocket-service-/inventory/internal/api/inventory/v1"
	apimocks "github.com/Anton119/rocket-service-/inventory/internal/api/inventory/v1/mocks"
	errs "github.com/Anton119/rocket-service-/inventory/internal/errors"
	invinterceptor "github.com/Anton119/rocket-service-/inventory/internal/interceptor"
	"github.com/Anton119/rocket-service-/inventory/internal/model"
	inventoryv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/inventory/v1"
)

func TestGetPart(t *testing.T) {
	type expected struct {
		code codes.Code
	}

	var (
		ctx      = context.Background()
		partUUID = uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")

		storedPart = model.Part{
			UUID:          partUUID.String(),
			Name:          "Hull",
			Description:   "Test hull",
			Price:         500_000,
			PartType:      model.PartTypeHull,
			StockQuantity: 10,
			CreatedAt:     time.Now(),
		}
	)

	tests := []struct {
		name      string
		req       *inventoryv1.GetPartRequest
		setupMock func(svc *apimocks.PartService)
		expected  expected
	}{
		{
			name: "успешное получение",
			req:  &inventoryv1.GetPartRequest{Uuid: partUUID.String()},
			setupMock: func(svc *apimocks.PartService) {
				svc.EXPECT().
					GetPart(ctx, partUUID).
					Return(storedPart, nil)
			},
		},
		{
			name:     "неверный uuid",
			req:      &inventoryv1.GetPartRequest{Uuid: "not-a-uuid"},
			expected: expected{code: codes.InvalidArgument},
		},
		{
			name: "деталь не найдена",
			req:  &inventoryv1.GetPartRequest{Uuid: partUUID.String()},
			setupMock: func(svc *apimocks.PartService) {
				svc.EXPECT().
					GetPart(ctx, partUUID).
					Return(model.Part{}, errs.ErrPartNotFound)
			},
			expected: expected{code: codes.NotFound},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := apimocks.NewPartService(t)
			if tc.setupMock != nil {
				tc.setupMock(svc)
			}

			res, err := v1.NewAPI(svc).GetPart(ctx, tc.req)
			err = invinterceptor.ToGRPCError(err)

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
			assert.Equal(t, partUUID.String(), res.GetPart().GetUuid())
		})
	}
}
