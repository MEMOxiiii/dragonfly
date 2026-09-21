package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewBreeze creates a new Breeze at the position passed.
func NewBreeze(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(BreezeType, breezeConf)
}

var breezeConf = MobConfig{
	MaxHealth:       30,
	Speed:           0.4,
	ExperienceDrops: [2]int{10, 10},
	LavaDamage:      4,
	Breathable:      Breathable{Air: true, TotalSupply: 15},
	Goals: []Goal{
		Float(0),
		NearestAttackableTarget(1, 24, 24, false),
		HurtByTarget(4),
		RandomStroll(6, 1),
		LookAtPlayer(7, 16, 1, 2),
		RandomLookAround(8),
	},
}

// BreezeType is a world.EntityType implementation for breezes.
var BreezeType breezeType

type breezeType struct{}

func (breezeType) EncodeEntity() string { return "minecraft:breeze" }
func (breezeType) SpawnEggName() string { return "breeze" }
func (breezeType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.3, 0, -0.3, 0.3, 1.77, 0.3)
}

func (breezeType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (breezeType) Apply(data *world.EntityData) { breezeConf.Apply(data) }

func (breezeType) DecodeNBT(m map[string]any, data *world.EntityData) {
	breezeConf.Apply(data)
	decodeMob(m, data)
}

func (breezeType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
