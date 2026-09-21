package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewCreaking creates a new Creaking at the position passed.
func NewCreaking(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(CreakingType, creakingConf)
}

var creakingConf = MobConfig{
	MaxHealth:  1,
	Speed:      0.4,
	LavaDamage: 4,
	Goals: []Goal{
		Float(0),
		RandomStroll(7, 0.3),
	},
}

// CreakingType is a world.EntityType implementation for creakings.
var CreakingType creakingType

type creakingType struct{}

func (creakingType) EncodeEntity() string { return "minecraft:creaking" }
func (creakingType) SpawnEggName() string { return "creaking" }
func (creakingType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.45, 0, -0.45, 0.45, 2.7, 0.45)
}

func (creakingType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (creakingType) Apply(data *world.EntityData) { creakingConf.Apply(data) }

func (creakingType) DecodeNBT(m map[string]any, data *world.EntityData) {
	creakingConf.Apply(data)
	decodeMob(m, data)
}

func (creakingType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
