package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewZombiePigman creates a new ZombiePigman at the position passed.
func NewZombiePigman(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(ZombiePigmanType, zombiePigmanConf)
}

var zombiePigmanConf = MobConfig{
	MaxHealth:       20,
	Speed:           0.23,
	ExperienceDrops: [2]int{5, 5},
	FireImmune:      true,
	Breathable:      Breathable{Air: true, Water: true, TotalSupply: 15},
	BurnsInDaylight: true,
	Drops: LootTable{
		Always("minecraft:rotten_flesh", 0, 1, Looting(1)),
		Always("minecraft:gold_nugget", 0, 1, Looting(1)),
		Always("minecraft:gold_ingot", 1, 1),
	},
	Goals: []Goal{
		HurtByTarget(1),
		MeleeAttack(3, 5, 1.5),
		RandomStroll(7, 1),
		LookAtPlayer(8, 6, 1, 2),
		RandomLookAround(9),
	},
}

// ZombiePigmanType is a world.EntityType implementation for zombie pigmans.
var ZombiePigmanType zombiePigmanType

type zombiePigmanType struct{}

func (zombiePigmanType) EncodeEntity() string { return "minecraft:zombie_pigman" }
func (zombiePigmanType) SpawnEggName() string { return "zombie_pigman" }
func (zombiePigmanType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.3, 0, -0.3, 0.3, 1.9, 0.3)
}

func (zombiePigmanType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (zombiePigmanType) Apply(data *world.EntityData) { zombiePigmanConf.Apply(data) }

func (zombiePigmanType) DecodeNBT(m map[string]any, data *world.EntityData) {
	zombiePigmanConf.Apply(data)
	decodeMob(m, data)
}

func (zombiePigmanType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
