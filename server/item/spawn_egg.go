package item

import (
	"math/rand/v2"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

// SpawnEggType is the type of mob a SpawnEgg spawns.
type SpawnEggType interface {
	world.EntityType
	world.EntityConfig
	// SpawnEggName returns the name of the spawn egg of the mob without the
	// '_spawn_egg' suffix, such as 'pig'.
	SpawnEggName() string
}

// SpawnEgg spawns a mob of the type it holds where it is used.
type SpawnEgg struct {
	Type SpawnEggType
}

// UseOnBlock ...
func (e SpawnEgg) UseOnBlock(pos cube.Pos, face cube.Face, clickPos mgl64.Vec3, tx *world.Tx, _ User, ctx *UseContext) bool {
	spawn := pos.Side(face).Vec3Middle()
	if face == cube.FaceUp {
		spawn = pos.Vec3Middle().Add(mgl64.Vec3{0, clickPos[1]})
	}
	opts := world.EntitySpawnOpts{Position: spawn, Rotation: cube.Rotation{rand.Float64()*360 - 180, 0}}
	tx.AddEntity(opts.New(e.Type, e.Type))
	ctx.SubtractFromCount(1)
	return true
}

// EncodeItem ...
func (e SpawnEgg) EncodeItem() (name string, meta int16) {
	return "minecraft:" + e.Type.SpawnEggName() + "_spawn_egg", 0
}
