package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewPhantom creates a new Phantom at the position passed.
func NewPhantom(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(PhantomType, phantomConf)
}

var phantomConf = MobConfig{
	MaxHealth:       20,
	Speed:           1.8,
	ExperienceDrops: [2]int{5, 5},
	LavaDamage:      4,
	Breathable:      Breathable{Air: true, TotalSupply: 15},
	BurnsInDaylight: true,
	Drops: LootTable{
		Always("minecraft:phantom_membrane", 0, 1, Looting(1)),
	},
	Goals: []Goal{
		AvoidMobType(0, 6, 1, 1, "minecraft:cat", "minecraft:ocelot"),
		NearestAttackableTarget(1, 64, 64, true),
	},
}

// PhantomType is a world.EntityType implementation for phantoms.
var PhantomType phantomType

type phantomType struct{}

func (phantomType) EncodeEntity() string { return "minecraft:phantom" }
func (phantomType) SpawnEggName() string { return "phantom" }
func (phantomType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.45, 0, -0.45, 0.45, 0.5, 0.45)
}

func (phantomType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (phantomType) Apply(data *world.EntityData) { phantomConf.Apply(data) }

func (phantomType) DecodeNBT(m map[string]any, data *world.EntityData) {
	phantomConf.Apply(data)
	decodeMob(m, data)
}

func (phantomType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
