package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewBogged creates a new Bogged at the position passed.
func NewBogged(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(BoggedType, boggedConf)
}

var boggedConf = MobConfig{
	MaxHealth:       16,
	Speed:           0.25,
	ExperienceDrops: [2]int{5, 5},
	LavaDamage:      4,
	Breathable:      Breathable{Air: true, Water: true, TotalSupply: 15},
	BurnsInDaylight: true,
	Goals: []Goal{
		RangedAttack(0, 15, 1, 3, ShootArrow),
		HurtByTarget(1),
		FleeSun(2, 1),
		NearestAttackableTarget(2, 16, 16, true),
		AvoidMobType(4, 6, 1.2, 1.2, "minecraft:wolf"),
		RandomStroll(6, 1),
		LookAtPlayer(7, 8, 2, 4),
		RandomLookAround(8),
	},
}

// BoggedType is a world.EntityType implementation for boggeds.
var BoggedType boggedType

type boggedType struct{}

func (boggedType) EncodeEntity() string { return "minecraft:bogged" }
func (boggedType) SpawnEggName() string { return "bogged" }
func (boggedType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.3, 0, -0.3, 0.3, 1.9, 0.3)
}

func (boggedType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (boggedType) Apply(data *world.EntityData) { boggedConf.Apply(data) }

func (boggedType) DecodeNBT(m map[string]any, data *world.EntityData) {
	boggedConf.Apply(data)
	decodeMob(m, data)
}

func (boggedType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
