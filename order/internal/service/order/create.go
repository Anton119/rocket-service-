package order

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

	errs "github.com/Anton119/rocket-service-/order/internal/errors"
	"github.com/Anton119/rocket-service-/order/internal/model"
	"github.com/Anton119/rocket-service-/order/internal/service/input"
	"github.com/Anton119/rocket-service-/platform/pkg/auth"
)

// CreateOrder проверяет наличие деталей и создаёт заказ.
func (s *Service) CreateOrder(ctx context.Context, in input.CreateOrderInput) (result *input.CreateOrderResult, err error) {
	ctx, span := otel.Tracer(tracerName).Start(ctx, "order.Create")
	defer func() {
		if err != nil {
			recordSpanError(span, err)
		}
		span.End()
	}()

	span.SetAttributes(
		attribute.Stringer("order.hull_uuid", in.HullUUID),
		attribute.Stringer("order.engine_uuid", in.EngineUUID),
	)

	userUUIDStr, ok := auth.UserUUIDFromContext(ctx)
	if !ok || userUUIDStr == "" {
		return nil, errs.ErrUnauthorized
	}

	userUUID, err := uuid.Parse(userUUIDStr)
	if err != nil {
		return nil, errs.ErrUnauthorized
	}

	uuids := partUUIDsFromInput(in)

	parts, err := s.inv.ListParts(ctx, uuids)
	if err != nil {
		return nil, err
	}

	hull := in.HullUUID
	engine := in.EngineUUID

	byUUID := indexPartsByUUID(parts)

	if err := validatePartsStock(byUUID, uuids); err != nil {
		return nil, err
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
		UserUUID:   userUUID,
		Items:      items,
		HullUUID:   hull,
		EngineUUID: engine,
		TotalPrice: totalSum,
		Status:     model.OrderStatusPendingPayment,
		CreatedAt:  now,
	}
	o.ShieldUUID = in.ShieldUUID
	o.WeaponUUID = in.WeaponUUID

	if err := s.persistOrder(ctx, o, uuids, userUUID); err != nil {
		return nil, err
	}

	recordOrderCreated(ctx, orderUUID, userUUID, totalSum)

	span.SetAttributes(
		attribute.Stringer("order.uuid", orderUUID),
		attribute.Int("order.items_count", len(items)),
		attribute.Int64("order.total_price", totalSum),
	)
	span.SetStatus(codes.Ok, "")

	return &input.CreateOrderResult{
		OrderUUID:  orderUUID,
		TotalPrice: totalSum,
	}, nil
}
