package order_producer_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/Anton119/rocket-service-/order/internal/model"
	orderproducer "github.com/Anton119/rocket-service-/order/internal/producer/order_producer"
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

	kp := &fakeKafkaProducer{}
	p := orderproducer.NewProducer(kp)
	err := p.Produce(context.Background(), model.OrderPaidEvent{
		EventUUID:       "evt-1",
		OrderUUID:       "ord-1",
		TransactionUUID: "tx-1",
		PaymentMethod:   "CARD",
		UserUUID:        "usr-1",
		PaidAt:          time.Unix(1700000000, 0).UTC(),
	})
	require.NoError(t, err)
	require.Equal(t, []byte("ord-1"), kp.last.Key)

	var pb eventsv1.OrderPaid
	require.NoError(t, proto.Unmarshal(kp.last.Value, &pb))
	require.Equal(t, "usr-1", pb.GetUserUuid())
	require.Equal(t, "tx-1", pb.GetTransactionUuid())
}

func TestProduce_PropagatesError(t *testing.T) {
	t.Parallel()
	want := errors.New("kafka down")
	p := orderproducer.NewProducer(&fakeKafkaProducer{err: want})
	err := p.Produce(context.Background(), model.OrderPaidEvent{OrderUUID: "ord-1"})
	require.ErrorIs(t, err, want)
}
