package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewGlowSquid creates a new GlowSquid at the position passed.
func NewGlowSquid(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(GlowSquidType, glowSquidConf)
}

var glowSquidConf = MobConfig{
	MaxHealth:       10,
	Speed:           0.2,
	ExperienceDrops: [2]int{1, 3},
	LavaDamage:      4,
	Breathable:      Breathable{Water: true, TotalSupply: 15},
	Goals: []Goal{
		SquidSwim(1),
	},
}

// GlowSquidType is a world.EntityType implementation for glow squids.
var GlowSquidType glowSquidType

type glowSquidType struct{}

func (glowSquidType) EncodeEntity() string { return "minecraft:glow_squid" }
func (glowSquidType) SpawnEggName() string { return "glow_squid" }
func (glowSquidType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.4, 0, -0.4, 0.4, 0.8, 0.4)
}

func (glowSquidType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (glowSquidType) Apply(data *world.EntityData) { glowSquidConf.Apply(data) }

func (glowSquidType) DecodeNBT(m map[string]any, data *world.EntityData) {
	glowSquidConf.Apply(data)
	decodeMob(m, data)
}

func (glowSquidType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
