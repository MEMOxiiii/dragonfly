package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewCaveSpider creates a new CaveSpider at the position passed.
func NewCaveSpider(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(CaveSpiderType, caveSpiderConf)
}

var caveSpiderConf = MobConfig{
	MaxHealth:       12,
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

// CaveSpiderType is a world.EntityType implementation for cave spiders.
var CaveSpiderType caveSpiderType

type caveSpiderType struct{}

func (caveSpiderType) EncodeEntity() string { return "minecraft:cave_spider" }
func (caveSpiderType) SpawnEggName() string { return "cave_spider" }
func (caveSpiderType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.35, 0, -0.35, 0.35, 0.5, 0.35)
}

func (caveSpiderType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (caveSpiderType) Apply(data *world.EntityData) { caveSpiderConf.Apply(data) }

func (caveSpiderType) DecodeNBT(m map[string]any, data *world.EntityData) {
	caveSpiderConf.Apply(data)
	decodeMob(m, data)
}

func (caveSpiderType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
