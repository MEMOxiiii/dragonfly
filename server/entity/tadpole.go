package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewTadpole creates a new Tadpole at the position passed.
func NewTadpole(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(TadpoleType, tadpoleConf)
}

var tadpoleConf = MobConfig{
	MaxHealth:  6,
	Speed:      0.1,
	Travel:     TravelSwim,
	LavaDamage: 4,
	Breathable: Breathable{Water: true, TotalSupply: 8},
	Goals: []Goal{
		Panic(1, 2),
		RandomSwim(2, 16, 4, 100, 1),
		LookAtPlayer(3, 6, 2, 4),
		Tempt(5, 1.25, 10, "minecraft:slime_ball"),
	},
}

// TadpoleType is a world.EntityType implementation for tadpoles.
var TadpoleType tadpoleType

type tadpoleType struct{}

func (tadpoleType) EncodeEntity() string { return "minecraft:tadpole" }
func (tadpoleType) SpawnEggName() string { return "tadpole" }
func (tadpoleType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.4, 0, -0.4, 0.4, 0.6, 0.4)
}

func (tadpoleType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (tadpoleType) Apply(data *world.EntityData) { tadpoleConf.Apply(data) }

func (tadpoleType) DecodeNBT(m map[string]any, data *world.EntityData) {
	tadpoleConf.Apply(data)
	decodeMob(m, data)
}

func (tadpoleType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
