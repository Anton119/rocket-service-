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

	apimocks "github.com/Anton119/rocket-service-/inventory/internal/api/inventory/v1/mocks"
	errs "github.com/Anton119/rocket-service-/inventory/internal/errors"
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
			testPart(t, hullUUID, "Hull", "", model.PartTypeHull, 500_000, 10),
			testPart(t, engineUUID, "Engine", "", model.PartTypeEngine, 300_000, 5),
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

			res, err := withErrorInterceptor[*inventoryv1.ListPartsResponse](ctx, tc.req, func(ctx context.Context, req any) (any, error) {
				return NewAPI(svc).ListParts(ctx, req.(*inventoryv1.ListPartsRequest))
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
			assert.Len(t, res.GetParts(), tc.expected.count)
		})
	}
}
