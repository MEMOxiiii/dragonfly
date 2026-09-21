package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewEvocationIllager creates a new EvocationIllager at the position passed.
func NewEvocationIllager(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(EvocationIllagerType, evocationIllagerConf)
}

var evocationIllagerConf = MobConfig{
	MaxHealth:       24,
	Speed:           0.5,
	ExperienceDrops: [2]int{10, 10},
	LavaDamage:      4,
	Breathable:      Breathable{Air: true, TotalSupply: 15},
	Drops: LootTable{
		Always("minecraft:emerald", 0, 1, Looting(1)),
		Always("minecraft:totem_of_undying", 1, 1),
	},
	Goals: []Goal{
		Float(0),
		HurtByTarget(1),
		NearestAttackableTarget(2, 16, 20, false),
		AvoidMobType(5, 8, 0.6, 0.6, "minecraft:player"),
		RandomStroll(8, 0.6),
		LookAtPlayer(9, 3, 1, 2),
	},
}

// EvocationIllagerType is a world.EntityType implementation for evocation illagers.
var EvocationIllagerType evocationIllagerType

type evocationIllagerType struct{}

func (evocationIllagerType) EncodeEntity() string { return "minecraft:evocation_illager" }
func (evocationIllagerType) SpawnEggName() string { return "" }
func (evocationIllagerType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.3, 0, -0.3, 0.3, 1.9, 0.3)
}

func (evocationIllagerType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (evocationIllagerType) Apply(data *world.EntityData) { evocationIllagerConf.Apply(data) }

func (evocationIllagerType) DecodeNBT(m map[string]any, data *world.EntityData) {
	evocationIllagerConf.Apply(data)
	decodeMob(m, data)
}

func (evocationIllagerType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
