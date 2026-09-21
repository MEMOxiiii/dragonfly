package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewPig creates a new Pig at the position passed.
func NewPig(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(PigType, pigConf)
}

var pigConf = MobConfig{
	MaxHealth:       10,
	Speed:           0.25,
	ExperienceDrops: [2]int{1, 3},
	LavaDamage:      4,
	Breathable:      Breathable{Air: true, TotalSupply: 15},
	BreedItems:      []string{"minecraft:carrot", "minecraft:beetroot", "minecraft:potato"},
	Drops: LootTable{
		Always("minecraft:porkchop", 1, 3, Smelted("minecraft:cooked_porkchop"), Looting(1)),
	},
	Goals: []Goal{
		Float(2),
		Panic(3, 1.25),
		Breed(4, 1),
		Tempt(5, 1.2, 10, "minecraft:potato", "minecraft:carrot", "minecraft:beetroot", "minecraft:carrot_on_a_stick"),
		RandomStroll(7, 1),
		LookAtPlayer(8, 6, 1, 2),
		RandomLookAround(9),
	},
}

// PigType is a world.EntityType implementation for pigs.
var PigType pigType

type pigType struct{}

func (pigType) EncodeEntity() string { return "minecraft:pig" }
func (pigType) SpawnEggName() string { return "pig" }
func (pigType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.45, 0, -0.45, 0.45, 0.9, 0.45)
}

func (pigType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (pigType) Apply(data *world.EntityData) { pigConf.Apply(data) }

func (pigType) DecodeNBT(m map[string]any, data *world.EntityData) {
	pigConf.Apply(data)
	decodeMob(m, data)
}

func (pigType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
