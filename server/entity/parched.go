package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewParched creates a new Parched at the position passed.
func NewParched(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(ParchedType, parchedConf)
}

var parchedConf = MobConfig{
	MaxHealth:       16,
	Speed:           0.25,
	ExperienceDrops: [2]int{5, 5},
	LavaDamage:      4,
	Breathable:      Breathable{Air: true, Water: true, TotalSupply: 15},
	Goals: []Goal{
		RangedAttack(1, 15, 1, 3, ShootArrow),
		HurtByTarget(2),
		NearestAttackableTarget(2, 16, 16, true),
		FleeSun(3, 1),
		AvoidMobType(5, 6, 1.2, 1.2, "minecraft:wolf"),
		RandomStroll(7, 1),
		LookAtPlayer(8, 8, 2, 4),
		RandomLookAround(9),
	},
}

// ParchedType is a world.EntityType implementation for parcheds.
var ParchedType parchedType

type parchedType struct{}

func (parchedType) EncodeEntity() string { return "minecraft:parched" }
func (parchedType) SpawnEggName() string { return "parched" }
func (parchedType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.3, 0, -0.3, 0.3, 1.9, 0.3)
}

func (parchedType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (parchedType) Apply(data *world.EntityData) { parchedConf.Apply(data) }

func (parchedType) DecodeNBT(m map[string]any, data *world.EntityData) {
	parchedConf.Apply(data)
	decodeMob(m, data)
}

func (parchedType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
