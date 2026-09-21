package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewHappyGhast creates a new HappyGhast at the position passed.
func NewHappyGhast(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(HappyGhastType, happyGhastConf)
}

var happyGhastConf = MobConfig{
	MaxHealth:       20,
	Speed:           0.016,
	Travel:          TravelFloat,
	ExperienceDrops: [2]int{1, 3},
	LavaDamage:      4,
	Breathable:      Breathable{Air: true, TotalSupply: 5},
	Goals: []Goal{
		Float(0),
		RandomFly(7, 10, 7, 0, 0, 0, 1),
	},
}

// HappyGhastType is a world.EntityType implementation for happy ghasts.
var HappyGhastType happyGhastType

type happyGhastType struct{}

func (happyGhastType) EncodeEntity() string { return "minecraft:happy_ghast" }
func (happyGhastType) SpawnEggName() string { return "happy_ghast" }
func (happyGhastType) BBox(world.Entity) cube.BBox {
	return cube.Box(-2, 0, -2, 2, 4, 2)
}

func (happyGhastType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (happyGhastType) Apply(data *world.EntityData) { happyGhastConf.Apply(data) }

func (happyGhastType) DecodeNBT(m map[string]any, data *world.EntityData) {
	happyGhastConf.Apply(data)
	decodeMob(m, data)
}

func (happyGhastType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
