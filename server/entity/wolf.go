package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewWolf creates a new Wolf at the position passed.
func NewWolf(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(WolfType, wolfConf)
}

var wolfConf = MobConfig{
	MaxHealth:       20,
	Speed:           0.3,
	ExperienceDrops: [2]int{1, 3},
	LavaDamage:      4,
	Breathable:      Breathable{Air: true, TotalSupply: 15},
	Healable:        Healable{Items: map[string]float64{"minecraft:porkchop": 3, "minecraft:cooked_porkchop": 8, "minecraft:cod": 2, "minecraft:salmon": 2, "minecraft:tropical_fish": 1, "minecraft:pufferfish": 1, "minecraft:cooked_cod": 5, "minecraft:cooked_salmon": 6, "minecraft:beef": 3, "minecraft:cooked_beef": 8, "minecraft:chicken": 2, "minecraft:cooked_chicken": 6, "minecraft:mutton": 2, "minecraft:cooked_mutton": 6, "minecraft:rotten_flesh": 4, "minecraft:rabbit": 3, "minecraft:cooked_rabbit": 5, "minecraft:rabbit_stew": 10}},
	Tameable:        Tameable{Items: []string{"minecraft:bone"}, Probability: 0.33},
	Colours:         []VariantChoice{{14, 1}},
	BreedItems:      []string{"minecraft:chicken", "minecraft:cooked_chicken", "minecraft:beef", "minecraft:cooked_beef", "minecraft:mutton", "minecraft:cooked_mutton", "minecraft:porkchop", "minecraft:cooked_porkchop", "minecraft:rabbit", "minecraft:cooked_rabbit", "minecraft:rotten_flesh"},
	Goals: []Goal{
		Float(0),
		HurtByTarget(3),
		AvoidMobType(3, 24, 1.5, 1.5, "minecraft:llama"),
		MeleeAttack(5, 4, 1),
		NearestAttackableTarget(5, 16, 16, false),
		LookAtPlayer(6, 8, 1, 2),
		Breed(7, 1),
		RandomStroll(8, 1),
	},
}

// WolfType is a world.EntityType implementation for wolfs.
var WolfType wolfType

type wolfType struct{}

func (wolfType) EncodeEntity() string { return "minecraft:wolf" }
func (wolfType) SpawnEggName() string { return "wolf" }
func (wolfType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.3, 0, -0.3, 0.3, 0.8, 0.3)
}

func (wolfType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (wolfType) Apply(data *world.EntityData) { wolfConf.Apply(data) }

func (wolfType) DecodeNBT(m map[string]any, data *world.EntityData) {
	wolfConf.Apply(data)
	decodeMob(m, data)
}

func (wolfType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
