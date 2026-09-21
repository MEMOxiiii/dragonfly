package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewVindicator creates a new Vindicator at the position passed.
func NewVindicator(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(VindicatorType, vindicatorConf)
}

var vindicatorConf = MobConfig{
	MaxHealth:  24,
	Speed:      0.35,
	LavaDamage: 4,
	Breathable: Breathable{Air: true, TotalSupply: 15},
	Variants:   []VariantChoice{{0, 1}},
	Drops: LootTable{
		Always("minecraft:emerald", 0, 1, Looting(1)),
	},
	Goals: []Goal{
		Float(0),
		HurtByTarget(1),
		NearestAttackableTarget(2, 12, 12, false),
		MeleeAttack(3, 8, 1),
		RandomStroll(9, 1),
		LookAtPlayer(10, 8, 1, 2),
	},
}

// VindicatorType is a world.EntityType implementation for vindicators.
var VindicatorType vindicatorType

type vindicatorType struct{}

func (vindicatorType) EncodeEntity() string { return "minecraft:vindicator" }
func (vindicatorType) SpawnEggName() string { return "vindicator" }
func (vindicatorType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.3, 0, -0.3, 0.3, 1.9, 0.3)
}

func (vindicatorType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (vindicatorType) Apply(data *world.EntityData) { vindicatorConf.Apply(data) }

func (vindicatorType) DecodeNBT(m map[string]any, data *world.EntityData) {
	vindicatorConf.Apply(data)
	decodeMob(m, data)
}

func (vindicatorType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
