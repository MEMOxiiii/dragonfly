package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewHusk creates a new Husk at the position passed.
func NewHusk(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(HuskType, huskConf)
}

var huskConf = MobConfig{
	MaxHealth:       20,
	Speed:           0.23,
	ExperienceDrops: [2]int{5, 5},
	LavaDamage:      4,
	Breathable:      Breathable{Air: true, Water: true, TotalSupply: 15},
	Variants:        []VariantChoice{{2, 1}},
	Drops: LootTable{
		Always("minecraft:rotten_flesh", 0, 2, Looting(1)),
		OneOf(Entry("minecraft:iron_ingot", 1, 1), Entry("minecraft:carrot", 1, 1), Entry("minecraft:potato", 1, 1)),
	},
	Goals: []Goal{
		HurtByTarget(1),
		NearestAttackableTarget(2, 25, 35, true),
		MeleeAttack(3, 3, 1),
		RandomStroll(6, 1),
		LookAtPlayer(7, 6, 1, 2),
		RandomLookAround(7),
	},
}

// HuskType is a world.EntityType implementation for husks.
var HuskType huskType

type huskType struct{}

func (huskType) EncodeEntity() string { return "minecraft:husk" }
func (huskType) SpawnEggName() string { return "husk" }
func (huskType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.3, 0, -0.3, 0.3, 1.9, 0.3)
}

func (huskType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (huskType) Apply(data *world.EntityData) { huskConf.Apply(data) }

func (huskType) DecodeNBT(m map[string]any, data *world.EntityData) {
	huskConf.Apply(data)
	decodeMob(m, data)
}

func (huskType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
