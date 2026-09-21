package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewEndermite creates a new Endermite at the position passed.
func NewEndermite(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(EndermiteType, endermiteConf)
}

var endermiteConf = MobConfig{
	MaxHealth:       8,
	Speed:           0.25,
	ExperienceDrops: [2]int{3, 3},
	LavaDamage:      4,
	Breathable:      Breathable{Air: true, TotalSupply: 15},
	Goals: []Goal{
		Float(0),
		MeleeAttack(3, 2, 1),
		NearestAttackableTarget(5, 16, 16, false),
		RandomStroll(6, 1),
	},
}

// EndermiteType is a world.EntityType implementation for endermites.
var EndermiteType endermiteType

type endermiteType struct{}

func (endermiteType) EncodeEntity() string { return "minecraft:endermite" }
func (endermiteType) SpawnEggName() string { return "endermite" }
func (endermiteType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.2, 0, -0.2, 0.2, 0.3, 0.2)
}

func (endermiteType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (endermiteType) Apply(data *world.EntityData) { endermiteConf.Apply(data) }

func (endermiteType) DecodeNBT(m map[string]any, data *world.EntityData) {
	endermiteConf.Apply(data)
	decodeMob(m, data)
}

func (endermiteType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
