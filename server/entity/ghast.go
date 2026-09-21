package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewGhast creates a new Ghast at the position passed.
func NewGhast(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(GhastType, ghastConf)
}

var ghastConf = MobConfig{
	MaxHealth:       10,
	Speed:           0.03,
	Travel:          TravelFloat,
	ExperienceDrops: [2]int{5, 5},
	FireImmune:      true,
	Breathable:      Breathable{Air: true, TotalSupply: 15},
	Drops: LootTable{
		Always("minecraft:ghast_tear", 0, 1, Looting(1)),
		Always("minecraft:gunpowder", 0, 2, Looting(1)),
	},
	Goals: []Goal{
		Float(0),
		HurtByTarget(1),
		RandomFly(2, 10, 7, 0, 0, 0, 1),
		NearestAttackableTarget(2, 16, 28, false),
	},
}

// GhastType is a world.EntityType implementation for ghasts.
var GhastType ghastType

type ghastType struct{}

func (ghastType) EncodeEntity() string { return "minecraft:ghast" }
func (ghastType) SpawnEggName() string { return "ghast" }
func (ghastType) BBox(world.Entity) cube.BBox {
	return cube.Box(-2, 0, -2, 2, 4, 2)
}

func (ghastType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (ghastType) Apply(data *world.EntityData) { ghastConf.Apply(data) }

func (ghastType) DecodeNBT(m map[string]any, data *world.EntityData) {
	ghastConf.Apply(data)
	decodeMob(m, data)
}

func (ghastType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
