package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewSpider creates a new Spider at the position passed.
func NewSpider(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(SpiderType, spiderConf)
}

var spiderConf = MobConfig{
	MaxHealth:       16,
	Speed:           0.3,
	Travel:          TravelClimb,
	ExperienceDrops: [2]int{5, 5},
	LavaDamage:      4,
	Breathable:      Breathable{Air: true, TotalSupply: 15},
	Drops: LootTable{
		Always("minecraft:string", 0, 2, Looting(1)),
		Always("minecraft:spider_eye", 0, 1, Looting(1)),
	},
	Goals: []Goal{
		Float(1),
		HurtByTarget(1),
		RandomStroll(6, 0.8),
		LookAtPlayer(7, 6, 2, 4),
		RandomLookAround(7),
	},
}

// SpiderType is a world.EntityType implementation for spiders.
var SpiderType spiderType

type spiderType struct{}

func (spiderType) EncodeEntity() string { return "minecraft:spider" }
func (spiderType) SpawnEggName() string { return "spider" }
func (spiderType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.7, 0, -0.7, 0.7, 0.9, 0.7)
}

func (spiderType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (spiderType) Apply(data *world.EntityData) { spiderConf.Apply(data) }

func (spiderType) DecodeNBT(m map[string]any, data *world.EntityData) {
	spiderConf.Apply(data)
	decodeMob(m, data)
}

func (spiderType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
