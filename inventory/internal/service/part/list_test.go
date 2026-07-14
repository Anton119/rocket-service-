package part

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errs "github.com/Anton119/rocket-service-/inventory/internal/errors"
	"github.com/Anton119/rocket-service-/inventory/internal/model"
	"github.com/Anton119/rocket-service-/inventory/internal/service/domain"
	"github.com/Anton119/rocket-service-/inventory/internal/service/input"
	"github.com/Anton119/rocket-service-/inventory/internal/service/part/mocks"
)

func TestListParts(t *testing.T) {
	ctx := context.Background()

	hullUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
	engineUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440003")

	parts := []model.Part{
		testPart(t, hullUUID, "Hull", model.PartTypeHull, 10),
		testPart(t, engineUUID, "Engine", model.PartTypeEngine, 5),
	}

	tests := []struct {
		name      string
		in        input.ListPartsInput
		setupMock func(repo *mocks.PartRepository)
		wantErr   error
	}{
		{
			name: "успешный список по UUID",
			in:   input.ListPartsInput{IDs: []uuid.UUID{hullUUID, engineUUID}},
			setupMock: func(repo *mocks.PartRepository) {
				repo.EXPECT().List(ctx, input.PartFilter{UUIDs: []uuid.UUID{hullUUID, engineUUID}}).Return(parts, nil)
			},
		},
		{
			name: "деталь не найдена",
			in:   input.ListPartsInput{IDs: []uuid.UUID{hullUUID}},
			setupMock: func(repo *mocks.PartRepository) {
				repo.EXPECT().List(ctx, input.PartFilter{UUIDs: []uuid.UUID{hullUUID}}).Return(nil, errs.ErrPartNotFound)
			},
			wantErr: errs.ErrPartNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewPartRepository(t)
			tc.setupMock(repo)

			svc := NewService(repo, noopTxManager{}, domain.NewCompatibilityChecker())
			got, err := svc.ListParts(ctx, tc.in)

			if tc.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
				require.Len(t, got, len(parts))
			}
		})
	}
}
