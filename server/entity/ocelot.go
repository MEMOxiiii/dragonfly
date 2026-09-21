package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewOcelot creates a new Ocelot at the position passed.
func NewOcelot(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(OcelotType, ocelotConf)
}

var ocelotConf = MobConfig{
	MaxHealth:       10,
	Speed:           0.3,
	ExperienceDrops: [2]int{1, 3},
	LavaDamage:      4,
	Breathable:      Breathable{Air: true, TotalSupply: 15},
	BreedItems:      []string{"minecraft:cod", "minecraft:salmon"},
	Goals: []Goal{
		Float(0),
		NearestAttackableTarget(1, 16, 8, true),
		Panic(1, 1.25),
		Breed(3, 1),
		Tempt(4, 0.5, 16, "minecraft:cod", "minecraft:salmon"),
		AvoidMobType(5, 10, 0.8, 1.33, "minecraft:player"),
		RandomStroll(8, 0.8),
		LookAtPlayer(9, 8, 1, 2),
	},
}

// OcelotType is a world.EntityType implementation for ocelots.
var OcelotType ocelotType

type ocelotType struct{}

func (ocelotType) EncodeEntity() string { return "minecraft:ocelot" }
func (ocelotType) SpawnEggName() string { return "ocelot" }
func (ocelotType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.3, 0, -0.3, 0.3, 0.7, 0.3)
}

func (ocelotType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (ocelotType) Apply(data *world.EntityData) { ocelotConf.Apply(data) }

func (ocelotType) DecodeNBT(m map[string]any, data *world.EntityData) {
	ocelotConf.Apply(data)
	decodeMob(m, data)
}

func (ocelotType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
