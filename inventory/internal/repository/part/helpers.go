package part

import (
	"sort"

	"github.com/google/uuid"

	errs "github.com/Anton119/rocket-service-/inventory/internal/errors"
	"github.com/Anton119/rocket-service-/inventory/internal/model"
	"github.com/Anton119/rocket-service-/inventory/internal/repository/converter"
	"github.com/Anton119/rocket-service-/inventory/internal/repository/record"
)

func listPartsByOrderedIDs(parts map[uuid.UUID]record.Part, ids []uuid.UUID) ([]model.Part, error) {
	out := make([]model.Part, 0, len(ids))
	for _, id := range ids {
		p, ok := parts[id]
		if !ok {
			return nil, errs.ErrPartNotFound
		}
		out = append(out, converter.PartRecordToModel(p))
	}
	return out, nil
}

func filterPartsByType(parts map[uuid.UUID]record.Part, partType model.PartType) []record.Part {
	list := make([]record.Part, 0, len(parts))
	for _, p := range parts {
		if partType != model.PartTypeUnspecified && model.PartType(p.PartType) != partType {
			continue
		}
		list = append(list, p)
	}
	return list
}

func sortRecordsByName(list []record.Part) {
	sort.Slice(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})
}
