package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewTropicalfish creates a new Tropicalfish at the position passed.
func NewTropicalfish(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(TropicalfishType, tropicalfishConf)
}

var tropicalfishConf = MobConfig{
	MaxHealth:       6,
	Speed:           0.12,
	Travel:          TravelSwim,
	ExperienceDrops: [2]int{1, 3},
	LavaDamage:      4,
	Breathable:      Breathable{Water: true, TotalSupply: 15},
	Colours:         []VariantChoice{{0, 1}, {1, 1}, {2, 1}, {3, 1}, {4, 1}, {5, 1}, {6, 1}, {7, 1}, {8, 1}, {9, 1}, {10, 1}, {11, 1}, {12, 1}, {13, 1}, {14, 1}},
	Variants:        []VariantChoice{{0, 1}, {1, 1}},
	Drops: LootTable{
		Always("minecraft:tropical_fish", 1, 1),
		Always("minecraft:bone", 1, 1, Looting(2)),
	},
	Goals: []Goal{
		AvoidMobType(1, 6, 1.5, 2, "minecraft:player"),
		RandomSwim(3, 16, 4, 0, 1),
	},
}

// TropicalfishType is a world.EntityType implementation for tropicalfishs.
var TropicalfishType tropicalfishType

type tropicalfishType struct{}

func (tropicalfishType) EncodeEntity() string { return "minecraft:tropicalfish" }
func (tropicalfishType) SpawnEggName() string { return "" }
func (tropicalfishType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.2, 0, -0.2, 0.2, 0.4, 0.2)
}

func (tropicalfishType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (tropicalfishType) Apply(data *world.EntityData) { tropicalfishConf.Apply(data) }

func (tropicalfishType) DecodeNBT(m map[string]any, data *world.EntityData) {
	tropicalfishConf.Apply(data)
	decodeMob(m, data)
}

func (tropicalfishType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
