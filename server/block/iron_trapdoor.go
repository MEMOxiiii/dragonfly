package block

import (
	"math"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
)

// IronTrapdoor is a trapdoor variant that, unlike wooden trapdoors, cannot be opened or closed by hand: it
// only responds to redstone power.
type IronTrapdoor struct {
	transparent
	bass
	sourceWaterDisplacer

	// Facing is the direction the trapdoor is facing.
	Facing cube.Direction
	// Open is whether the trapdoor is open.
	Open bool
	// Top is whether the trapdoor occupies the top or bottom part of a block.
	Top bool
}

// Model ...
func (t IronTrapdoor) Model() world.BlockModel {
	return model.Trapdoor{Facing: t.Facing, Top: t.Top, Open: t.Open}
}

// UseOnBlock handles the directional placing of iron trapdoors and makes sure they are properly placed
// upside down when needed.
func (t IronTrapdoor) UseOnBlock(pos cube.Pos, face cube.Face, clickPos mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	pos, face, used := firstReplaceable(tx, pos, face, t)
	if !used {
		return false
	}
	t.Facing = user.Rotation().Direction().Opposite()
	t.Top = (clickPos.Y() > 0.5 && face != cube.FaceUp) || face == cube.FaceDown

	place(tx, pos, t, user, ctx)
	return placed(ctx)
}

// RedstoneUpdate opens or closes the trapdoor to match the redstone power currently received.
func (t IronTrapdoor) RedstoneUpdate(pos cube.Pos, tx *world.Tx) {
	powered := t.powered(pos, tx)
	if powered == t.Open {
		return
	}
	t.Open = powered
	tx.SetBlock(pos, t, nil)
	if t.Open {
		tx.PlaySound(pos.Vec3Centre(), sound.TrapdoorOpen{Block: t})
		return
	}
	tx.PlaySound(pos.Vec3Centre(), sound.TrapdoorClose{Block: t})
}

// powered reports whether the trapdoor is currently receiving redstone power.
func (t IronTrapdoor) powered(pos cube.Pos, tx *world.Tx) bool {
	for _, face := range cube.Faces() {
		if tx.RedstonePower(pos.Side(face), face, true) > 0 {
			return true
		}
	}
	return false
}

// BreakInfo ...
func (t IronTrapdoor) BreakInfo() BreakInfo {
	return newBreakInfo(5, pickaxeHarvestable, pickaxeEffective, oneOf(t)).withBlastResistance(5)
}

// SideClosed ...
func (t IronTrapdoor) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool {
	return false
}

// EncodeItem ...
func (t IronTrapdoor) EncodeItem() (name string, meta int16) {
	return "minecraft:iron_trapdoor", 0
}

// EncodeBlock ...
func (t IronTrapdoor) EncodeBlock() (name string, properties map[string]any) {
	return "minecraft:iron_trapdoor", map[string]any{"direction": int32(math.Abs(float64(t.Facing) - 3)), "open_bit": t.Open, "upside_down_bit": t.Top}
}

// allIronTrapdoors returns a list of all iron trapdoor block states.
func allIronTrapdoors() (trapdoors []world.Block) {
	for i := cube.Direction(0); i <= 3; i++ {
		trapdoors = append(trapdoors, IronTrapdoor{Facing: i, Open: false, Top: false})
		trapdoors = append(trapdoors, IronTrapdoor{Facing: i, Open: false, Top: true})
		trapdoors = append(trapdoors, IronTrapdoor{Facing: i, Open: true, Top: true})
		trapdoors = append(trapdoors, IronTrapdoor{Facing: i, Open: true, Top: false})
	}
	return
}
