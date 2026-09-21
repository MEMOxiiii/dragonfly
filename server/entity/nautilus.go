package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewNautilus creates a new Nautilus at the position passed.
func NewNautilus(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(NautilusType, nautilusConf)
}

var nautilusConf = MobConfig{
	MaxHealth:           15,
	Travel:              TravelSwim,
	ExperienceDrops:     [2]int{1, 3},
	KnockBackResistance: 0.3,
	LavaDamage:          4,
	Breathable:          Breathable{Water: true, TotalSupply: 15},
	BreedItems:          []string{"minecraft:pufferfish_bucket", "minecraft:cod_bucket", "minecraft:salmon_bucket", "minecraft:tropical_fish_bucket", "minecraft:pufferfish", "minecraft:cod", "minecraft:salmon", "minecraft:tropical_fish", "minecraft:cooked_cod", "minecraft:cooked_salmon"},
	Goals: []Goal{
		HurtByTarget(1),
		Breed(3, 1),
		RandomSwim(6, 16, 4, 0, 1.5),
	},
}

// NautilusType is a world.EntityType implementation for nautiluss.
var NautilusType nautilusType

type nautilusType struct{}

func (nautilusType) EncodeEntity() string { return "minecraft:nautilus" }
func (nautilusType) SpawnEggName() string { return "nautilus" }
func (nautilusType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.4375, 0, -0.4375, 0.4375, 0.95, 0.4375)
}

func (nautilusType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (nautilusType) Apply(data *world.EntityData) { nautilusConf.Apply(data) }

func (nautilusType) DecodeNBT(m map[string]any, data *world.EntityData) {
	nautilusConf.Apply(data)
	decodeMob(m, data)
}

func (nautilusType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
