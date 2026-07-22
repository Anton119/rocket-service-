package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errs "github.com/Anton119/rocket-service-/inventory/internal/errors"
	"github.com/Anton119/rocket-service-/inventory/internal/model"
	"github.com/Anton119/rocket-service-/inventory/internal/service/domain"
)

func TestCompatibilityChecker_Check(t *testing.T) {
	checker := domain.NewCompatibilityChecker()

	hullProps, err := model.NewHullProperties(50)
	require.NoError(t, err)
	engineCProps, err := model.NewEngineProperties(model.EngineClassC)
	require.NoError(t, err)
	engineBProps, err := model.NewEngineProperties(model.EngineClassB)
	require.NoError(t, err)
	plasmaProps, err := model.NewShieldProperties(model.ShieldTypePlasma)
	require.NoError(t, err)
	laserProps, err := model.NewWeaponProperties(model.WeaponTypeLaser)
	require.NoError(t, err)

	tests := []struct {
		name    string
		slots   model.ResolvedShipSlots
		wantErr error
	}{
		{
			name: "корпус выдерживает двигатель класса C",
			slots: model.ResolvedShipSlots{
				Hull:   model.Part{PartType: model.PartTypeHull, Properties: hullProps},
				Engine: model.Part{PartType: model.PartTypeEngine, Properties: engineCProps},
			},
		},
		{
			name: "корпус не выдерживает двигатель класса B",
			slots: model.ResolvedShipSlots{
				Hull:   model.Part{PartType: model.PartTypeHull, Properties: hullProps},
				Engine: model.Part{PartType: model.PartTypeEngine, Properties: engineBProps},
			},
			wantErr: errs.ErrIncompatibleParts,
		},
		{
			name: "plasma щит конфликтует с laser оружием",
			slots: model.ResolvedShipSlots{
				Hull:   model.Part{PartType: model.PartTypeHull, Properties: hullProps},
				Engine: model.Part{PartType: model.PartTypeEngine, Properties: engineCProps},
				Shield: ptr(model.Part{PartType: model.PartTypeShield, Properties: plasmaProps}),
				Weapon: ptr(model.Part{PartType: model.PartTypeWeapon, Properties: laserProps}),
			},
			wantErr: errs.ErrIncompatibleParts,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := checker.Check(tc.slots)
			if tc.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

func ptr(p model.Part) *model.Part { return &p }
