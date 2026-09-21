package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewCreeper creates a new Creeper at the position passed.
func NewCreeper(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(CreeperType, creeperConf)
}

var creeperConf = MobConfig{
	MaxHealth:       20,
	Speed:           0.2,
	ExperienceDrops: [2]int{5, 5},
	LavaDamage:      4,
	Breathable:      Breathable{Air: true, TotalSupply: 15},
	Drops: LootTable{
		Always("minecraft:gunpowder", 0, 2, Looting(1)),
		OneOf(Entry("minecraft:music_disc_13", 1, 1), Entry("minecraft:music_disc_cat", 1, 1), Entry("minecraft:music_disc_blocks", 1, 1), Entry("minecraft:music_disc_chirp", 1, 1), Entry("minecraft:music_disc_far", 1, 1), Entry("minecraft:music_disc_mall", 1, 1), Entry("minecraft:music_disc_mellohi", 1, 1), Entry("minecraft:music_disc_stal", 1, 1), Entry("minecraft:music_disc_strad", 1, 1), Entry("minecraft:music_disc_ward", 1, 1), Entry("minecraft:music_disc_11", 1, 1), Entry("minecraft:music_disc_wait", 1, 1)),
		OneOf(Entry("minecraft:music_disc_13", 1, 1), Entry("minecraft:music_disc_cat", 1, 1), Entry("minecraft:music_disc_blocks", 1, 1), Entry("minecraft:music_disc_chirp", 1, 1), Entry("minecraft:music_disc_far", 1, 1), Entry("minecraft:music_disc_mall", 1, 1), Entry("minecraft:music_disc_mellohi", 1, 1), Entry("minecraft:music_disc_stal", 1, 1), Entry("minecraft:music_disc_strad", 1, 1), Entry("minecraft:music_disc_ward", 1, 1), Entry("minecraft:music_disc_11", 1, 1), Entry("minecraft:music_disc_wait", 1, 1)),
	},
	Goals: []Goal{
		Float(0),
		NearestAttackableTarget(1, 16, 16, false),
		HurtByTarget(2),
		Swell(2, 2.5, 6, 30, 3),
		AvoidMobType(3, 6, 1, 1.2, "minecraft:cat", "minecraft:ocelot"),
		MeleeAttack(4, 3, 1.25),
		RandomStroll(5, 1),
		LookAtPlayer(6, 8, 1, 2),
		RandomLookAround(6),
	},
}

// CreeperType is a world.EntityType implementation for creepers.
var CreeperType creeperType

type creeperType struct{}

func (creeperType) EncodeEntity() string { return "minecraft:creeper" }
func (creeperType) SpawnEggName() string { return "creeper" }
func (creeperType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.3, 0, -0.3, 0.3, 1.8, 0.3)
}

func (creeperType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (creeperType) Apply(data *world.EntityData) { creeperConf.Apply(data) }

func (creeperType) DecodeNBT(m map[string]any, data *world.EntityData) {
	creeperConf.Apply(data)
	decodeMob(m, data)
}

func (creeperType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
