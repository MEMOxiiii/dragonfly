package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewSalmon creates a new Salmon at the position passed.
func NewSalmon(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(SalmonType, salmonConf)
}

var salmonConf = MobConfig{
	MaxHealth:       6,
	Speed:           0.12,
	Travel:          TravelSwim,
	ExperienceDrops: [2]int{1, 3},
	LavaDamage:      4,
	Breathable:      Breathable{Water: true, TotalSupply: 15},
	Drops: LootTable{
		Always("minecraft:salmon", 1, 1),
		Always("minecraft:bone", 1, 1, Looting(2)),
	},
	Goals: []Goal{
		AvoidMobType(1, 3, 1.5, 2, "minecraft:player"),
		RandomSwim(3, 16, 4, 0, 1),
	},
}

// SalmonType is a world.EntityType implementation for salmons.
var SalmonType salmonType

type salmonType struct{}

func (salmonType) EncodeEntity() string { return "minecraft:salmon" }
func (salmonType) SpawnEggName() string { return "salmon" }
func (salmonType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.25, 0, -0.25, 0.25, 0.5, 0.25)
}

func (salmonType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (salmonType) Apply(data *world.EntityData) { salmonConf.Apply(data) }

func (salmonType) DecodeNBT(m map[string]any, data *world.EntityData) {
	salmonConf.Apply(data)
	decodeMob(m, data)
}

func (salmonType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
