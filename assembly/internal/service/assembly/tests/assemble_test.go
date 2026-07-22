package assembly_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/Anton119/rocket-service-/assembly/internal/model"
	assemblysvc "github.com/Anton119/rocket-service-/assembly/internal/service/assembly"
)

type fixedAssemblerCfg struct {
	min, max int
}

func (c fixedAssemblerCfg) MinSec() int { return c.min }
func (c fixedAssemblerCfg) MaxSec() int { return c.max }

func TestAssemble_RespectsContextCancel(t *testing.T) {
	t.Parallel()

	svc := assemblysvc.NewService(fixedAssemblerCfg{min: 5, max: 15})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := svc.Assemble(ctx, model.OrderPaidEvent{
		OrderUUID: "order-1",
		UserUUID:  "user-1",
	})
	require.ErrorIs(t, err, context.Canceled)
}

func TestAssemble_ReturnsEvent(t *testing.T) {
	t.Parallel()

	svc := assemblysvc.NewService(fixedAssemblerCfg{min: 0, max: 0})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	got, err := svc.Assemble(ctx, model.OrderPaidEvent{
		OrderUUID: "order-1",
		UserUUID:  "user-1",
	})
	require.NoError(t, err)
	require.Equal(t, "order-1", got.OrderUUID)
	require.Equal(t, "user-1", got.UserUUID)
	require.NotEmpty(t, got.EventUUID)
}
