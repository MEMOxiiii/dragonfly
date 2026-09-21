package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewPufferfish creates a new Pufferfish at the position passed.
func NewPufferfish(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(PufferfishType, pufferfishConf)
}

var pufferfishConf = MobConfig{
	MaxHealth:       6,
	Speed:           0.13,
	Travel:          TravelSwim,
	ExperienceDrops: [2]int{1, 3},
	LavaDamage:      4,
	Breathable:      Breathable{Water: true, TotalSupply: 15},
	Variants:        []VariantChoice{{0, 1}},
	Drops: LootTable{
		Always("minecraft:pufferfish", 1, 1),
		Always("minecraft:bone", 1, 1, Looting(2)),
	},
	Goals: []Goal{
		RandomSwim(3, 16, 4, 0, 1),
	},
}

// PufferfishType is a world.EntityType implementation for pufferfishs.
var PufferfishType pufferfishType

type pufferfishType struct{}

func (pufferfishType) EncodeEntity() string { return "minecraft:pufferfish" }
func (pufferfishType) SpawnEggName() string { return "pufferfish" }
func (pufferfishType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.4, 0, -0.4, 0.4, 0.8, 0.4)
}

func (pufferfishType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (pufferfishType) Apply(data *world.EntityData) { pufferfishConf.Apply(data) }

func (pufferfishType) DecodeNBT(m map[string]any, data *world.EntityData) {
	pufferfishConf.Apply(data)
	decodeMob(m, data)
}

func (pufferfishType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
