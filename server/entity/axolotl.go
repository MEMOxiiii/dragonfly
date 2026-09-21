package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewAxolotl creates a new Axolotl at the position passed.
func NewAxolotl(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(AxolotlType, axolotlConf)
}

var axolotlConf = MobConfig{
	MaxHealth:       14,
	Speed:           0.2,
	Travel:          TravelSwim,
	ExperienceDrops: [2]int{1, 3},
	LavaDamage:      4,
	Breathable:      Breathable{Air: true, Water: true, TotalSupply: 15},
	Variants:        []VariantChoice{{1, 25}, {2, 25}, {0, 25}, {3, 25}},
	BreedItems:      []string{"minecraft:tropical_fish_bucket"},
	Goals: []Goal{
		Breed(1, 1),
		Tempt(2, 1.1, 10, "minecraft:tropical_fish_bucket"),
		NearestAttackableTarget(3, 20, 8, true),
		MeleeAttack(4, 2, 1),
		MoveToWater(6, 16, 1),
		RandomSwim(8, 30, 15, 0, 1),
		RandomStroll(9, 1),
		LookAtPlayer(10, 6, 2, 4),
	},
}

// AxolotlType is a world.EntityType implementation for axolotls.
var AxolotlType axolotlType

type axolotlType struct{}

func (axolotlType) EncodeEntity() string { return "minecraft:axolotl" }
func (axolotlType) SpawnEggName() string { return "axolotl" }
func (axolotlType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.375, 0, -0.375, 0.375, 0.42, 0.375)
}

func (axolotlType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (axolotlType) Apply(data *world.EntityData) { axolotlConf.Apply(data) }

func (axolotlType) DecodeNBT(m map[string]any, data *world.EntityData) {
	axolotlConf.Apply(data)
	decodeMob(m, data)
}

func (axolotlType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
