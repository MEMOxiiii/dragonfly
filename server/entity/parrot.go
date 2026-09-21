package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewParrot creates a new Parrot at the position passed.
func NewParrot(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(ParrotType, parrotConf)
}

var parrotConf = MobConfig{
	MaxHealth:       6,
	Speed:           0.4,
	Travel:          TravelFly,
	ExperienceDrops: [2]int{1, 3},
	LavaDamage:      4,
	Breathable:      Breathable{Air: true, TotalSupply: 15},
	Healable:        Healable{Items: map[string]float64{"minecraft:cookie": 0}},
	Tameable:        Tameable{Items: []string{"minecraft:wheat_seeds", "minecraft:pumpkin_seeds", "minecraft:melon_seeds", "minecraft:beetroot_seeds"}, Probability: 0.33},
	Variants:        []VariantChoice{{0, 20}, {1, 20}, {2, 20}, {3, 20}, {4, 20}},
	Drops: LootTable{
		Always("minecraft:feather", 1, 2),
	},
	Goals: []Goal{
		Float(0),
		Panic(0, 1.25),
		LookAtPlayer(1, 8, 1, 2),
		RandomFly(2, 15, 1, 0, 0, 0, 1),
	},
}

// ParrotType is a world.EntityType implementation for parrots.
var ParrotType parrotType

type parrotType struct{}

func (parrotType) EncodeEntity() string { return "minecraft:parrot" }
func (parrotType) SpawnEggName() string { return "parrot" }
func (parrotType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.25, 0, -0.25, 0.25, 1, 0.25)
}

func (parrotType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (parrotType) Apply(data *world.EntityData) { parrotConf.Apply(data) }

func (parrotType) DecodeNBT(m map[string]any, data *world.EntityData) {
	parrotConf.Apply(data)
	decodeMob(m, data)
}

func (parrotType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
