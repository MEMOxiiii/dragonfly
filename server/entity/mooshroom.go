package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewMooshroom creates a new Mooshroom at the position passed.
func NewMooshroom(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(MooshroomType, mooshroomConf)
}

var mooshroomConf = MobConfig{
	MaxHealth:       10,
	Speed:           0.25,
	ExperienceDrops: [2]int{1, 3},
	LavaDamage:      4,
	Breathable:      Breathable{Air: true, TotalSupply: 15},
	Variants:        []VariantChoice{{0, 95}, {0, 5}},
	BreedItems:      []string{"minecraft:wheat"},
	Drops: LootTable{
		Always("minecraft:leather", 0, 2, Looting(1)),
		Always("minecraft:beef", 1, 3, Smelted("minecraft:cooked_beef"), Looting(1)),
	},
	Goals: []Goal{
		Float(0),
		Panic(1, 1.25),
		Breed(3, 1),
		Tempt(4, 1.25, 10, "minecraft:wheat"),
		FollowParent(5, 1.1),
		RandomStroll(6, 0.8),
		LookAtPlayer(7, 6, 1, 2),
		RandomLookAround(9),
	},
}

// MooshroomType is a world.EntityType implementation for mooshrooms.
var MooshroomType mooshroomType

type mooshroomType struct{}

func (mooshroomType) EncodeEntity() string { return "minecraft:mooshroom" }
func (mooshroomType) SpawnEggName() string { return "mooshroom" }
func (mooshroomType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.45, 0, -0.45, 0.45, 1.3, 0.45)
}

func (mooshroomType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (mooshroomType) Apply(data *world.EntityData) { mooshroomConf.Apply(data) }

func (mooshroomType) DecodeNBT(m map[string]any, data *world.EntityData) {
	mooshroomConf.Apply(data)
	decodeMob(m, data)
}

func (mooshroomType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
