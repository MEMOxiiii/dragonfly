package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewShulker creates a new Shulker at the position passed.
func NewShulker(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(ShulkerType, shulkerConf)
}

var shulkerConf = MobConfig{
	MaxHealth:       30,
	Speed:           0,
	ExperienceDrops: [2]int{5, 5},
	LavaDamage:      4,
	Breathable:      Breathable{Air: true, TotalSupply: 15},
	Variants:        []VariantChoice{{16, 1}},
	Drops: LootTable{
		Always("minecraft:shulker_shell", 0, 1, Looting(1)),
	},
	Goals: []Goal{
		LookAtPlayer(1, 6, 1, 2),
		HurtByTarget(2),
		NearestAttackableTarget(3, 16, 16, false),
		RandomLookAround(8),
	},
}

// ShulkerType is a world.EntityType implementation for shulkers.
var ShulkerType shulkerType

type shulkerType struct{}

func (shulkerType) EncodeEntity() string { return "minecraft:shulker" }
func (shulkerType) SpawnEggName() string { return "shulker" }
func (shulkerType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.3, 0, -0.3, 0.3, 1.8, 0.3)
}

func (shulkerType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (shulkerType) Apply(data *world.EntityData) { shulkerConf.Apply(data) }

func (shulkerType) DecodeNBT(m map[string]any, data *world.EntityData) {
	shulkerConf.Apply(data)
	decodeMob(m, data)
}

func (shulkerType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
