package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewEnderman creates a new Enderman at the position passed.
func NewEnderman(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(EndermanType, endermanConf)
}

var endermanConf = MobConfig{
	MaxHealth:       40,
	Speed:           0.3,
	ExperienceDrops: [2]int{5, 5},
	LavaDamage:      4,
	Breathable:      Breathable{Air: true, TotalSupply: 15},
	Drops: LootTable{
		Always("minecraft:ender_pearl", 0, 1, Looting(1)),
	},
	Goals: []Goal{
		Float(0),
		HurtByTarget(3),
		NearestAttackableTarget(5, 16, 64, false),
		RandomStroll(7, 1),
		LookAtPlayer(8, 8, 2, 4),
		RandomLookAround(8),
	},
}

// EndermanType is a world.EntityType implementation for endermans.
var EndermanType endermanType

type endermanType struct{}

func (endermanType) EncodeEntity() string { return "minecraft:enderman" }
func (endermanType) SpawnEggName() string { return "enderman" }
func (endermanType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.3, 0, -0.3, 0.3, 2.9, 0.3)
}

func (endermanType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (endermanType) Apply(data *world.EntityData) { endermanConf.Apply(data) }

func (endermanType) DecodeNBT(m map[string]any, data *world.EntityData) {
	endermanConf.Apply(data)
	decodeMob(m, data)
}

func (endermanType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
