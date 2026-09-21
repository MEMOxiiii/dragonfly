package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewDrowned creates a new Drowned at the position passed.
func NewDrowned(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(DrownedType, drownedConf)
}

var drownedConf = MobConfig{
	MaxHealth:       20,
	Speed:           0.06,
	Travel:          TravelSwim,
	ExperienceDrops: [2]int{5, 5},
	LavaDamage:      4,
	Breathable:      Breathable{Air: true, Water: true, TotalSupply: 15},
	BurnsInDaylight: true,
	Drops: LootTable{
		Always("minecraft:rotten_flesh", 0, 2, Looting(1)),
		OneOf(Entry("minecraft:trident", 1, 1, Weight(5)), Entry("minecraft:gold_ingot", 1, 1, Weight(5))),
	},
	Goals: []Goal{
		HurtByTarget(1),
		FleeSun(2, 1),
		NearestAttackableTarget(2, 12, 20, true),
		MeleeAttack(3, 3, 1),
		RandomStroll(6, 1),
		LookAtPlayer(7, 6, 1, 2),
		RandomLookAround(7),
	},
}

// DrownedType is a world.EntityType implementation for drowneds.
var DrownedType drownedType

type drownedType struct{}

func (drownedType) EncodeEntity() string { return "minecraft:drowned" }
func (drownedType) SpawnEggName() string { return "drowned" }
func (drownedType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.3, 0, -0.3, 0.3, 1.9, 0.3)
}

func (drownedType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (drownedType) Apply(data *world.EntityData) { drownedConf.Apply(data) }

func (drownedType) DecodeNBT(m map[string]any, data *world.EntityData) {
	drownedConf.Apply(data)
	decodeMob(m, data)
}

func (drownedType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
