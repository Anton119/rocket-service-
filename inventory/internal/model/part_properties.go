package model

// PartProperties — типоспецифичные свойства детали (pointer union).
type PartProperties struct {
	hull   *HullProperties
	engine *EngineProperties
	shield *ShieldProperties
	weapon *WeaponProperties
}

// Hull возвращает свойства корпуса, если они заданы.
func (p PartProperties) Hull() *HullProperties { return p.hull }

// Engine возвращает свойства двигателя, если они заданы.
func (p PartProperties) Engine() *EngineProperties { return p.engine }

// Shield возвращает свойства щита, если они заданы.
func (p PartProperties) Shield() *ShieldProperties { return p.shield }

// Weapon возвращает свойства оружия, если они заданы.
func (p PartProperties) Weapon() *WeaponProperties { return p.weapon }

// NewPartProperties собирает union из типоспецифичных свойств.
func NewPartProperties(
	hull *HullProperties,
	engine *EngineProperties,
	shield *ShieldProperties,
	weapon *WeaponProperties,
) PartProperties {
	return PartProperties{
		hull:   hull,
		engine: engine,
		shield: shield,
		weapon: weapon,
	}
}
