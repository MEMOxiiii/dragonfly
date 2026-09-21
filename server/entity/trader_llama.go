package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewTraderLlama creates a new TraderLlama at the position passed.
func NewTraderLlama(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(TraderLlamaType, traderLlamaConf)
}

var traderLlamaConf = MobConfig{
	MaxHealth:       22.5,
	Speed:           0.25,
	ExperienceDrops: [2]int{1, 3},
	LavaDamage:      4,
	Breathable:      Breathable{Air: true, TotalSupply: 15},
	Healable:        Healable{Items: map[string]float64{"minecraft:wheat": 2, "minecraft:hay_block": 10}},
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
		Tempt(5, 1.2, 10, "minecraft:hay_block"),
		RandomStroll(6, 0.7),
		LookAtPlayer(7, 6, 2, 4),
		RandomLookAround(8),
	},
}

// TraderLlamaType is a world.EntityType implementation for trader llamas.
var TraderLlamaType traderLlamaType

type traderLlamaType struct{}

func (traderLlamaType) EncodeEntity() string { return "minecraft:trader_llama" }
func (traderLlamaType) SpawnEggName() string { return "trader_llama" }
func (traderLlamaType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.45, 0, -0.45, 0.45, 1.87, 0.45)
}

func (traderLlamaType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (traderLlamaType) Apply(data *world.EntityData) { traderLlamaConf.Apply(data) }

func (traderLlamaType) DecodeNBT(m map[string]any, data *world.EntityData) {
	traderLlamaConf.Apply(data)
	decodeMob(m, data)
}

func (traderLlamaType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
