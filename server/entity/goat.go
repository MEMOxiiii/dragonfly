package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewGoat creates a new Goat at the position passed.
func NewGoat(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(GoatType, goatConf)
}

var goatConf = MobConfig{
	MaxHealth:       10,
	Speed:           0.4,
	ExperienceDrops: [2]int{1, 3},
	LavaDamage:      4,
	Breathable:      Breathable{Air: true, TotalSupply: 15},
	BreedItems:      []string{"minecraft:wheat"},
	Goals: []Goal{
		Float(0),
		Panic(1, 1),
		Breed(3, 0.6),
		Tempt(4, 0.75, 10, "minecraft:wheat"),
		NearestAttackableTarget(6, 16, 16, false),
		RandomStroll(9, 0.6),
		LookAtPlayer(10, 6, 2, 4),
		RandomLookAround(11),
	},
}

// GoatType is a world.EntityType implementation for goats.
var GoatType goatType

type goatType struct{}

func (goatType) EncodeEntity() string { return "minecraft:goat" }
func (goatType) SpawnEggName() string { return "goat" }
func (goatType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.45, 0, -0.45, 0.45, 1.3, 0.45)
}

func (goatType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (goatType) Apply(data *world.EntityData) { goatConf.Apply(data) }

func (goatType) DecodeNBT(m map[string]any, data *world.EntityData) {
	goatConf.Apply(data)
	decodeMob(m, data)
}

func (goatType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
