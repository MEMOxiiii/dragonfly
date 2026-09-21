package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewZoglin creates a new Zoglin at the position passed.
func NewZoglin(opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(ZoglinType, zoglinConf)
}

var zoglinConf = MobConfig{
	MaxHealth:           40,
	Speed:               0.25,
	ExperienceDrops:     [2]int{5, 5},
	FireImmune:          true,
	KnockBackResistance: 0.6,
	Breathable:          Breathable{Air: true, Water: true, TotalSupply: 15},
	Goals: []Goal{
		HurtByTarget(1),
		NearestAttackableTarget(3, 16, 16, false),
		MeleeAttack(4, 3, 1.15),
		RandomStroll(7, 1),
		LookAtPlayer(8, 6, 2, 4),
		RandomLookAround(9),
	},
}

// ZoglinType is a world.EntityType implementation for zoglins.
var ZoglinType zoglinType

type zoglinType struct{}

func (zoglinType) EncodeEntity() string { return "minecraft:zoglin" }
func (zoglinType) SpawnEggName() string { return "zoglin" }
func (zoglinType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.7, 0, -0.7, 0.7, 1.4, 0.7)
}

func (zoglinType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (zoglinType) Apply(data *world.EntityData) { zoglinConf.Apply(data) }

func (zoglinType) DecodeNBT(m map[string]any, data *world.EntityData) {
	zoglinConf.Apply(data)
	decodeMob(m, data)
}

func (zoglinType) EncodeNBT(data *world.EntityData) map[string]any { return encodeMob(data) }
