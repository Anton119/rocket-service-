package converter

import (
	"errors"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Anton119/rocket-service-/inventory/internal/model"
	"github.com/Anton119/rocket-service-/inventory/internal/service/input"
	inventoryv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/inventory/v1"
)

var (
	errEmptyUUID   = errors.New("uuid не может быть пустым")
	errInvalidUUID = errors.New("неверный формат uuid")
)

// ListPartsInputFromRequest собирает вход use case из rpc ListParts.
func ListPartsInputFromRequest(req *inventoryv1.ListPartsRequest) (input.ListPartsInput, error) {
	in := input.ListPartsInput{
		PartType: PartTypeFromProto(req.GetPartType()),
	}
	if len(req.GetUuids()) == 0 {
		return in, nil
	}

	ids := make([]uuid.UUID, 0, len(req.GetUuids()))
	for _, idStr := range req.GetUuids() {
		id, err := uuid.Parse(idStr)
		if err != nil {
			return input.ListPartsInput{}, errInvalidUUID
		}
		ids = append(ids, id)
	}
	in.IDs = ids

	return in, nil
}

// ParsePartUUID разбирает строковый UUID из rpc GetPart.
func ParsePartUUID(s string) (uuid.UUID, error) {
	if s == "" {
		return uuid.Nil, errEmptyUUID
	}
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, errInvalidUUID
	}

	return id, nil
}

// IsInvalidUUID возвращает true, если err — ошибка разбора UUID в конвертере.
func IsInvalidUUID(err error) bool {
	return errors.Is(err, errInvalidUUID) || errors.Is(err, errEmptyUUID)
}

// PartToProto переводит доменную модель в protobuf Part.
func PartToProto(p model.Part) *inventoryv1.Part {
	return &inventoryv1.Part{
		Uuid:          p.UUID,
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		PartType:      inventoryv1.PartType(p.PartType),
		StockQuantity: p.StockQuantity,
		CreatedAt:     timestamppb.New(p.CreatedAt),
	}
}

// PartsToProto переводит срез доменных моделей в protobuf.
func PartsToProto(parts []model.Part) []*inventoryv1.Part {
	out := make([]*inventoryv1.Part, 0, len(parts))
	for _, p := range parts {
		out = append(out, PartToProto(p))
	}

	return out
}

// PartTypeFromProto переводит protobuf enum в доменный тип.
func PartTypeFromProto(p inventoryv1.PartType) model.PartType {
	return model.PartType(p)
}
