package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewSheep creates a new Sheep at the position passed.
func NewSheep(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(SheepType, sheepConf)
}

var sheepConf = MobConfig{
	MaxHealth:       8,
	Speed:           0.25,
	ExperienceDrops: [2]int{1, 3},
	LavaDamage:      4,
	Breathable:      Breathable{Air: true, TotalSupply: 15},
	Colours:         []VariantChoice{{0, 81836}, {15, 5000}, {7, 5000}, {8, 5000}, {12, 3000}, {6, 164}},
	ShearDrops: LootTable{
		Always("minecraft:white_wool", 1, 3),
	},
	BreedItems: []string{"minecraft:wheat"},
	Drops: LootTable{
		Always("minecraft:white_wool", 1, 1),
		Always("minecraft:mutton", 1, 2, Smelted("minecraft:cooked_mutton"), Looting(1)),
	},
	Goals: []Goal{
		Float(0),
		Panic(1, 1.25),
		Breed(3, 1),
		Tempt(4, 1.25, 10, "minecraft:wheat"),
		FollowParent(5, 1.1),
		RandomStroll(7, 0.8),
		LookAtPlayer(8, 6, 1, 2),
		RandomLookAround(9),
	},
}

// SheepType is a world.EntityType implementation for sheeps.
var SheepType sheepType

type sheepType struct{}

func (sheepType) EncodeEntity() string { return "minecraft:sheep" }
func (sheepType) SpawnEggName() string { return "sheep" }
func (sheepType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.45, 0, -0.45, 0.45, 1.3, 0.45)
}

func (sheepType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (sheepType) Apply(data *world.EntityData) { sheepConf.Apply(data) }

func (sheepType) DecodeNBT(m map[string]any, data *world.EntityData) {
	sheepConf.Apply(data)
	decodeMob(m, data)
}

func (sheepType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
