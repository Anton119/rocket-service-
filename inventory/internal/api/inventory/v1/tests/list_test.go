package v1_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	v1 "github.com/Anton119/rocket-service-/inventory/internal/api/inventory/v1"
	apimocks "github.com/Anton119/rocket-service-/inventory/internal/api/inventory/v1/mocks"
	errs "github.com/Anton119/rocket-service-/inventory/internal/errors"
	invinterceptor "github.com/Anton119/rocket-service-/inventory/internal/interceptor"
	"github.com/Anton119/rocket-service-/inventory/internal/model"
	"github.com/Anton119/rocket-service-/inventory/internal/service/input"
	inventoryv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/inventory/v1"
)

func TestListParts(t *testing.T) {
	type expected struct {
		code  codes.Code
		count int
	}

	var (
		ctx        = context.Background()
		hullUUID   = uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
		engineUUID = uuid.MustParse("550e8400-e29b-41d4-a716-446655440003")

		parts = []model.Part{
			{
				UUID:          hullUUID.String(),
				Name:          "Hull",
				Price:         500_000,
				PartType:      model.PartTypeHull,
				StockQuantity: 10,
				CreatedAt:     time.Now(),
			},
			{
				UUID:          engineUUID.String(),
				Name:          "Engine",
				Price:         300_000,
				PartType:      model.PartTypeEngine,
				StockQuantity: 5,
				CreatedAt:     time.Now(),
			},
		}
	)

	tests := []struct {
		name      string
		req       *inventoryv1.ListPartsRequest
		setupMock func(svc *apimocks.PartService)
		expected  expected
	}{
		{
			name: "успешный список",
			req: &inventoryv1.ListPartsRequest{
				Uuids: []string{hullUUID.String(), engineUUID.String()},
			},
			setupMock: func(svc *apimocks.PartService) {
				svc.EXPECT().
					ListParts(ctx, mock.MatchedBy(func(in input.ListPartsInput) bool {
						return len(in.IDs) == 2
					})).
					Return(parts, nil)
			},
			expected: expected{count: 2},
		},
		{
			name: "неверный uuid в списке",
			req: &inventoryv1.ListPartsRequest{
				Uuids: []string{"bad-uuid"},
			},
			expected: expected{code: codes.InvalidArgument},
		},
		{
			name: "деталь не найдена",
			req: &inventoryv1.ListPartsRequest{
				Uuids: []string{hullUUID.String()},
			},
			setupMock: func(svc *apimocks.PartService) {
				svc.EXPECT().
					ListParts(ctx, mock.Anything).
					Return(nil, errs.ErrPartNotFound)
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

			res, err := v1.NewAPI(svc).ListParts(ctx, tc.req)
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
			assert.Len(t, res.GetParts(), tc.expected.count)
		})
	}
}
