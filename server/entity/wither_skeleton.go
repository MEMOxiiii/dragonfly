package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewWitherSkeleton creates a new WitherSkeleton at the position passed.
func NewWitherSkeleton(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(WitherSkeletonType, witherSkeletonConf)
}

var witherSkeletonConf = MobConfig{
	MaxHealth:       20,
	Speed:           0.25,
	ExperienceDrops: [2]int{5, 5},
	FireImmune:      true,
	Breathable:      Breathable{Air: true, Water: true, TotalSupply: 15},
	Drops: LootTable{
		Always("minecraft:coal", 0, 1, Looting(1)),
		Always("minecraft:bone", 0, 2, Looting(1)),
		Always("minecraft:wither_skeleton_skull", 1, 1),
	},
	Goals: []Goal{
		HurtByTarget(1),
		NearestAttackableTarget(2, 16, 16, true),
		MeleeAttack(4, 4, 1.25),
		RandomStroll(5, 1),
		LookAtPlayer(6, 8, 1, 2),
		RandomLookAround(6),
	},
}

// WitherSkeletonType is a world.EntityType implementation for wither skeletons.
var WitherSkeletonType witherSkeletonType

type witherSkeletonType struct{}

func (witherSkeletonType) EncodeEntity() string { return "minecraft:wither_skeleton" }
func (witherSkeletonType) SpawnEggName() string { return "wither_skeleton" }
func (witherSkeletonType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.36, 0, -0.36, 0.36, 2.01, 0.36)
}

func (witherSkeletonType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (witherSkeletonType) Apply(data *world.EntityData) { witherSkeletonConf.Apply(data) }

func (witherSkeletonType) DecodeNBT(m map[string]any, data *world.EntityData) {
	witherSkeletonConf.Apply(data)
	decodeMob(m, data)
}

func (witherSkeletonType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
