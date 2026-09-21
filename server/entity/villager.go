package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewVillager creates a new Villager at the position passed.
func NewVillager(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(VillagerType, villagerConf)
}

var villagerConf = MobConfig{
	MaxHealth:  20,
	Speed:      0.5,
	LavaDamage: 4,
	Breathable: Breathable{Air: true, TotalSupply: 15},
	Variants:   []VariantChoice{{0, 5}, {0, 5}, {0, 5}, {0, 5}, {1, 20}, {1, 20}, {2, 20}, {3, 6}, {3, 6}, {3, 6}, {4, 10}, {4, 10}},
	Goals: []Goal{
		Float(0),
		AvoidMobType(3, 8, 0.6, 0.6, "minecraft:illager", "minecraft:vex", "minecraft:zombie", "minecraft:zombie_pigman", "minecraft:zombie_villager"),
		Panic(3, 0.6),
		RandomStroll(11, 0.6),
		LookAtPlayer(12, 8, 1, 2),
	},
}

// VillagerType is a world.EntityType implementation for villagers.
var VillagerType villagerType

type villagerType struct{}

func (villagerType) EncodeEntity() string { return "minecraft:villager" }
func (villagerType) SpawnEggName() string { return "villager" }
func (villagerType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.3, 0, -0.3, 0.3, 1.9, 0.3)
}

func (villagerType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (villagerType) Apply(data *world.EntityData) { villagerConf.Apply(data) }

func (villagerType) DecodeNBT(m map[string]any, data *world.EntityData) {
	villagerConf.Apply(data)
	decodeMob(m, data)
}

func (villagerType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
