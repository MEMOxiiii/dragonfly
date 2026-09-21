package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewElderGuardian creates a new ElderGuardian at the position passed.
func NewElderGuardian(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(ElderGuardianType, elderGuardianConf)
}

var elderGuardianConf = MobConfig{
	MaxHealth:       80,
	Speed:           0.3,
	Travel:          TravelSwim,
	ExperienceDrops: [2]int{10, 10},
	LavaDamage:      4,
	Breathable:      Breathable{Air: true, Water: true},
	Drops: LootTable{
		Always("minecraft:prismarine_shard", 0, 2, Looting(1)),
		OneOf(Entry("minecraft:cod", 1, 1, Weight(3), Looting(1)), Entry("minecraft:prismarine_crystals", 1, 1, Weight(2), Looting(1))),
		Always("minecraft:sponge", 1, 1),
	},
	Goals: []Goal{
		NearestAttackableTarget(1, 16, 16, false),
		RandomStroll(7, 0.5),
		LookAtPlayer(8, 12, 1, 2),
		RandomLookAround(9),
	},
}

// ElderGuardianType is a world.EntityType implementation for elder guardians.
var ElderGuardianType elderGuardianType

type elderGuardianType struct{}

func (elderGuardianType) EncodeEntity() string { return "minecraft:elder_guardian" }
func (elderGuardianType) SpawnEggName() string { return "elder_guardian" }
func (elderGuardianType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.995, 0, -0.995, 0.995, 1.99, 0.995)
}

func (elderGuardianType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (elderGuardianType) Apply(data *world.EntityData) { elderGuardianConf.Apply(data) }

func (elderGuardianType) DecodeNBT(m map[string]any, data *world.EntityData) {
	elderGuardianConf.Apply(data)
	decodeMob(m, data)
}

func (elderGuardianType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
