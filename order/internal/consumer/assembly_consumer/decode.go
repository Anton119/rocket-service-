package assembly_consumer

import (
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/Anton119/rocket-service-/order/internal/model"
	eventsv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/events/v1"
)

func decodeShipAssembled(data []byte) (model.ShipAssembledEvent, error) {
	var pb eventsv1.ShipAssembled
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.ShipAssembledEvent{}, fmt.Errorf("десериализовать ShipAssembled: %w", err)
	}

	event := model.ShipAssembledEvent{
		EventUUID:    pb.GetEventUuid(),
		OrderUUID:    pb.GetOrderUuid(),
		UserUUID:     pb.GetUserUuid(),
		BuildTimeSec: pb.GetBuildTimeSec(),
	}
	if pb.GetAssembledAt() != nil {
		event.AssembledAt = pb.GetAssembledAt().AsTime()
	} else {
		event.AssembledAt = time.Time{}
	}

	return event, nil
}
