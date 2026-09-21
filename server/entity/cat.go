package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewCat creates a new Cat at the position passed.
func NewCat(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(CatType, catConf)
}

var catConf = MobConfig{
	MaxHealth:       20,
	Speed:           0.3,
	ExperienceDrops: [2]int{1, 3},
	LavaDamage:      4,
	Breathable:      Breathable{Air: true, TotalSupply: 15},
	Healable:        Healable{Items: map[string]float64{"minecraft:cod": 2, "minecraft:salmon": 2}},
	Tameable:        Tameable{Items: []string{"minecraft:cod", "minecraft:salmon"}, Probability: 0.33},
	Colours:         []VariantChoice{{14, 1}},
	Variants:        []VariantChoice{{0, 15}, {1, 15}, {2, 15}, {3, 15}, {4, 15}, {5, 15}, {6, 15}, {7, 15}, {8, 15}, {9, 15}, {10, 15}},
	BreedItems:      []string{"minecraft:cod", "minecraft:salmon"},
	Drops: LootTable{
		Always("minecraft:string", 0, 2),
	},
	Goals: []Goal{
		Float(0),
		Panic(1, 1.25),
		NearestAttackableTarget(1, 16, 8, true),
		Breed(3, 1),
		Tempt(5, 0.5, 16, "minecraft:cod", "minecraft:salmon"),
		AvoidMobType(6, 10, 0.8, 1.33, "minecraft:player"),
		RandomStroll(8, 0.8),
		LookAtPlayer(9, 8, 1, 2),
	},
}

// CatType is a world.EntityType implementation for cats.
var CatType catType

type catType struct{}

func (catType) EncodeEntity() string { return "minecraft:cat" }
func (catType) SpawnEggName() string { return "cat" }
func (catType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.3, 0, -0.3, 0.3, 0.7, 0.3)
}

func (catType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (catType) Apply(data *world.EntityData) { catConf.Apply(data) }

func (catType) DecodeNBT(m map[string]any, data *world.EntityData) {
	catConf.Apply(data)
	decodeMob(m, data)
}

func (catType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
