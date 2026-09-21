package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewPiglinBrute creates a new PiglinBrute at the position passed.
func NewPiglinBrute(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(PiglinBruteType, piglinBruteConf)
}

var piglinBruteConf = MobConfig{
	MaxHealth:       50,
	Speed:           0.35,
	ExperienceDrops: [2]int{20, 20},
	LavaDamage:      4,
	Breathable:      Breathable{Air: true, TotalSupply: 15},
	Variants:        []VariantChoice{{1, 1}},
	Goals: []Goal{
		HurtByTarget(1),
		MeleeAttack(4, 7, 1),
		RandomStroll(7, 0.6),
		LookAtPlayer(8, 8, 1, 2),
		RandomLookAround(9),
	},
}

// PiglinBruteType is a world.EntityType implementation for piglin brutes.
var PiglinBruteType piglinBruteType

type piglinBruteType struct{}

func (piglinBruteType) EncodeEntity() string { return "minecraft:piglin_brute" }
func (piglinBruteType) SpawnEggName() string { return "piglin_brute" }
func (piglinBruteType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.3, 0, -0.3, 0.3, 1.9, 0.3)
}

func (piglinBruteType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (piglinBruteType) Apply(data *world.EntityData) { piglinBruteConf.Apply(data) }

func (piglinBruteType) DecodeNBT(m map[string]any, data *world.EntityData) {
	piglinBruteConf.Apply(data)
	decodeMob(m, data)
}

func (piglinBruteType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
