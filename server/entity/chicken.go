package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewChicken creates a new Chicken at the position passed.
func NewChicken(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(ChickenType, chickenConf)
}

var chickenConf = MobConfig{
	MaxHealth:       4,
	Speed:           0.25,
	ExperienceDrops: [2]int{1, 3},
	LavaDamage:      4,
	Breathable:      Breathable{Air: true, TotalSupply: 15},
	BreedItems:      []string{"minecraft:wheat_seeds", "minecraft:beetroot_seeds", "minecraft:melon_seeds", "minecraft:pumpkin_seeds"},
	Drops: LootTable{
		Always("minecraft:feather", 0, 2, Looting(1)),
		Always("minecraft:chicken", 1, 1, Smelted("minecraft:cooked_chicken"), Looting(1)),
	},
	Goals: []Goal{
		Float(0),
		Panic(1, 1.5),
		Breed(3, 1),
		Tempt(4, 1, 10, "minecraft:wheat_seeds", "minecraft:beetroot_seeds", "minecraft:melon_seeds", "minecraft:pumpkin_seeds"),
		RandomStroll(6, 1),
		LookAtPlayer(7, 6, 1, 2),
		RandomLookAround(8),
	},
}

// ChickenType is a world.EntityType implementation for chickens.
var ChickenType chickenType

type chickenType struct{}

func (chickenType) EncodeEntity() string { return "minecraft:chicken" }
func (chickenType) SpawnEggName() string { return "chicken" }
func (chickenType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.3, 0, -0.3, 0.3, 0.8, 0.3)
}

func (chickenType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (chickenType) Apply(data *world.EntityData) { chickenConf.Apply(data) }

func (chickenType) DecodeNBT(m map[string]any, data *world.EntityData) {
	chickenConf.Apply(data)
	decodeMob(m, data)
}

func (chickenType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
