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
	"github.com/Anton119/rocket-service-/inventory/internal/service/part/mocks"
)

func TestGetPart(t *testing.T) {
	type expected struct {
		err  error
		part model.Part
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
		setupMock func(repo *mocks.PartRepository)
		expected  expected
	}{
		{
			name: "успешное получение детали",
			setupMock: func(repo *mocks.PartRepository) {
				repo.EXPECT().
					Get(ctx, partUUID).
					Return(storedPart, nil)
			},
			expected: expected{err: nil, part: storedPart},
		},
		{
			name: "деталь не найдена",
			setupMock: func(repo *mocks.PartRepository) {
				repo.EXPECT().
					Get(ctx, partUUID).
					Return(model.Part{}, errs.ErrPartNotFound)
			},
			expected: expected{err: errs.ErrPartNotFound},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewPartRepository(t)
			tc.setupMock(repo)

			svc := NewService(repo)
			part, err := svc.GetPart(ctx, partUUID)

			if tc.expected.err != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.err)
				assert.Empty(t, part.UUID)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected.part, part)
			}
		})
	}
}
