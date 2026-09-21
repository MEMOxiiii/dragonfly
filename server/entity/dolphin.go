package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewDolphin creates a new Dolphin at the position passed.
func NewDolphin(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(DolphinType, dolphinConf)
}

var dolphinConf = MobConfig{
	MaxHealth:       10,
	Speed:           0.15,
	Travel:          TravelSwim,
	ExperienceDrops: [2]int{1, 4},
	LavaDamage:      4,
	Breathable:      Breathable{Air: true, TotalSupply: 240},
	Drops: LootTable{
		Always("minecraft:cod", 0, 1, Smelted("minecraft:cooked_cod"), Looting(1)),
	},
	Goals: []Goal{
		HurtByTarget(1),
		MoveToWater(1, 15, 1),
		AvoidMobType(2, 8, 1, 1, "minecraft:guardian", "minecraft:guardian_elder"),
		MeleeAttack(2, 3, 1),
		RandomSwim(5, 20, 4, 0, 1),
		RandomLookAround(7),
		LookAtPlayer(8, 6, 1, 2),
	},
}

// DolphinType is a world.EntityType implementation for dolphins.
var DolphinType dolphinType

type dolphinType struct{}

func (dolphinType) EncodeEntity() string { return "minecraft:dolphin" }
func (dolphinType) SpawnEggName() string { return "dolphin" }
func (dolphinType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.45, 0, -0.45, 0.45, 0.6, 0.45)
}

func (dolphinType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (dolphinType) Apply(data *world.EntityData) { dolphinConf.Apply(data) }

func (dolphinType) DecodeNBT(m map[string]any, data *world.EntityData) {
	dolphinConf.Apply(data)
	decodeMob(m, data)
}

func (dolphinType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
