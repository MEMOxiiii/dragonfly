package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewHorse creates a new Horse at the position passed.
func NewHorse(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(HorseType, horseConf)
}

var horseConf = MobConfig{
	MaxHealth:       22.5,
	Speed:           0.225,
	ExperienceDrops: [2]int{1, 3},
	LavaDamage:      4,
	Breathable:      Breathable{Air: true, TotalSupply: 15},
	Healable:        Healable{Items: map[string]float64{"minecraft:wheat": 2, "minecraft:sugar": 1, "minecraft:hay_block": 20, "minecraft:apple": 3, "minecraft:golden_carrot": 4, "minecraft:golden_apple": 10, "minecraft:enchanted_golden_apple": 10}},
	Variants:        []VariantChoice{{0, 1}, {1, 1}, {2, 1}, {3, 1}, {4, 1}, {5, 1}, {6, 1}},
	BreedItems:      []string{"minecraft:golden_carrot", "minecraft:golden_apple", "minecraft:enchanted_golden_apple"},
	Drops: LootTable{
		Always("minecraft:leather", 0, 2, Looting(1)),
	},
	Goals: []Goal{
		Float(0),
		Breed(2, 1),
		Panic(3, 1.2),
		RandomStroll(6, 0.7),
		LookAtPlayer(7, 6, 1, 2),
		RandomLookAround(8),
	},
}

// HorseType is a world.EntityType implementation for horses.
var HorseType horseType

type horseType struct{}

func (horseType) EncodeEntity() string { return "minecraft:horse" }
func (horseType) SpawnEggName() string { return "horse" }
func (horseType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.7, 0, -0.7, 0.7, 1.6, 0.7)
}

func (horseType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (horseType) Apply(data *world.EntityData) { horseConf.Apply(data) }

func (horseType) DecodeNBT(m map[string]any, data *world.EntityData) {
	horseConf.Apply(data)
	decodeMob(m, data)
}

func (horseType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
