package order

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	errs "github.com/Anton119/rocket-service-/order/internal/errors"
	"github.com/Anton119/rocket-service-/order/internal/model"
	"github.com/Anton119/rocket-service-/order/internal/service/input"
)

// CreateOrder проверяет наличие деталей и создаёт заказ.
func (s *Service) CreateOrder(ctx context.Context, in input.CreateOrderInput) (*input.CreateOrderResult, error) {
	hull := in.HullUUID
	engine := in.EngineUUID

	uuids := []string{hull.String(), engine.String()}
	if in.ShieldUUID != nil {
		uuids = append(uuids, in.ShieldUUID.String())
	}
	if in.WeaponUUID != nil {
		uuids = append(uuids, in.WeaponUUID.String())
	}

	parts, err := s.inv.ListParts(ctx, uuids)
	if err != nil {
		return nil, err
	}

	byUUID := make(map[string]model.Part, len(parts))
	for i := range parts {
		p := parts[i]
		byUUID[p.UUID] = p
	}

	for _, id := range uuids {
		p, ok := byUUID[id]
		if !ok {
			return nil, errs.ErrPartNotFound
		}
		if p.StockQuantity <= 0 {
			return nil, errs.ErrPartOutOfStock
		}
	}

	if err := s.inv.ReserveParts(ctx, uuids); err != nil {
		return nil, err
	}

	items := mergeOrderItems([]orderSlot{
		{partUUID: hull, partType: "HULL", price: byUUID[hull.String()].Price},
		{partUUID: engine, partType: "ENGINE", price: byUUID[engine.String()].Price},
	})

	var totalSum int64
	for _, item := range items {
		totalSum += item.Price
	}

	if in.ShieldUUID != nil {
		shield := *in.ShieldUUID
		items = appendMergedItem(items, shield, "SHIELD", byUUID[shield.String()].Price)
		totalSum += byUUID[shield.String()].Price
	}
	if in.WeaponUUID != nil {
		weapon := *in.WeaponUUID
		items = appendMergedItem(items, weapon, "WEAPON", byUUID[weapon.String()].Price)
		totalSum += byUUID[weapon.String()].Price
	}

	orderUUID := uuid.New()
	now := time.Now()

	o := model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   in.UserUUID,
		Items:      items,
		HullUUID:   hull,
		EngineUUID: engine,
		TotalPrice: totalSum,
		Status:     model.OrderStatusPendingPayment,
		CreatedAt:  now,
	}
	o.ShieldUUID = in.ShieldUUID
	o.WeaponUUID = in.WeaponUUID

	if err := s.repo.Create(ctx, o); err != nil {
		if releaseErr := s.inv.ReleaseParts(ctx, uuids); releaseErr != nil {
			return nil, errors.Join(fmt.Errorf("сохранение заказа: %w", err), fmt.Errorf("освободить резерв: %w", releaseErr))
		}
		return nil, fmt.Errorf("сохранение заказа: %w", err)
	}

	return &input.CreateOrderResult{
		OrderUUID:  orderUUID,
		TotalPrice: totalSum,
	}, nil
}
