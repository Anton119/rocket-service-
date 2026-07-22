package assembly

import (
	"context"
	"math/rand/v2"
	"time"

	"github.com/google/uuid"

	"github.com/Anton119/rocket-service-/assembly/internal/model"
)

// Assemble эмулирует сборку корабля и возвращает событие ShipAssembled.
func (s *Service) Assemble(ctx context.Context, paid model.OrderPaidEvent) (model.ShipAssembledEvent, error) {
	minSec := s.cfg.MinSec()
	maxSec := s.cfg.MaxSec()
	buildTimeSec := int64(minSec + rand.IntN(maxSec-minSec+1)) //nolint:gosec // G404: учебная эмуляция сборки, crypto/rand не нужен

	if err := sleepWithContext(ctx, time.Duration(buildTimeSec)*time.Second); err != nil {
		return model.ShipAssembledEvent{}, err
	}

	return model.ShipAssembledEvent{
		EventUUID:    uuid.New().String(),
		OrderUUID:    paid.OrderUUID,
		UserUUID:     paid.UserUUID,
		BuildTimeSec: buildTimeSec,
		AssembledAt:  time.Now().UTC(),
	}, nil
}

func sleepWithContext(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
