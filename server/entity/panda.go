package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewPanda creates a new Panda at the position passed.
func NewPanda(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(PandaType, pandaConf)
}

var pandaConf = MobConfig{
	MaxHealth:       20,
	Speed:           0.15,
	ExperienceDrops: [2]int{1, 4},
	LavaDamage:      4,
	Breathable:      Breathable{Air: true, TotalSupply: 15},
	Variants:        []VariantChoice{{0, 1}},
	BreedItems:      []string{"minecraft:bamboo"},
	Drops: LootTable{
		Always("minecraft:bamboo", 0, 2, Looting(1)),
	},
	Goals: []Goal{
		Float(0),
		HurtByTarget(1),
		Panic(1, 1.25),
		MeleeAttack(2, 2, 1),
		Breed(3, 1),
		Tempt(4, 1.25, 10, "minecraft:bamboo"),
		LookAtPlayer(8, 6, 1, 2),
		RandomLookAround(9),
		RandomStroll(14, 0.8),
	},
}

// PandaType is a world.EntityType implementation for pandas.
var PandaType pandaType

type pandaType struct{}

func (pandaType) EncodeEntity() string { return "minecraft:panda" }
func (pandaType) SpawnEggName() string { return "panda" }
func (pandaType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.85, 0, -0.85, 0.85, 1.5, 0.85)
}

func (pandaType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (pandaType) Apply(data *world.EntityData) { pandaConf.Apply(data) }

func (pandaType) DecodeNBT(m map[string]any, data *world.EntityData) {
	pandaConf.Apply(data)
	decodeMob(m, data)
}

func (pandaType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
