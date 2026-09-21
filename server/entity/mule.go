package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewMule creates a new Mule at the position passed.
func NewMule(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(MuleType, muleConf)
}

var muleConf = MobConfig{
	MaxHealth:       22.5,
	Speed:           0.175,
	ExperienceDrops: [2]int{1, 3},
	LavaDamage:      4,
	Breathable:      Breathable{Air: true, TotalSupply: 15},
	Healable:        Healable{Items: map[string]float64{"minecraft:wheat": 2, "minecraft:sugar": 1, "minecraft:hay_block": 20, "minecraft:apple": 3, "minecraft:golden_carrot": 4, "minecraft:golden_apple": 10, "minecraft:enchanted_golden_apple": 10}},
	Drops: LootTable{
		Always("minecraft:leather", 0, 2, Looting(1)),
	},
	Goals: []Goal{
		Float(0),
		Panic(1, 1.2),
		RandomStroll(6, 0.7),
		LookAtPlayer(7, 6, 1, 2),
		RandomLookAround(8),
	},
}

// MuleType is a world.EntityType implementation for mules.
var MuleType muleType

type muleType struct{}

func (muleType) EncodeEntity() string { return "minecraft:mule" }
func (muleType) SpawnEggName() string { return "mule" }
func (muleType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.7, 0, -0.7, 0.7, 1.6, 0.7)
}

func (muleType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (muleType) Apply(data *world.EntityData) { muleConf.Apply(data) }

func (muleType) DecodeNBT(m map[string]any, data *world.EntityData) {
	muleConf.Apply(data)
	decodeMob(m, data)
}

func (muleType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
