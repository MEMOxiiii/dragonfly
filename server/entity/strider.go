package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewStrider creates a new Strider at the position passed.
func NewStrider(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(StriderType, striderConf)
}

var striderConf = MobConfig{
	MaxHealth:       20,
	Speed:           0.16,
	ExperienceDrops: [2]int{1, 3},
	FireImmune:      true,
	BreedItems:      []string{"minecraft:warped_fungus"},
	Goals: []Goal{
		Panic(3, 1.1),
		Breed(4, 1),
		Tempt(5, 1.2, 10, "minecraft:warped_fungus", "minecraft:warped_fungus_on_a_stick"),
		LookAtPlayer(9, 6, 2, 4),
		RandomLookAround(10),
	},
}

// StriderType is a world.EntityType implementation for striders.
var StriderType striderType

type striderType struct{}

func (striderType) EncodeEntity() string { return "minecraft:strider" }
func (striderType) SpawnEggName() string { return "strider" }
func (striderType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.45, 0, -0.45, 0.45, 1.7, 0.45)
}

func (striderType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (striderType) Apply(data *world.EntityData) { striderConf.Apply(data) }

func (striderType) DecodeNBT(m map[string]any, data *world.EntityData) {
	striderConf.Apply(data)
	decodeMob(m, data)
}

func (striderType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
