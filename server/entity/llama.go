package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewLlama creates a new Llama at the position passed.
func NewLlama(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(LlamaType, llamaConf)
}

var llamaConf = MobConfig{
	MaxHealth:       22.5,
	Speed:           0.25,
	ExperienceDrops: [2]int{1, 3},
	LavaDamage:      4,
	Breathable:      Breathable{Air: true, TotalSupply: 15},
	Healable:        Healable{Items: map[string]float64{"minecraft:wheat": 2, "minecraft:hay_block": 10}},
	Variants:        []VariantChoice{{0, 25}, {1, 25}, {2, 25}, {3, 25}},
	BreedItems:      []string{"minecraft:hay_block"},
	Drops: LootTable{
		Always("minecraft:leather", 0, 2, Looting(1)),
	},
	Goals: []Goal{
		Float(0),
		HurtByTarget(1),
		NearestAttackableTarget(2, 16, 10, false),
		Panic(4, 1.2),
		Breed(4, 1),
		RandomStroll(6, 0.7),
		LookAtPlayer(7, 6, 1, 2),
		RandomLookAround(8),
	},
}

// LlamaType is a world.EntityType implementation for llamas.
var LlamaType llamaType

type llamaType struct{}

func (llamaType) EncodeEntity() string { return "minecraft:llama" }
func (llamaType) SpawnEggName() string { return "llama" }
func (llamaType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.45, 0, -0.45, 0.45, 1.87, 0.45)
}

func (llamaType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (llamaType) Apply(data *world.EntityData) { llamaConf.Apply(data) }

func (llamaType) DecodeNBT(m map[string]any, data *world.EntityData) {
	llamaConf.Apply(data)
	decodeMob(m, data)
}

func (llamaType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
