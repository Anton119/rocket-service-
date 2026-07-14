package model

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errs "github.com/Anton119/rocket-service-/inventory/internal/errors"
)

func TestPartReserveRelease(t *testing.T) {
	hullProps, err := NewHullProperties(50)
	require.NoError(t, err)

	part := RestorePart(
		uuid.New(),
		"Hull",
		"desc",
		PartTypeHull,
		100,
		2,
		0,
		hullProps,
		time.Now(),
	)

	require.NoError(t, part.Reserve())
	assert.Equal(t, 1, part.Reserved())
	assert.Equal(t, 1, part.Available())

	require.NoError(t, part.Reserve())
	assert.Equal(t, 0, part.Available())

	err = part.Reserve()
	require.ErrorIs(t, err, errs.ErrOutOfStock)

	require.NoError(t, part.Release())
	assert.Equal(t, 1, part.Reserved())

	err = part.Release()
	require.NoError(t, err)

	err = part.Release()
	require.ErrorIs(t, err, errs.ErrNothingToRelease)
}

func TestHullCanSupportEngine(t *testing.T) {
	hullProps, err := NewHullProperties(50)
	require.NoError(t, err)

	engineProps, err := NewEngineProperties(EngineClassC, 30)
	require.NoError(t, err)

	assert.True(t, hullProps.Hull().CanSupport(engineProps.Engine()))

	engineB, err := NewEngineProperties(EngineClassB, 70)
	require.NoError(t, err)
	assert.False(t, hullProps.Hull().CanSupport(engineB.Engine()))
}

func TestShieldConflictsWithWeapon(t *testing.T) {
	plasma, err := NewShieldProperties(ShieldTypePlasma)
	require.NoError(t, err)

	laser, err := NewWeaponProperties(WeaponTypeLaser)
	require.NoError(t, err)

	missile, err := NewWeaponProperties(WeaponTypeMissile)
	require.NoError(t, err)

	assert.True(t, plasma.Shield().ConflictsWith(laser.Weapon()))
	assert.False(t, plasma.Shield().ConflictsWith(missile.Weapon()))
}
