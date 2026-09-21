package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewDonkey creates a new Donkey at the position passed.
func NewDonkey(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(DonkeyType, donkeyConf)
}

var donkeyConf = MobConfig{
	MaxHealth:       22.5,
	Speed:           0.175,
	ExperienceDrops: [2]int{1, 3},
	LavaDamage:      4,
	Breathable:      Breathable{Air: true, TotalSupply: 15},
	Healable:        Healable{Items: map[string]float64{"minecraft:wheat": 2, "minecraft:sugar": 1, "minecraft:hay_block": 20, "minecraft:apple": 3, "minecraft:golden_carrot": 4, "minecraft:golden_apple": 10, "minecraft:enchanted_golden_apple": 10}},
	BreedItems:      []string{"minecraft:golden_carrot", "minecraft:golden_apple", "minecraft:enchanted_golden_apple"},
	Drops: LootTable{
		Always("minecraft:leather", 0, 2, Looting(1)),
	},
	Goals: []Goal{
		Float(0),
		Panic(1, 1.2),
		Breed(2, 1),
		RandomStroll(6, 0.7),
		LookAtPlayer(7, 6, 1, 2),
		RandomLookAround(8),
	},
}

// DonkeyType is a world.EntityType implementation for donkeys.
var DonkeyType donkeyType

type donkeyType struct{}

func (donkeyType) EncodeEntity() string { return "minecraft:donkey" }
func (donkeyType) SpawnEggName() string { return "donkey" }
func (donkeyType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.7, 0, -0.7, 0.7, 1.6, 0.7)
}

func (donkeyType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (donkeyType) Apply(data *world.EntityData) { donkeyConf.Apply(data) }

func (donkeyType) DecodeNBT(m map[string]any, data *world.EntityData) {
	donkeyConf.Apply(data)
	decodeMob(m, data)
}

func (donkeyType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
