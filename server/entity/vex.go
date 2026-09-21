package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewVex creates a new Vex at the position passed.
func NewVex(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(VexType, vexConf)
}

var vexConf = MobConfig{
	MaxHealth:       14,
	Speed:           1,
	ExperienceDrops: [2]int{5, 5},
	FireImmune:      true,
	Goals: []Goal{
		Float(0),
		HurtByTarget(1),
		NearestAttackableTarget(3, 16, 70, false),
		LookAtPlayer(9, 6, 1, 2),
	},
}

// VexType is a world.EntityType implementation for vexs.
var VexType vexType

type vexType struct{}

func (vexType) EncodeEntity() string { return "minecraft:vex" }
func (vexType) SpawnEggName() string { return "vex" }
func (vexType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.2, 0, -0.2, 0.2, 0.8, 0.2)
}

func (vexType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (vexType) Apply(data *world.EntityData) { vexConf.Apply(data) }

func (vexType) DecodeNBT(m map[string]any, data *world.EntityData) {
	vexConf.Apply(data)
	decodeMob(m, data)
}

func (vexType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
