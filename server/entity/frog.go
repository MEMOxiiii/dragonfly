package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewFrog creates a new Frog at the position passed.
func NewFrog(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(FrogType, frogConf)
}

var frogConf = MobConfig{
	MaxHealth:       10,
	Speed:           0.15,
	Travel:          TravelSwim,
	ExperienceDrops: [2]int{1, 3},
	LavaDamage:      4,
	Breathable:      Breathable{Air: true, Water: true, TotalSupply: 15},
	Variants:        []VariantChoice{{1, 1}},
	Goals: []Goal{
		Panic(1, 2),
		Breed(4, 1),
		Tempt(5, 1.25, 10, "minecraft:slime_ball"),
		NearestAttackableTarget(8, 16, 16, false),
		RandomStroll(11, 1),
		LookAtPlayer(12, 6, 2, 4),
	},
}

// FrogType is a world.EntityType implementation for frogs.
var FrogType frogType

type frogType struct{}

func (frogType) EncodeEntity() string { return "minecraft:frog" }
func (frogType) SpawnEggName() string { return "frog" }
func (frogType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.25, 0, -0.25, 0.25, 0.55, 0.25)
}

func (frogType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (frogType) Apply(data *world.EntityData) { frogConf.Apply(data) }

func (frogType) DecodeNBT(m map[string]any, data *world.EntityData) {
	frogConf.Apply(data)
	decodeMob(m, data)
}

func (frogType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
