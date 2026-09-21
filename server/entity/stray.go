package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewStray creates a new Stray at the position passed.
func NewStray(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(StrayType, strayConf)
}

var strayConf = MobConfig{
	MaxHealth:       20,
	Speed:           0.25,
	ExperienceDrops: [2]int{5, 5},
	LavaDamage:      4,
	Breathable:      Breathable{Air: true, Water: true, TotalSupply: 15},
	BurnsInDaylight: true,
	Drops: LootTable{
		Always("minecraft:arrow", 0, 2, Looting(1)),
		Always("minecraft:bone", 0, 2, Looting(1)),
		Always("minecraft:arrow", 0, 1, Looting(1)),
	},
	Goals: []Goal{
		RangedAttack(0, 15, 1, 3, ShootArrow),
		HurtByTarget(1),
		FleeSun(2, 1),
		NearestAttackableTarget(2, 16, 16, true),
		AvoidMobType(4, 6, 1.2, 1.2, "minecraft:wolf"),
		RandomStroll(5, 1),
		LookAtPlayer(6, 8, 1, 2),
		RandomLookAround(6),
	},
}

// StrayType is a world.EntityType implementation for strays.
var StrayType strayType

type strayType struct{}

func (strayType) EncodeEntity() string { return "minecraft:stray" }
func (strayType) SpawnEggName() string { return "stray" }
func (strayType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.3, 0, -0.3, 0.3, 1.9, 0.3)
}

func (strayType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (strayType) Apply(data *world.EntityData) { strayConf.Apply(data) }

func (strayType) DecodeNBT(m map[string]any, data *world.EntityData) {
	strayConf.Apply(data)
	decodeMob(m, data)
}

func (strayType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
