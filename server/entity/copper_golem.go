package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewCopperGolem creates a new CopperGolem at the position passed.
func NewCopperGolem(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(CopperGolemType, copperGolemConf)
}

var copperGolemConf = MobConfig{
	MaxHealth:  12,
	Speed:      0.2,
	LavaDamage: 4,
	Goals: []Goal{
		Panic(2, 1.5),
		RandomStroll(5, 1),
		LookAtPlayer(6, 6, 1, 2),
		RandomLookAround(7),
	},
}

// CopperGolemType is a world.EntityType implementation for copper golems.
var CopperGolemType copperGolemType

type copperGolemType struct{}

func (copperGolemType) EncodeEntity() string { return "minecraft:copper_golem" }
func (copperGolemType) SpawnEggName() string { return "copper_golem" }
func (copperGolemType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.3, 0, -0.3, 0.3, 0.98, 0.3)
}

func (copperGolemType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (copperGolemType) Apply(data *world.EntityData) { copperGolemConf.Apply(data) }

func (copperGolemType) DecodeNBT(m map[string]any, data *world.EntityData) {
	copperGolemConf.Apply(data)
	decodeMob(m, data)
}

func (copperGolemType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
