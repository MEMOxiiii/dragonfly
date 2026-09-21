package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewGuardian creates a new Guardian at the position passed.
func NewGuardian(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(GuardianType, guardianConf)
}

var guardianConf = MobConfig{
	MaxHealth:       30,
	Speed:           0.12,
	Travel:          TravelSwim,
	ExperienceDrops: [2]int{10, 10},
	LavaDamage:      4,
	Breathable:      Breathable{Air: true, Water: true},
	Drops: LootTable{
		Always("minecraft:prismarine_shard", 0, 2, Looting(1)),
		OneOf(Entry("minecraft:cod", 1, 1, Weight(2), Looting(1)), Entry("minecraft:prismarine_crystals", 1, 1, Weight(2), Looting(1))),
	},
	Goals: []Goal{
		NearestAttackableTarget(1, 16, 16, false),
		RandomStroll(7, 1),
		LookAtPlayer(8, 12, 1, 2),
		RandomLookAround(9),
	},
}

// GuardianType is a world.EntityType implementation for guardians.
var GuardianType guardianType

type guardianType struct{}

func (guardianType) EncodeEntity() string { return "minecraft:guardian" }
func (guardianType) SpawnEggName() string { return "guardian" }
func (guardianType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.425, 0, -0.425, 0.425, 0.85, 0.425)
}

func (guardianType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (guardianType) Apply(data *world.EntityData) { guardianConf.Apply(data) }

func (guardianType) DecodeNBT(m map[string]any, data *world.EntityData) {
	guardianConf.Apply(data)
	decodeMob(m, data)
}

func (guardianType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
