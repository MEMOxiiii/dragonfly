package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewPillager creates a new Pillager at the position passed.
func NewPillager(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(PillagerType, pillagerConf)
}

var pillagerConf = MobConfig{
	MaxHealth:  24,
	Speed:      0.35,
	LavaDamage: 4,
	Breathable: Breathable{Air: true, TotalSupply: 15},
	Variants:   []VariantChoice{{0, 1}},
	Drops: LootTable{
		Always("minecraft:arrow", 0, 2, Looting(1)),
	},
	Goals: []Goal{
		Float(0),
		HurtByTarget(1),
		NearestAttackableTarget(2, 16, 16, true),
		RangedAttack(4, 8, 1, 1, ShootArrow),
		RandomStroll(7, 1),
		LookAtPlayer(8, 8, 1, 2),
		RandomLookAround(8),
	},
}

// PillagerType is a world.EntityType implementation for pillagers.
var PillagerType pillagerType

type pillagerType struct{}

func (pillagerType) EncodeEntity() string { return "minecraft:pillager" }
func (pillagerType) SpawnEggName() string { return "pillager" }
func (pillagerType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.3, 0, -0.3, 0.3, 1.9, 0.3)
}

func (pillagerType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (pillagerType) Apply(data *world.EntityData) { pillagerConf.Apply(data) }

func (pillagerType) DecodeNBT(m map[string]any, data *world.EntityData) {
	pillagerConf.Apply(data)
	decodeMob(m, data)
}

func (pillagerType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
