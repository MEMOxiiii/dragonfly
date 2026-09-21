package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewTurtle creates a new Turtle at the position passed.
func NewTurtle(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(TurtleType, turtleConf)
}

var turtleConf = MobConfig{
	MaxHealth:       30,
	Speed:           0.12,
	Travel:          TravelSwim,
	ExperienceDrops: [2]int{1, 4},
	LavaDamage:      4,
	Breathable:      Breathable{Air: true, Water: true, TotalSupply: 15},
	BreedItems:      []string{"minecraft:seagrass"},
	Drops: LootTable{
		Always("minecraft:seagrass", 0, 2, Looting(1)),
	},
	Goals: []Goal{
		Panic(0, 1.2),
		HurtByTarget(1),
		Breed(2, 1),
		Tempt(3, 1.1, 10, "minecraft:seagrass"),
		MoveToWater(4, 16, 1),
		RandomSwim(7, 30, 15, 0, 1),
		LookAtPlayer(8, 8, 1, 2),
		RandomStroll(9, 1),
	},
}

// TurtleType is a world.EntityType implementation for turtles.
var TurtleType turtleType

type turtleType struct{}

func (turtleType) EncodeEntity() string { return "minecraft:turtle" }
func (turtleType) SpawnEggName() string { return "turtle" }
func (turtleType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.6, 0, -0.6, 0.6, 0.4, 0.6)
}

func (turtleType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (turtleType) Apply(data *world.EntityData) { turtleConf.Apply(data) }

func (turtleType) DecodeNBT(m map[string]any, data *world.EntityData) {
	turtleConf.Apply(data)
	decodeMob(m, data)
}

func (turtleType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
