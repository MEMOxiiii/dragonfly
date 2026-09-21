package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewRavager creates a new Ravager at the position passed.
func NewRavager(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(RavagerType, ravagerConf)
}

var ravagerConf = MobConfig{
	MaxHealth:           100,
	Speed:               0.3,
	ExperienceDrops:     [2]int{20, 20},
	KnockBackResistance: 0.5,
	LavaDamage:          4,
	Breathable:          Breathable{Air: true, TotalSupply: 15},
	Drops: LootTable{
		Always("minecraft:saddle", 1, 1),
	},
	Goals: []Goal{
		Float(0),
		HurtByTarget(2),
		NearestAttackableTarget(3, 16, 16, false),
		RandomStroll(6, 0.4),
		LookAtPlayer(7, 6, 2, 4),
	},
}

// RavagerType is a world.EntityType implementation for ravagers.
var RavagerType ravagerType

type ravagerType struct{}

func (ravagerType) EncodeEntity() string { return "minecraft:ravager" }
func (ravagerType) SpawnEggName() string { return "ravager" }
func (ravagerType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.6, 0, -0.6, 0.6, 1.9, 0.6)
}

func (ravagerType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (ravagerType) Apply(data *world.EntityData) { ravagerConf.Apply(data) }

func (ravagerType) DecodeNBT(m map[string]any, data *world.EntityData) {
	ravagerConf.Apply(data)
	decodeMob(m, data)
}

func (ravagerType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
