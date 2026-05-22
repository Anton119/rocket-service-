package part

import (
	"context"
	"sync"

	"github.com/google/uuid"

	errs "github.com/Anton119/rocket-service-/inventory/internal/errors"
	"github.com/Anton119/rocket-service-/inventory/internal/model"
	repoconv "github.com/Anton119/rocket-service-/inventory/internal/repository/converter"
	"github.com/Anton119/rocket-service-/inventory/internal/repository/record"
)

// Repository — потокобезопасное in-memory хранилище деталей.
type Repository struct {
	mu    sync.RWMutex
	parts map[uuid.UUID]record.Part
}

// NewRepository создаёт репозиторий из доменных сущностей (конвертирует в record).
func NewRepository(parts map[uuid.UUID]model.Part) *Repository {
	rec := make(map[uuid.UUID]record.Part, len(parts))
	for id, m := range parts {
		r := repoconv.PartModelToRecord(m)
		rec[id] = r
	}

	return &Repository{parts: rec}
}

func (r *Repository) Get(_ context.Context, id uuid.UUID) (model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, ok := r.parts[id]
	if !ok {
		return model.Part{}, errs.ErrPartNotFound
	}

	return repoconv.PartRecordToModel(p), nil
}

// ListParts возвращает детали по списку ids (если ids непустой — порядок как в ids).
// Если ids пустой — выборка по partType и сортировка по имени.
func (r *Repository) ListParts(_ context.Context, partType model.PartType, ids []uuid.UUID) ([]model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(ids) > 0 {
		return listPartsByOrderedIDs(r.parts, ids)
	}

	list := filterPartsByType(r.parts, partType)
	sortRecordsByName(list)

	out := make([]model.Part, 0, len(list))
	for i := range list {
		out = append(out, repoconv.PartRecordToModel(list[i]))
	}

	return out, nil
}
