package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewCow creates a new Cow at the position passed.
func NewCow(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(CowType, cowConf)
}

var cowConf = MobConfig{
	MaxHealth:       10,
	Speed:           0.25,
	ExperienceDrops: [2]int{1, 3},
	LavaDamage:      4,
	Breathable:      Breathable{Air: true, TotalSupply: 15},
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

// CowType is a world.EntityType implementation for cows.
var CowType cowType

type cowType struct{}

func (cowType) EncodeEntity() string { return "minecraft:cow" }
func (cowType) SpawnEggName() string { return "cow" }
func (cowType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.45, 0, -0.45, 0.45, 1.3, 0.45)
}

func (cowType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (cowType) Apply(data *world.EntityData) { cowConf.Apply(data) }

func (cowType) DecodeNBT(m map[string]any, data *world.EntityData) {
	cowConf.Apply(data)
	decodeMob(m, data)
}

func (cowType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
