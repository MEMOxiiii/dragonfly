package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewFox creates a new Fox at the position passed.
func NewFox(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(FoxType, foxConf)
}

var foxConf = MobConfig{
	MaxHealth:       20,
	Speed:           0.3,
	ExperienceDrops: [2]int{1, 3},
	Breathable:      Breathable{Air: true},
	Variants:        []VariantChoice{{0, 1}},
	BreedItems:      []string{"minecraft:sweet_berries"},
	Goals: []Goal{
		Float(0),
		Panic(1, 1.25),
		Tempt(3, 0.5, 16, "minecraft:sweet_berries"),
		Breed(3, 1),
		AvoidMobType(5, 10, 1, 1.5, "minecraft:player", "minecraft:polarbear", "minecraft:wolf"),
		MeleeAttack(10, 2, 1),
		RandomStroll(13, 0.8),
		LookAtPlayer(14, 6, 1, 2),
		RandomLookAround(15),
	},
}

// FoxType is a world.EntityType implementation for foxs.
var FoxType foxType

type foxType struct{}

func (foxType) EncodeEntity() string { return "minecraft:fox" }
func (foxType) SpawnEggName() string { return "fox" }
func (foxType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.3, 0, -0.3, 0.3, 0.7, 0.3)
}

func (foxType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (foxType) Apply(data *world.EntityData) { foxConf.Apply(data) }

func (foxType) DecodeNBT(m map[string]any, data *world.EntityData) {
	foxConf.Apply(data)
	decodeMob(m, data)
}

func (foxType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
