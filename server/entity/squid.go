package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewSquid creates a new Squid at the position passed.
func NewSquid(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(SquidType, squidConf)
}

var squidConf = MobConfig{
	MaxHealth:       10,
	Speed:           0.2,
	ExperienceDrops: [2]int{1, 3},
	LavaDamage:      4,
	Breathable:      Breathable{Water: true, TotalSupply: 15},
	Drops: LootTable{
		Always("minecraft:ink_sac", 1, 3, Looting(1)),
	},
	Goals: []Goal{
		SquidSwim(1),
	},
}

// SquidType is a world.EntityType implementation for squids.
var SquidType squidType

type squidType struct{}

func (squidType) EncodeEntity() string { return "minecraft:squid" }
func (squidType) SpawnEggName() string { return "squid" }
func (squidType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.475, 0, -0.475, 0.475, 0.95, 0.475)
}

func (squidType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (squidType) Apply(data *world.EntityData) { squidConf.Apply(data) }

func (squidType) DecodeNBT(m map[string]any, data *world.EntityData) {
	squidConf.Apply(data)
	decodeMob(m, data)
}

func (squidType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
