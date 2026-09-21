package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewBlaze creates a new Blaze at the position passed.
func NewBlaze(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(BlazeType, blazeConf)
}

var blazeConf = MobConfig{
	MaxHealth:       20,
	Speed:           0.23,
	ExperienceDrops: [2]int{10, 10},
	FireImmune:      true,
	Drops: LootTable{
		Always("minecraft:blaze_rod", 0, 1, Looting(1)),
	},
	Goals: []Goal{
		Float(0),
		HurtByTarget(1),
		NearestAttackableTarget(2, 16, 48, false),
		MeleeAttack(3, 6, 1),
		RangedBurstAttack(4, 16, 3, 5, 3, 0.3, ShootSmallFireball),
		RandomStroll(5, 1),
		RandomLookAround(6),
	},
}

// BlazeType is a world.EntityType implementation for blazes.
var BlazeType blazeType

type blazeType struct{}

func (blazeType) EncodeEntity() string { return "minecraft:blaze" }
func (blazeType) SpawnEggName() string { return "blaze" }
func (blazeType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.25, 0, -0.25, 0.25, 1.8, 0.25)
}

func (blazeType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (blazeType) Apply(data *world.EntityData) { blazeConf.Apply(data) }

func (blazeType) DecodeNBT(m map[string]any, data *world.EntityData) {
	blazeConf.Apply(data)
	decodeMob(m, data)
}

func (blazeType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
