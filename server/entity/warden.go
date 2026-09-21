package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewWarden creates a new Warden at the position passed.
func NewWarden(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(WardenType, wardenConf)
}

var wardenConf = MobConfig{
	MaxHealth:           500,
	Speed:               0.3,
	ExperienceDrops:     [2]int{5, 5},
	FireImmune:          true,
	KnockBackResistance: 1,
	Breathable:          Breathable{Air: true, TotalSupply: 15},
	Goals: []Goal{
		Float(0),
		MeleeAttack(4, 30, 1.2),
		RandomStroll(9, 0.5),
		RandomLookAround(11),
	},
}

// WardenType is a world.EntityType implementation for wardens.
var WardenType wardenType

type wardenType struct{}

func (wardenType) EncodeEntity() string { return "minecraft:warden" }
func (wardenType) SpawnEggName() string { return "warden" }
func (wardenType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.45, 0, -0.45, 0.45, 2.9, 0.45)
}

func (wardenType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (wardenType) Apply(data *world.EntityData) { wardenConf.Apply(data) }

func (wardenType) DecodeNBT(m map[string]any, data *world.EntityData) {
	wardenConf.Apply(data)
	decodeMob(m, data)
}

func (wardenType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
