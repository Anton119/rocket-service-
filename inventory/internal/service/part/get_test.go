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
	"github.com/Anton119/rocket-service-/inventory/internal/service/domain"
	"github.com/Anton119/rocket-service-/inventory/internal/service/part/mocks"
)

type noopTxManager struct{}

func (noopTxManager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func testPart(t *testing.T, id uuid.UUID, name string, pt model.PartType, stock int) model.Part {
	t.Helper()

	props, err := model.NewHullProperties(50)
	require.NoError(t, err)
	if pt == model.PartTypeEngine {
		props, err = model.NewEngineProperties(model.EngineClassC, 30)
		require.NoError(t, err)
	}

	return model.RestorePart(id, name, "desc", pt, 100, stock, 0, props, time.Now())
}

func TestGetPart(t *testing.T) {
	ctx := context.Background()
	partUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
	storedPart := testPart(t, partUUID, "Hull", model.PartTypeHull, 10)

	tests := []struct {
		name      string
		setupMock func(repo *mocks.PartRepository)
		wantErr   error
	}{
		{
			name: "успешное получение детали",
			setupMock: func(repo *mocks.PartRepository) {
				repo.EXPECT().Get(ctx, partUUID).Return(storedPart, nil)
			},
		},
		{
			name: "деталь не найдена",
			setupMock: func(repo *mocks.PartRepository) {
				repo.EXPECT().Get(ctx, partUUID).Return(model.Part{}, errs.ErrPartNotFound)
			},
			wantErr: errs.ErrPartNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewPartRepository(t)
			tc.setupMock(repo)

			svc := NewService(repo, noopTxManager{}, domain.NewCompatibilityChecker())
			part, err := svc.GetPart(ctx, partUUID)

			if tc.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, storedPart.UUID(), part.UUID())
				assert.Equal(t, storedPart.Name(), part.Name())
			}
		})
	}
}
