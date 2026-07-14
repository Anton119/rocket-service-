package v1

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	apimocks "github.com/Anton119/rocket-service-/inventory/internal/api/inventory/v1/mocks"
	errs "github.com/Anton119/rocket-service-/inventory/internal/errors"
	"github.com/Anton119/rocket-service-/inventory/internal/model"
	inventoryv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/inventory/v1"
)

func testPart(t *testing.T, id uuid.UUID, name, desc string, pt model.PartType, price int64, stock int) model.Part {
	t.Helper()

	var props model.PartProperties
	var err error

	switch pt {
	case model.PartTypeHull:
		props, err = model.NewHullProperties(50)
	case model.PartTypeEngine:
		props, err = model.NewEngineProperties(model.EngineClassC, 30)
	default:
		props, err = model.NewHullProperties(50)
	}
	require.NoError(t, err)

	return model.RestorePart(id, name, desc, pt, price, stock, 0, props, time.Now())
}

func TestGetPart(t *testing.T) {
	type expected struct {
		code codes.Code
	}

	var (
		ctx      = context.Background()
		partUUID = uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")

		storedPart = testPart(t, partUUID, "Hull", "Test hull", model.PartTypeHull, 500_000, 10)
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

			res, err := withErrorInterceptor[*inventoryv1.GetPartResponse](ctx, tc.req, func(ctx context.Context, req any) (any, error) {
				return NewAPI(svc).GetPart(ctx, req.(*inventoryv1.GetPartRequest))
			})

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
