package v1

import (
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func parseUUIDs(raw []string) ([]uuid.UUID, error) {
	uuids := make([]uuid.UUID, 0, len(raw))
	for _, s := range raw {
		id, err := uuid.Parse(s)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "невалидный uuid: %s", s)
		}
		uuids = append(uuids, id)
	}
	return uuids, nil
}
