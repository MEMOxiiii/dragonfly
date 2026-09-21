package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewZombieNautilus creates a new ZombieNautilus at the position passed.
func NewZombieNautilus(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(ZombieNautilusType, zombieNautilusConf)
}

var zombieNautilusConf = MobConfig{
	MaxHealth:           15,
	Travel:              TravelSwim,
	ExperienceDrops:     [2]int{1, 3},
	KnockBackResistance: 0.3,
	LavaDamage:          4,
	Breathable:          Breathable{Air: true, Water: true, TotalSupply: 15},
	BurnsInDaylight:     true,
	Goals: []Goal{
		HurtByTarget(1),
		RandomSwim(4, 16, 4, 0, 1.5),
	},
}

// ZombieNautilusType is a world.EntityType implementation for zombie nautiluss.
var ZombieNautilusType zombieNautilusType

type zombieNautilusType struct{}

func (zombieNautilusType) EncodeEntity() string { return "minecraft:zombie_nautilus" }
func (zombieNautilusType) SpawnEggName() string { return "zombie_nautilus" }
func (zombieNautilusType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.4375, 0, -0.4375, 0.4375, 0.95, 0.4375)
}

func (zombieNautilusType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (zombieNautilusType) Apply(data *world.EntityData) { zombieNautilusConf.Apply(data) }

func (zombieNautilusType) DecodeNBT(m map[string]any, data *world.EntityData) {
	zombieNautilusConf.Apply(data)
	decodeMob(m, data)
}

func (zombieNautilusType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
