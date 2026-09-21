package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewHoglin creates a new Hoglin at the position passed.
func NewHoglin(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(HoglinType, hoglinConf)
}

var hoglinConf = MobConfig{
	MaxHealth:           40,
	Speed:               0.3,
	ExperienceDrops:     [2]int{5, 5},
	KnockBackResistance: 0.6,
	LavaDamage:          4,
	Breathable:          Breathable{Air: true, TotalSupply: 15},
	BreedItems:          []string{"minecraft:crimson_fungus"},
	Goals: []Goal{
		AvoidMobType(0, 10, 1, 1.2, "minecraft:piglin"),
		HurtByTarget(2),
		Breed(3, 0.6),
		NearestAttackableTarget(4, 16, 16, false),
		MeleeAttack(4, 3, 1),
		RandomStroll(7, 0.4),
		LookAtPlayer(8, 6, 2, 4),
		RandomLookAround(9),
	},
}

// HoglinType is a world.EntityType implementation for hoglins.
var HoglinType hoglinType

type hoglinType struct{}

func (hoglinType) EncodeEntity() string { return "minecraft:hoglin" }
func (hoglinType) SpawnEggName() string { return "hoglin" }
func (hoglinType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.7, 0, -0.7, 0.7, 1.4, 0.7)
}

func (hoglinType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (hoglinType) Apply(data *world.EntityData) { hoglinConf.Apply(data) }

func (hoglinType) DecodeNBT(m map[string]any, data *world.EntityData) {
	hoglinConf.Apply(data)
	decodeMob(m, data)
}

func (hoglinType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
