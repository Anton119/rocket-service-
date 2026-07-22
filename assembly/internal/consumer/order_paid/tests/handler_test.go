package order_paid_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Anton119/rocket-service-/assembly/internal/consumer/order_paid"
	"github.com/Anton119/rocket-service-/assembly/internal/model"
	"github.com/Anton119/rocket-service-/platform/pkg/kafka"
	eventsv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/events/v1"
)

type stubAssembler struct {
	event model.ShipAssembledEvent
	err   error
}

func (s stubAssembler) Assemble(_ context.Context, paid model.OrderPaidEvent) (model.ShipAssembledEvent, error) {
	if s.err != nil {
		return model.ShipAssembledEvent{}, s.err
	}
	s.event.OrderUUID = paid.OrderUUID
	s.event.UserUUID = paid.UserUUID
	return s.event, nil
}

type stubProducer struct {
	last model.ShipAssembledEvent
	err  error
}

func (s *stubProducer) Produce(_ context.Context, event model.ShipAssembledEvent) error {
	s.last = event
	return s.err
}

type stubConsumer struct{}

func (stubConsumer) Consume(context.Context, kafka.MessageHandler) error { return nil }

func marshalOrderPaid(t *testing.T, orderUUID, userUUID string) []byte {
	t.Helper()
	payload, err := proto.Marshal(&eventsv1.OrderPaid{
		EventUuid:       "evt",
		OrderUuid:       orderUUID,
		TransactionUuid: "tx",
		PaymentMethod:   "CARD",
		UserUuid:        userUUID,
		PaidAt:          timestamppb.New(time.Unix(1700000000, 0).UTC()),
	})
	require.NoError(t, err)
	return payload
}

func TestOrderPaidHandler_HappyPath(t *testing.T) {
	t.Parallel()

	prod := &stubProducer{}
	svc := order_paid.NewService(
		stubConsumer{},
		stubAssembler{event: model.ShipAssembledEvent{EventUUID: "a1", BuildTimeSec: 5}},
		prod,
	)

	err := svc.OrderPaidHandler(context.Background(), kafka.Message{
		Value: marshalOrderPaid(t, "ord-1", "usr-1"),
	})
	require.NoError(t, err)
	require.Equal(t, "ord-1", prod.last.OrderUUID)
	require.Equal(t, "usr-1", prod.last.UserUUID)
}

func TestOrderPaidHandler_BadPayloadCommits(t *testing.T) {
	t.Parallel()

	svc := order_paid.NewService(stubConsumer{}, stubAssembler{}, &stubProducer{})
	err := svc.OrderPaidHandler(context.Background(), kafka.Message{Value: []byte("not-proto")})
	require.NoError(t, err)
}

func TestOrderPaidHandler_ProduceError(t *testing.T) {
	t.Parallel()

	want := errors.New("send failed")
	svc := order_paid.NewService(
		stubConsumer{},
		stubAssembler{event: model.ShipAssembledEvent{EventUUID: "a1"}},
		&stubProducer{err: want},
	)
	err := svc.OrderPaidHandler(context.Background(), kafka.Message{
		Value: marshalOrderPaid(t, "ord-1", "usr-1"),
	})
	require.ErrorIs(t, err, want)
}
