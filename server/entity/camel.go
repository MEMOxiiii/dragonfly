package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewCamel creates a new Camel at the position passed.
func NewCamel(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(CamelType, camelConf)
}

var camelConf = MobConfig{
	MaxHealth:       32,
	Speed:           0.09,
	ExperienceDrops: [2]int{1, 3},
	LavaDamage:      4,
	Breathable:      Breathable{Air: true, TotalSupply: 15},
	Healable:        Healable{Items: map[string]float64{"minecraft:cactus": 2}},
	BreedItems:      []string{"minecraft:cactus"},
	Goals: []Goal{
		Float(0),
		Panic(1, 4),
		Breed(2, 1),
		Tempt(3, 2.5, 10, "minecraft:cactus"),
		RandomStroll(6, 2),
		LookAtPlayer(7, 6, 2, 4),
		RandomLookAround(8),
	},
}

// CamelType is a world.EntityType implementation for camels.
var CamelType camelType

type camelType struct{}

func (camelType) EncodeEntity() string { return "minecraft:camel" }
func (camelType) SpawnEggName() string { return "camel" }
func (camelType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.85, 0, -0.85, 0.85, 2.375, 0.85)
}

func (camelType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (camelType) Apply(data *world.EntityData) { camelConf.Apply(data) }

func (camelType) DecodeNBT(m map[string]any, data *world.EntityData) {
	camelConf.Apply(data)
	decodeMob(m, data)
}

func (camelType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
