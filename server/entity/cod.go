package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewCod creates a new Cod at the position passed.
func NewCod(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(CodType, codConf)
}

var codConf = MobConfig{
	MaxHealth:       6,
	Speed:           0.1,
	Travel:          TravelSwim,
	ExperienceDrops: [2]int{1, 3},
	LavaDamage:      4,
	Breathable:      Breathable{Water: true, TotalSupply: 15},
	Drops: LootTable{
		Always("minecraft:cod", 1, 1),
		Always("minecraft:bone", 1, 1, Looting(2)),
	},
	Goals: []Goal{
		AvoidMobType(1, 6, 1.5, 2, "minecraft:player"),
		RandomSwim(3, 16, 4, 0, 1),
	},
}

// CodType is a world.EntityType implementation for cods.
var CodType codType

type codType struct{}

func (codType) EncodeEntity() string { return "minecraft:cod" }
func (codType) SpawnEggName() string { return "cod" }
func (codType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.3, 0, -0.3, 0.3, 0.3, 0.3)
}

func (codType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (codType) Apply(data *world.EntityData) { codConf.Apply(data) }

func (codType) DecodeNBT(m map[string]any, data *world.EntityData) {
	codConf.Apply(data)
	decodeMob(m, data)
}

func (codType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
