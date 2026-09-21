package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewWitch creates a new Witch at the position passed.
func NewWitch(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(WitchType, witchConf)
}

var witchConf = MobConfig{
	MaxHealth:  26,
	Speed:      0.25,
	LavaDamage: 4,
	Breathable: Breathable{Air: true, TotalSupply: 15},
	Drops: LootTable{
		OneOf(Entry("minecraft:glowstone_dust", 0, 2, Looting(1)), Entry("minecraft:sugar", 0, 2, Looting(1)), Entry("minecraft:redstone", 0, 2, Looting(1)), Entry("minecraft:spider_eye", 0, 2, Looting(1)), Entry("minecraft:glass_bottle", 0, 2, Looting(1)), Entry("minecraft:gunpowder", 0, 2, Looting(1)), Entry("minecraft:stick", 0, 2, Weight(2), Looting(1))),
	},
	Goals: []Goal{
		Float(1),
		HurtByTarget(1),
		NearestAttackableTarget(2, 16, 10, false),
		RandomStroll(4, 1),
		LookAtPlayer(5, 8, 1, 2),
		RandomLookAround(5),
	},
}

// WitchType is a world.EntityType implementation for witchs.
var WitchType witchType

type witchType struct{}

func (witchType) EncodeEntity() string { return "minecraft:witch" }
func (witchType) SpawnEggName() string { return "witch" }
func (witchType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.3, 0, -0.3, 0.3, 1.9, 0.3)
}

func (witchType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (witchType) Apply(data *world.EntityData) { witchConf.Apply(data) }

func (witchType) DecodeNBT(m map[string]any, data *world.EntityData) {
	witchConf.Apply(data)
	decodeMob(m, data)
}

func (witchType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
