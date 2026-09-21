package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewZombieVillager creates a new ZombieVillager at the position passed.
func NewZombieVillager(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(ZombieVillagerType, zombieVillagerConf)
}

var zombieVillagerConf = MobConfig{
	MaxHealth:       20,
	Speed:           0.23,
	ExperienceDrops: [2]int{5, 5},
	LavaDamage:      4,
	Breathable:      Breathable{Air: true, Water: true, TotalSupply: 15},
	Variants:        []VariantChoice{{0, 5}, {0, 5}, {0, 5}, {0, 5}, {1, 20}, {1, 20}, {2, 20}, {3, 6}, {3, 6}, {3, 6}, {4, 10}, {4, 10}},
	BurnsInDaylight: true,
	Drops: LootTable{
		Always("minecraft:rotten_flesh", 0, 2, Looting(1)),
		OneOf(Entry("minecraft:iron_ingot", 1, 1), Entry("minecraft:carrot", 1, 1), Entry("minecraft:potato", 1, 1)),
	},
	Goals: []Goal{
		HurtByTarget(1),
		NearestAttackableTarget(2, 16, 35, true),
		MeleeAttack(6, 3, 1),
		RandomStroll(8, 1),
		LookAtPlayer(9, 6, 1, 2),
		RandomLookAround(9),
	},
}

// ZombieVillagerType is a world.EntityType implementation for zombie villagers.
var ZombieVillagerType zombieVillagerType

type zombieVillagerType struct{}

func (zombieVillagerType) EncodeEntity() string { return "minecraft:zombie_villager" }
func (zombieVillagerType) SpawnEggName() string { return "zombie_villager" }
func (zombieVillagerType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.3, 0, -0.3, 0.3, 1.9, 0.3)
}

func (zombieVillagerType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (zombieVillagerType) Apply(data *world.EntityData) { zombieVillagerConf.Apply(data) }

func (zombieVillagerType) DecodeNBT(m map[string]any, data *world.EntityData) {
	zombieVillagerConf.Apply(data)
	decodeMob(m, data)
}

func (zombieVillagerType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
