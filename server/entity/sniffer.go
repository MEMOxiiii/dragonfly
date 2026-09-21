package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewSniffer creates a new Sniffer at the position passed.
func NewSniffer(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(SnifferType, snifferConf)
}

var snifferConf = MobConfig{
	MaxHealth:       14,
	Speed:           0.09,
	ExperienceDrops: [2]int{1, 3},
	LavaDamage:      4,
	Breathable:      Breathable{Air: true, TotalSupply: 15},
	Healable:        Healable{Items: map[string]float64{"minecraft:torchflower_seeds": 2}},
	BreedItems:      []string{"minecraft:torchflower_seeds"},
	Goals: []Goal{
		Float(0),
		Panic(1, 2),
		Breed(3, 1),
		Tempt(4, 1.25, 10, "minecraft:torchflower_seeds"),
		RandomStroll(7, 1),
		LookAtPlayer(8, 6, 2, 4),
		RandomLookAround(9),
	},
}

// SnifferType is a world.EntityType implementation for sniffers.
var SnifferType snifferType

type snifferType struct{}

func (snifferType) EncodeEntity() string { return "minecraft:sniffer" }
func (snifferType) SpawnEggName() string { return "sniffer" }
func (snifferType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.95, 0, -0.95, 0.95, 1.75, 0.95)
}

func (snifferType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (snifferType) Apply(data *world.EntityData) { snifferConf.Apply(data) }

func (snifferType) DecodeNBT(m map[string]any, data *world.EntityData) {
	snifferConf.Apply(data)
	decodeMob(m, data)
}

func (snifferType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
