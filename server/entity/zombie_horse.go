package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewZombieHorse creates a new ZombieHorse at the position passed.
func NewZombieHorse(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(ZombieHorseType, zombieHorseConf)
}

var zombieHorseConf = MobConfig{
	MaxHealth:       15,
	Speed:           0.2,
	ExperienceDrops: [2]int{1, 3},
	LavaDamage:      4,
	Breathable:      Breathable{Air: true, Water: true, TotalSupply: 15},
	Drops: LootTable{
		Always("minecraft:rotten_flesh", 0, 2, Looting(1)),
	},
	Goals: []Goal{
		Panic(1, 1.2),
		RandomStroll(6, 0.7),
		LookAtPlayer(7, 6, 1, 2),
		RandomLookAround(8),
	},
}

// ZombieHorseType is a world.EntityType implementation for zombie horses.
var ZombieHorseType zombieHorseType

type zombieHorseType struct{}

func (zombieHorseType) EncodeEntity() string { return "minecraft:zombie_horse" }
func (zombieHorseType) SpawnEggName() string { return "zombie_horse" }
func (zombieHorseType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.7, 0, -0.7, 0.7, 1.6, 0.7)
}

func (zombieHorseType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (zombieHorseType) Apply(data *world.EntityData) { zombieHorseConf.Apply(data) }

func (zombieHorseType) DecodeNBT(m map[string]any, data *world.EntityData) {
	zombieHorseConf.Apply(data)
	decodeMob(m, data)
}

func (zombieHorseType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
