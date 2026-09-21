package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewPiglin creates a new Piglin at the position passed.
func NewPiglin(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(PiglinType, piglinConf)
}

var piglinConf = MobConfig{
	MaxHealth:       16,
	Speed:           0.35,
	ExperienceDrops: [2]int{5, 5},
	LavaDamage:      4,
	Breathable:      Breathable{Air: true, TotalSupply: 15},
	Goals: []Goal{
		HurtByTarget(1),
		AvoidMobType(4, 6, 1, 1.2, "minecraft:hoglin", "minecraft:zoglin", "minecraft:zombie_pigman"),
		RandomStroll(10, 0.6),
		LookAtPlayer(11, 8, 2, 4),
		RandomLookAround(12),
	},
}

// PiglinType is a world.EntityType implementation for piglins.
var PiglinType piglinType

type piglinType struct{}

func (piglinType) EncodeEntity() string { return "minecraft:piglin" }
func (piglinType) SpawnEggName() string { return "piglin" }
func (piglinType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.3, 0, -0.3, 0.3, 1.9, 0.3)
}

func (piglinType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (piglinType) Apply(data *world.EntityData) { piglinConf.Apply(data) }

func (piglinType) DecodeNBT(m map[string]any, data *world.EntityData) {
	piglinConf.Apply(data)
	decodeMob(m, data)
}

func (piglinType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
