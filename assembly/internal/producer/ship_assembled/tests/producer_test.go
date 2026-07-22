package ship_assembled_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/Anton119/rocket-service-/assembly/internal/model"
	shipassembled "github.com/Anton119/rocket-service-/assembly/internal/producer/ship_assembled"
	"github.com/Anton119/rocket-service-/platform/pkg/kafka"
	eventsv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/events/v1"
)

type fakeKafkaProducer struct {
	last *kafka.Message
	err  error
}

func (f *fakeKafkaProducer) Send(_ context.Context, msg *kafka.Message) error {
	f.last = msg
	return f.err
}

func TestProduce_MarshalsAndSends(t *testing.T) {
	t.Parallel()

	kafkaProd := &fakeKafkaProducer{}
	p := shipassembled.NewProducer(kafkaProd)

	err := p.Produce(context.Background(), model.ShipAssembledEvent{
		EventUUID:    "evt-1",
		OrderUUID:    "ord-1",
		UserUUID:     "usr-1",
		BuildTimeSec: 7,
		AssembledAt:  time.Unix(1700000000, 0).UTC(),
	})
	require.NoError(t, err)
	require.Equal(t, []byte("ord-1"), kafkaProd.last.Key)

	var pb eventsv1.ShipAssembled
	require.NoError(t, proto.Unmarshal(kafkaProd.last.Value, &pb))
	require.Equal(t, "evt-1", pb.GetEventUuid())
	require.Equal(t, "usr-1", pb.GetUserUuid())
	require.Equal(t, int64(7), pb.GetBuildTimeSec())
}

func TestProduce_PropagatesSendError(t *testing.T) {
	t.Parallel()

	want := errors.New("kafka down")
	p := shipassembled.NewProducer(&fakeKafkaProducer{err: want})
	err := p.Produce(context.Background(), model.ShipAssembledEvent{OrderUUID: "ord-1"})
	require.ErrorIs(t, err, want)
}
