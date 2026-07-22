package assembly_consumer_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Anton119/rocket-service-/order/internal/consumer/assembly_consumer"
	"github.com/Anton119/rocket-service-/order/internal/model"
	"github.com/Anton119/rocket-service-/platform/pkg/kafka"
	eventsv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/events/v1"
)

type stubAssembler struct {
	last model.ShipAssembledEvent
	err  error
}

func (s *stubAssembler) HandleShipAssembled(_ context.Context, event model.ShipAssembledEvent) error {
	s.last = event
	return s.err
}

type stubConsumer struct{}

func (stubConsumer) Consume(context.Context, kafka.MessageHandler) error { return nil }

func TestShipAssembledHandler_HappyPath(t *testing.T) {
	t.Parallel()
	asm := &stubAssembler{}
	svc := assembly_consumer.NewService(stubConsumer{}, asm)

	payload, err := proto.Marshal(&eventsv1.ShipAssembled{
		EventUuid:    "e1",
		OrderUuid:    "o1",
		UserUuid:     "u1",
		BuildTimeSec: 8,
		AssembledAt:  timestamppb.New(time.Unix(1700000000, 0).UTC()),
	})
	require.NoError(t, err)

	require.NoError(t, svc.ShipAssembledHandler(context.Background(), kafka.Message{Value: payload}))
	require.Equal(t, "o1", asm.last.OrderUUID)
	require.Equal(t, "u1", asm.last.UserUUID)
}

func TestShipAssembledHandler_BadPayloadCommits(t *testing.T) {
	t.Parallel()
	svc := assembly_consumer.NewService(stubConsumer{}, &stubAssembler{})
	require.NoError(t, svc.ShipAssembledHandler(context.Background(), kafka.Message{Value: []byte("bad")}))
}

func TestShipAssembledHandler_BusinessError(t *testing.T) {
	t.Parallel()
	want := errors.New("commit failed")
	svc := assembly_consumer.NewService(stubConsumer{}, &stubAssembler{err: want})
	payload, err := proto.Marshal(&eventsv1.ShipAssembled{OrderUuid: "o1", UserUuid: "u1"})
	require.NoError(t, err)
	require.ErrorIs(t, svc.ShipAssembledHandler(context.Background(), kafka.Message{Value: payload}), want)
}
