package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	errs "github.com/Anton119/rocket-service-/inventory/internal/errors"
	"github.com/Anton119/rocket-service-/inventory/internal/model"
)

func TestCompatibilityChecker(t *testing.T) {
	checker := NewCompatibilityChecker()

	hullProps, err := model.NewHullProperties(50)
	require.NoError(t, err)
	engineC, err := model.NewEngineProperties(model.EngineClassC, 30)
	require.NoError(t, err)
	engineB, err := model.NewEngineProperties(model.EngineClassB, 70)
	require.NoError(t, err)
	energyShield, err := model.NewShieldProperties(model.ShieldTypeEnergy)
	require.NoError(t, err)
	plasmaShield, err := model.NewShieldProperties(model.ShieldTypePlasma)
	require.NoError(t, err)
	laser, err := model.NewWeaponProperties(model.WeaponTypeLaser)
	require.NoError(t, err)

	hull := model.RestorePart(uuid.New(), "hull", "", model.PartTypeHull, 1, 1, 0, hullProps, time.Now())
	engineCPart := model.RestorePart(uuid.New(), "engine C", "", model.PartTypeEngine, 1, 1, 0, engineC, time.Now())
	engineBPart := model.RestorePart(uuid.New(), "engine B", "", model.PartTypeEngine, 1, 1, 0, engineB, time.Now())
	shieldEnergy := model.RestorePart(uuid.New(), "shield", "", model.PartTypeShield, 1, 1, 0, energyShield, time.Now())
	shieldPlasma := model.RestorePart(uuid.New(), "plasma", "", model.PartTypeShield, 1, 1, 0, plasmaShield, time.Now())
	weaponLaser := model.RestorePart(uuid.New(), "laser", "", model.PartTypeWeapon, 1, 1, 0, laser, time.Now())

	require.NoError(t, checker.Check(model.ResolvedShipSlots{Hull: hull, Engine: engineCPart}))

	err = checker.Check(model.ResolvedShipSlots{Hull: hull, Engine: engineBPart})
	require.ErrorIs(t, err, errs.ErrIncompatibleParts)

	require.NoError(t, checker.Check(model.ResolvedShipSlots{
		Hull: hull, Engine: engineCPart, Shield: &shieldEnergy, Weapon: &weaponLaser,
	}))

	err = checker.Check(model.ResolvedShipSlots{
		Hull: hull, Engine: engineCPart, Shield: &shieldPlasma, Weapon: &weaponLaser,
	})
	require.ErrorIs(t, err, errs.ErrIncompatibleParts)
}
