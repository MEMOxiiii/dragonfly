package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewRabbit creates a new Rabbit at the position passed.
func NewRabbit(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(RabbitType, rabbitConf)
}

var rabbitConf = MobConfig{
	MaxHealth:       3,
	Speed:           0.3,
	ExperienceDrops: [2]int{1, 3},
	LavaDamage:      4,
	Breathable:      Breathable{Air: true, TotalSupply: 15},
	Variants:        []VariantChoice{{0, 50}, {2, 40}, {5, 10}},
	BreedItems:      []string{"minecraft:golden_carrot", "minecraft:carrot", "minecraft:dandelion"},
	Drops: LootTable{
		Always("minecraft:rabbit_hide", 0, 1, Looting(1)),
		Always("minecraft:rabbit", 0, 1, Smelted("minecraft:cooked_rabbit")),
		Always("minecraft:rabbit_foot", 1, 1),
	},
	Goals: []Goal{
		Float(0),
		Panic(1, 2.2),
		Breed(2, 1),
		Tempt(3, 1, 10, "minecraft:golden_carrot", "minecraft:carrot", "minecraft:dandelion"),
		AvoidMobType(4, 8, 1.5, 1.8, "minecraft:monster", "minecraft:player", "minecraft:wolf"),
		RandomStroll(6, 0.6),
		LookAtPlayer(11, 8, 1, 2),
	},
}

// RabbitType is a world.EntityType implementation for rabbits.
var RabbitType rabbitType

type rabbitType struct{}

func (rabbitType) EncodeEntity() string { return "minecraft:rabbit" }
func (rabbitType) SpawnEggName() string { return "rabbit" }
func (rabbitType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.335, 0, -0.335, 0.335, 0.67, 0.335)
}

func (rabbitType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (rabbitType) Apply(data *world.EntityData) { rabbitConf.Apply(data) }

func (rabbitType) DecodeNBT(m map[string]any, data *world.EntityData) {
	rabbitConf.Apply(data)
	decodeMob(m, data)
}

func (rabbitType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
