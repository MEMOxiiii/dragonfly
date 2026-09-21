package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewWanderingTrader creates a new WanderingTrader at the position passed.
func NewWanderingTrader(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(WanderingTraderType, wanderingTraderConf)
}

var wanderingTraderConf = MobConfig{
	MaxHealth:  20,
	Speed:      0.5,
	LavaDamage: 4,
	Breathable: Breathable{Air: true, TotalSupply: 15},
	Goals: []Goal{
		Float(0),
		Panic(1, 0.6),
		AvoidMobType(2, 6, 0.6, 0.6, "minecraft:illager", "minecraft:vex", "minecraft:zombie", "minecraft:zombie_pigman", "minecraft:zombie_villager"),
		LookAtPlayer(8, 8, 1, 2),
		RandomLookAround(9),
	},
}

// WanderingTraderType is a world.EntityType implementation for wandering traders.
var WanderingTraderType wanderingTraderType

type wanderingTraderType struct{}

func (wanderingTraderType) EncodeEntity() string { return "minecraft:wandering_trader" }
func (wanderingTraderType) SpawnEggName() string { return "wandering_trader" }
func (wanderingTraderType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.3, 0, -0.3, 0.3, 1.9, 0.3)
}

func (wanderingTraderType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (wanderingTraderType) Apply(data *world.EntityData) { wanderingTraderConf.Apply(data) }

func (wanderingTraderType) DecodeNBT(m map[string]any, data *world.EntityData) {
	wanderingTraderConf.Apply(data)
	decodeMob(m, data)
}

func (wanderingTraderType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
