package order_paid

import (
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/Anton119/rocket-service-/assembly/internal/model"
	eventsv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/events/v1"
)

func decodeOrderPaid(data []byte) (model.OrderPaidEvent, error) {
	var pb eventsv1.OrderPaid
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.OrderPaidEvent{}, fmt.Errorf("десериализовать OrderPaid: %w", err)
	}

	event := model.OrderPaidEvent{
		EventUUID:       pb.GetEventUuid(),
		OrderUUID:       pb.GetOrderUuid(),
		TransactionUUID: pb.GetTransactionUuid(),
		PaymentMethod:   pb.GetPaymentMethod(),
		UserUUID:        pb.GetUserUuid(),
	}
	if pb.GetPaidAt() != nil {
		event.PaidAt = pb.GetPaidAt().AsTime()
	} else {
		event.PaidAt = time.Time{}
	}

	return event, nil
}
