package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewArmadillo creates a new Armadillo at the position passed.
func NewArmadillo(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(ArmadilloType, armadilloConf)
}

var armadilloConf = MobConfig{
	MaxHealth:       12,
	Speed:           0.14,
	ExperienceDrops: [2]int{1, 3},
	LavaDamage:      4,
	Breathable:      Breathable{Air: true, TotalSupply: 15},
	BreedItems:      []string{"minecraft:spider_eye"},
	Goals: []Goal{
		Float(0),
		Panic(1, 2),
		Breed(2, 1),
		Tempt(3, 1.25, 10, "minecraft:spider_eye"),
		RandomStroll(6, 1),
		LookAtPlayer(7, 6, 2, 4),
		RandomLookAround(8),
	},
}

// ArmadilloType is a world.EntityType implementation for armadillos.
var ArmadilloType armadilloType

type armadilloType struct{}

func (armadilloType) EncodeEntity() string { return "minecraft:armadillo" }
func (armadilloType) SpawnEggName() string { return "armadillo" }
func (armadilloType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.35, 0, -0.35, 0.35, 0.65, 0.35)
}

func (armadilloType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (armadilloType) Apply(data *world.EntityData) { armadilloConf.Apply(data) }

func (armadilloType) DecodeNBT(m map[string]any, data *world.EntityData) {
	armadilloConf.Apply(data)
	decodeMob(m, data)
}

func (armadilloType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
