package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewBat creates a new Bat at the position passed.
func NewBat(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(BatType, batConf)
}

var batConf = MobConfig{
	MaxHealth:  6,
	Speed:      0.1,
	Travel:     TravelFloat,
	LavaDamage: 4,
	Breathable: Breathable{Air: true, TotalSupply: 15},
	Goals: []Goal{
		Float(0),
		RandomFly(0, 10, 7, -2, 0, 0, 1),
	},
}

// BatType is a world.EntityType implementation for bats.
var BatType batType

type batType struct{}

func (batType) EncodeEntity() string { return "minecraft:bat" }
func (batType) SpawnEggName() string { return "bat" }
func (batType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.25, 0, -0.25, 0.25, 0.9, 0.25)
}

func (batType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (batType) Apply(data *world.EntityData) { batConf.Apply(data) }

func (batType) DecodeNBT(m map[string]any, data *world.EntityData) {
	batConf.Apply(data)
	decodeMob(m, data)
}

func (batType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
