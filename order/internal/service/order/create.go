package order

import (
	"context"
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

	var totalSum int64
	for _, id := range uuids {
		totalSum += byUUID[id].Price
	}

	orderUUID := uuid.New()
	now := time.Now()

	o := model.Order{
		OrderUUID:  orderUUID,
		HullUUID:   hull,
		EngineUUID: engine,
		TotalPrice: totalSum,
		Status:     model.OrderStatusPendingPayment,
		CreatedAt:  now,
	}
	o.ShieldUUID = in.ShieldUUID
	o.WeaponUUID = in.WeaponUUID

	if err := s.repo.Save(ctx, o); err != nil {
		return nil, err
	}

	return &input.CreateOrderResult{
		OrderUUID:  orderUUID,
		TotalPrice: totalSum,
	}, nil
}
