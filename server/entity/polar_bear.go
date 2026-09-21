package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewPolarBear creates a new PolarBear at the position passed.
func NewPolarBear(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(PolarBearType, polarBearConf)
}

var polarBearConf = MobConfig{
	MaxHealth:       30,
	Speed:           0.25,
	ExperienceDrops: [2]int{1, 4},
	LavaDamage:      4,
	Breathable:      Breathable{Air: true, TotalSupply: 15},
	Drops: LootTable{
		Always("minecraft:cod", 0, 2, Looting(1)),
		Always("minecraft:salmon", 0, 2, Looting(1)),
	},
	Goals: []Goal{
		Float(0),
		HurtByTarget(1),
		Panic(2, 2),
		NearestAttackableTarget(4, 16, 16, false),
		RandomStroll(5, 1),
		LookAtPlayer(6, 8, 1, 2),
		RandomLookAround(7),
	},
}

// PolarBearType is a world.EntityType implementation for polar bears.
var PolarBearType polarBearType

type polarBearType struct{}

func (polarBearType) EncodeEntity() string { return "minecraft:polar_bear" }
func (polarBearType) SpawnEggName() string { return "polar_bear" }
func (polarBearType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.65, 0, -0.65, 0.65, 1.4, 0.65)
}

func (polarBearType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (polarBearType) Apply(data *world.EntityData) { polarBearConf.Apply(data) }

func (polarBearType) DecodeNBT(m map[string]any, data *world.EntityData) {
	polarBearConf.Apply(data)
	decodeMob(m, data)
}

func (polarBearType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
