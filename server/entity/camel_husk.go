package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewCamelHusk creates a new CamelHusk at the position passed.
func NewCamelHusk(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(CamelHuskType, camelHuskConf)
}

var camelHuskConf = MobConfig{
	MaxHealth:       32,
	Speed:           0.09,
	ExperienceDrops: [2]int{1, 3},
	LavaDamage:      4,
	Breathable:      Breathable{Air: true, Water: true, TotalSupply: 15},
	Healable:        Healable{Items: map[string]float64{"minecraft:rabbit_foot": 2}},
	Goals: []Goal{
		Float(0),
		Tempt(4, 2.5, 10, "minecraft:rabbit_foot"),
		RandomStroll(7, 2),
		LookAtPlayer(8, 6, 2, 4),
		RandomLookAround(9),
	},
}

// CamelHuskType is a world.EntityType implementation for camel husks.
var CamelHuskType camelHuskType

type camelHuskType struct{}

func (camelHuskType) EncodeEntity() string { return "minecraft:camel_husk" }
func (camelHuskType) SpawnEggName() string { return "camel_husk" }
func (camelHuskType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.85, 0, -0.85, 0.85, 2.375, 0.85)
}

func (camelHuskType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (camelHuskType) Apply(data *world.EntityData) { camelHuskConf.Apply(data) }

func (camelHuskType) DecodeNBT(m map[string]any, data *world.EntityData) {
	camelHuskConf.Apply(data)
	decodeMob(m, data)
}

func (camelHuskType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
