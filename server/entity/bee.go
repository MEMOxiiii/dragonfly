package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewBee creates a new Bee at the position passed.
func NewBee(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(BeeType, beeConf)
}

var beeConf = MobConfig{
	MaxHealth:       10,
	Speed:           0.15,
	Travel:          TravelHover,
	ExperienceDrops: [2]int{1, 3},
	LavaDamage:      4,
	Breathable:      Breathable{Air: true},
	BreedItems:      []string{"minecraft:poppy", "minecraft:blue_orchid", "minecraft:allium", "minecraft:azure_bluet", "minecraft:red_tulip", "minecraft:orange_tulip", "minecraft:white_tulip", "minecraft:pink_tulip", "minecraft:oxeye_daisy", "minecraft:cornflower", "minecraft:lily_of_the_valley", "minecraft:dandelion", "minecraft:wither_rose", "minecraft:sunflower", "minecraft:lilac", "minecraft:rose_bush", "minecraft:peony", "minecraft:flowering_azalea", "minecraft:azalea_leaves_flowered", "minecraft:mangrove_propagule", "minecraft:pitcher_plant", "minecraft:torchflower", "minecraft:cherry_leaves", "minecraft:pink_petals", "minecraft:wildflowers", "minecraft:cactus_flower", "minecraft:chorus_flower", "minecraft:spore_blossom"},
	Goals: []Goal{
		HurtByTarget(1),
		Breed(4, 1),
		Tempt(5, 1.25, 8, "minecraft:poppy", "minecraft:blue_orchid", "minecraft:allium", "minecraft:azure_bluet", "minecraft:red_tulip", "minecraft:orange_tulip", "minecraft:white_tulip", "minecraft:pink_tulip", "minecraft:oxeye_daisy", "minecraft:cornflower", "minecraft:lily_of_the_valley", "minecraft:dandelion", "minecraft:wither_rose", "minecraft:sunflower", "minecraft:lilac", "minecraft:rose_bush", "minecraft:peony", "minecraft:flowering_azalea", "minecraft:azalea_leaves_flowered", "minecraft:mangrove_propagule", "minecraft:pitcher_plant", "minecraft:torchflower", "minecraft:cherry_leaves", "minecraft:pink_petals", "minecraft:open_eyeblossom", "minecraft:wildflowers", "minecraft:cactus_flower", "minecraft:chorus_flower", "minecraft:spore_blossom"),
		RandomFly(12, 8, 8, -1, 1, 4, 1),
		Float(19),
	},
}

// BeeType is a world.EntityType implementation for bees.
var BeeType beeType

type beeType struct{}

func (beeType) EncodeEntity() string { return "minecraft:bee" }
func (beeType) SpawnEggName() string { return "bee" }
func (beeType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.275, 0, -0.275, 0.275, 0.5, 0.275)
}

func (beeType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (beeType) Apply(data *world.EntityData) { beeConf.Apply(data) }

func (beeType) DecodeNBT(m map[string]any, data *world.EntityData) {
	beeConf.Apply(data)
	decodeMob(m, data)
}

func (beeType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
