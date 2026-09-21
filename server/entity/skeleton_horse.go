package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewSkeletonHorse creates a new SkeletonHorse at the position passed.
func NewSkeletonHorse(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(SkeletonHorseType, skeletonHorseConf)
}

var skeletonHorseConf = MobConfig{
	MaxHealth:       15,
	Speed:           0.2,
	ExperienceDrops: [2]int{1, 3},
	LavaDamage:      4,
	Breathable:      Breathable{Air: true, Water: true, TotalSupply: 15},
	Drops: LootTable{
		Always("minecraft:bone", 0, 2, Looting(1)),
	},
	Goals: []Goal{
		Panic(1, 1.2),
		RandomStroll(6, 0.7),
		LookAtPlayer(7, 6, 1, 2),
		RandomLookAround(8),
	},
}

// SkeletonHorseType is a world.EntityType implementation for skeleton horses.
var SkeletonHorseType skeletonHorseType

type skeletonHorseType struct{}

func (skeletonHorseType) EncodeEntity() string { return "minecraft:skeleton_horse" }
func (skeletonHorseType) SpawnEggName() string { return "skeleton_horse" }
func (skeletonHorseType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.7, 0, -0.7, 0.7, 1.6, 0.7)
}

func (skeletonHorseType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (skeletonHorseType) Apply(data *world.EntityData) { skeletonHorseConf.Apply(data) }

func (skeletonHorseType) DecodeNBT(m map[string]any, data *world.EntityData) {
	skeletonHorseConf.Apply(data)
	decodeMob(m, data)
}

func (skeletonHorseType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
