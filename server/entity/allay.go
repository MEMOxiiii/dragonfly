package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewAllay creates a new Allay at the position passed.
func NewAllay(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(AllayType, allayConf)
}

var allayConf = MobConfig{
	MaxHealth:  20,
	Speed:      0.1,
	Travel:     TravelHover,
	LavaDamage: 4,
	Breathable: Breathable{Air: true},
	Goals: []Goal{
		Panic(1, 2),
		Float(7),
		LookAtPlayer(8, 6, 2, 4),
		RandomLookAround(8),
		RandomFly(9, 8, 8, -1, 1, 4, 1),
	},
}

// AllayType is a world.EntityType implementation for allays.
var AllayType allayType

type allayType struct{}

func (allayType) EncodeEntity() string { return "minecraft:allay" }
func (allayType) SpawnEggName() string { return "allay" }
func (allayType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.175, 0, -0.175, 0.175, 0.6, 0.175)
}

func (allayType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (allayType) Apply(data *world.EntityData) { allayConf.Apply(data) }

func (allayType) DecodeNBT(m map[string]any, data *world.EntityData) {
	allayConf.Apply(data)
	decodeMob(m, data)
}

func (allayType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
