package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewZombie creates a new Zombie at the position passed.
func NewZombie(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(ZombieType, zombieConf)
}

var zombieConf = MobConfig{
	MaxHealth:       20,
	Speed:           0.23,
	ExperienceDrops: [2]int{5, 5},
	LavaDamage:      4,
	Breathable:      Breathable{Air: true, Water: true, TotalSupply: 15},
	BurnsInDaylight: true,
	Drops: LootTable{
		Always("minecraft:rotten_flesh", 0, 2, Looting(1)),
		OneOf(Entry("minecraft:iron_ingot", 1, 1), Entry("minecraft:carrot", 1, 1), Entry("minecraft:potato", 1, 1)),
	},
	Goals: []Goal{
		HurtByTarget(1),
		NearestAttackableTarget(2, 25, 35, true),
		MeleeAttack(3, 3, 1),
		RandomStroll(6, 1),
		LookAtPlayer(7, 6, 1, 2),
		RandomLookAround(7),
	},
}

// ZombieType is a world.EntityType implementation for zombies.
var ZombieType zombieType

type zombieType struct{}

func (zombieType) EncodeEntity() string { return "minecraft:zombie" }
func (zombieType) SpawnEggName() string { return "zombie" }
func (zombieType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.3, 0, -0.3, 0.3, 1.9, 0.3)
}

func (zombieType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (zombieType) Apply(data *world.EntityData) { zombieConf.Apply(data) }

func (zombieType) DecodeNBT(m map[string]any, data *world.EntityData) {
	zombieConf.Apply(data)
	decodeMob(m, data)
}

func (zombieType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
