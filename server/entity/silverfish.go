package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewSilverfish creates a new Silverfish at the position passed.
func NewSilverfish(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(SilverfishType, silverfishConf)
}

var silverfishConf = MobConfig{
	MaxHealth:       8,
	Speed:           0.25,
	ExperienceDrops: [2]int{5, 5},
	LavaDamage:      4,
	Breathable:      Breathable{Air: true, TotalSupply: 15},
	Goals: []Goal{
		Float(1),
		HurtByTarget(1),
		NearestAttackableTarget(2, 16, 8, false),
	},
}

// SilverfishType is a world.EntityType implementation for silverfishs.
var SilverfishType silverfishType

type silverfishType struct{}

func (silverfishType) EncodeEntity() string { return "minecraft:silverfish" }
func (silverfishType) SpawnEggName() string { return "silverfish" }
func (silverfishType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.2, 0, -0.2, 0.2, 0.3, 0.2)
}

func (silverfishType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (silverfishType) Apply(data *world.EntityData) { silverfishConf.Apply(data) }

func (silverfishType) DecodeNBT(m map[string]any, data *world.EntityData) {
	silverfishConf.Apply(data)
	decodeMob(m, data)
}

func (silverfishType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
