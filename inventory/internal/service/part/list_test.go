package part

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errs "github.com/Anton119/rocket-service-/inventory/internal/errors"
	"github.com/Anton119/rocket-service-/inventory/internal/model"
	"github.com/Anton119/rocket-service-/inventory/internal/service/input"
	"github.com/Anton119/rocket-service-/inventory/internal/service/part/mocks"
)

func TestListParts(t *testing.T) {
	type args struct {
		in input.ListPartsInput
	}

	type expected struct {
		err   error
		parts []model.Part
	}

	var (
		ctx = context.Background()

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
		args      args
		setupMock func(repo *mocks.PartRepository)
		expected  expected
	}{
		{
			name: "успешный список по UUID",
			args: args{
				in: input.ListPartsInput{
					IDs: []uuid.UUID{hullUUID, engineUUID},
				},
			},
			setupMock: func(repo *mocks.PartRepository) {
				repo.EXPECT().
					ListParts(ctx, model.PartTypeUnspecified, []uuid.UUID{hullUUID, engineUUID}).
					Return(parts, nil)
			},
			expected: expected{parts: parts},
		},
		{
			name: "деталь не найдена",
			args: args{
				in: input.ListPartsInput{
					IDs: []uuid.UUID{hullUUID},
				},
			},
			setupMock: func(repo *mocks.PartRepository) {
				repo.EXPECT().
					ListParts(ctx, model.PartTypeUnspecified, []uuid.UUID{hullUUID}).
					Return(nil, errs.ErrPartNotFound)
			},
			expected: expected{err: errs.ErrPartNotFound},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewPartRepository(t)
			tc.setupMock(repo)

			svc := NewService(repo)
			got, err := svc.ListParts(ctx, tc.args.in)

			if tc.expected.err != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.err)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected.parts, got)
			}
		})
	}
}
