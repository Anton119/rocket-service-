package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	errs "github.com/Anton119/rocket-service-/order/internal/errors"
	"github.com/Anton119/rocket-service-/order/internal/model"
)

// CancelOrder отменяет заказ в статусе ожидания оплаты (с блокировкой строки)
// и снимает резерв деталей в Inventory.
func (s *Service) CancelOrder(ctx context.Context, id uuid.UUID) error {
	var partUUIDs []string

	err := s.txManager.Do(ctx, func(txCtx context.Context) error {
		o, err := s.repo.GetForUpdate(txCtx, id)
		if err != nil {
			return fmt.Errorf("получить заказ: %w", err)
		}

		switch o.Status {
		case model.OrderStatusPendingPayment:
			// ok
		case model.OrderStatusPaid, model.OrderStatusAssembled:
			return errs.ErrOrderAlreadyPaid
		case model.OrderStatusCancelled:
			return errs.ErrOrderAlreadyCancelled
		default:
			return errs.ErrOrderCancelNotAllowed
		}

		partUUIDs = orderPartUUIDs(o)

		o.Status = model.OrderStatusCancelled
		if err := s.repo.Save(txCtx, o); err != nil {
			return fmt.Errorf("сохранить заказ: %w", err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	if len(partUUIDs) == 0 {
		return nil
	}

	if err := s.inv.ReleaseParts(ctx, partUUIDs); err != nil {
		return fmt.Errorf("освободить резерв деталей: %w", err)
	}

	return nil
}

func orderPartUUIDs(o model.Order) []string {
	if len(o.Items) > 0 {
		uuids := make([]string, 0, len(o.Items))
		seen := make(map[uuid.UUID]struct{}, len(o.Items))
		for _, item := range o.Items {
			if _, ok := seen[item.PartUUID]; ok {
				continue
			}
			seen[item.PartUUID] = struct{}{}
			uuids = append(uuids, item.PartUUID.String())
		}
		return uuids
	}

	uuids := []string{o.HullUUID.String(), o.EngineUUID.String()}
	if o.ShieldUUID != nil {
		uuids = append(uuids, o.ShieldUUID.String())
	}
	if o.WeaponUUID != nil {
		uuids = append(uuids, o.WeaponUUID.String())
	}
	return uuids
}
