package block

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
)

// IronDoor is a door variant that, unlike wooden doors, cannot be opened or closed by hand: it only
// responds to redstone power.
type IronDoor struct {
	transparent
	bass
	sourceWaterDisplacer

	// Facing is the direction that the door opens towards. When closed, the door sits on the side of its
	// block on the opposite direction.
	Facing cube.Direction
	// Open is whether the door is open.
	Open bool
	// Top is whether the block is the top or bottom half of a door.
	Top bool
	// Right is whether the door hinge is on the right side.
	Right bool
}

// Model ...
func (d IronDoor) Model() world.BlockModel {
	return model.Door{Facing: d.Facing, Open: d.Open, Right: d.Right}
}

// NeighbourUpdateTick ...
func (d IronDoor) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	if d.Top {
		if _, ok := tx.Block(pos.Side(cube.FaceDown)).(IronDoor); !ok {
			breakBlockNoDrops(d, pos, tx)
		}
	} else if solid := tx.Block(pos.Side(cube.FaceDown)).Model().FaceSolid(pos.Side(cube.FaceDown), cube.FaceUp, tx); !solid {
		breakBlock(d, pos, tx)
	} else if _, ok := tx.Block(pos.Side(cube.FaceUp)).(IronDoor); !ok {
		breakBlockNoDrops(d, pos, tx)
	}
}

// UseOnBlock handles the directional placing of iron doors.
func (d IronDoor) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	if face != cube.FaceUp {
		// Doors can only be placed when clicking the top face.
		return false
	}
	below := pos
	pos = pos.Side(cube.FaceUp)
	if !replaceableWith(tx, pos, d) || !replaceableWith(tx, pos.Side(cube.FaceUp), d) {
		return false
	}
	if !tx.Block(below).Model().FaceSolid(below, cube.FaceUp, tx) {
		return false
	}
	d.Facing = user.Rotation().Direction()
	left := tx.Block(pos.Side(d.Facing.RotateLeft().Face()))
	right := tx.Block(pos.Side(d.Facing.RotateRight().Face()))
	if _, ok := left.Model().(model.Door); ok {
		d.Right = true
	}
	// The side the door hinge is on can be affected by the blocks to the left and right of the door. In particular,
	// opaque blocks on the right side of the door with transparent blocks on the left side result in a right sided
	// door hinge.
	if diffuser, ok := right.(LightDiffuser); !ok || diffuser.LightDiffusionLevel() != 0 {
		if diffuser, ok := left.(LightDiffuser); ok && diffuser.LightDiffusionLevel() == 0 {
			d.Right = true
		}
	}

	ctx.IgnoreBBox = true
	place(tx, pos, d, user, ctx)
	place(tx, pos.Side(cube.FaceUp), IronDoor{Facing: d.Facing, Top: true, Right: d.Right}, user, ctx)
	ctx.CountSub = 1
	return placed(ctx)
}

// RedstoneUpdate opens or closes the door to match the redstone power currently received by either half.
func (d IronDoor) RedstoneUpdate(pos cube.Pos, tx *world.Tx) {
	powered := d.powered(pos, tx)
	if powered == d.Open {
		return
	}
	d.Open = powered
	tx.SetBlock(pos, d, nil)

	otherPos := pos.Side(cube.Face(boolByte(!d.Top)))
	if other, ok := tx.Block(otherPos).(IronDoor); ok && other.Open != d.Open {
		other.Open = d.Open
		tx.SetBlock(otherPos, other, nil)
	}
	if d.Open {
		tx.PlaySound(pos.Vec3Centre(), sound.DoorOpen{Block: d})
		return
	}
	tx.PlaySound(pos.Vec3Centre(), sound.DoorClose{Block: d})
}

// powered reports whether either half of the door is currently receiving redstone power.
func (d IronDoor) powered(pos cube.Pos, tx *world.Tx) bool {
	otherPos := pos.Side(cube.Face(boolByte(!d.Top)))
	for _, face := range cube.Faces() {
		if tx.RedstonePower(pos.Side(face), face, true) > 0 || tx.RedstonePower(otherPos.Side(face), face, true) > 0 {
			return true
		}
	}
	return false
}

// BreakInfo ...
func (d IronDoor) BreakInfo() BreakInfo {
	return newBreakInfo(5, pickaxeHarvestable, pickaxeEffective, oneOf(d)).withBlastResistance(5)
}

// SideClosed ...
func (d IronDoor) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool {
	return false
}

// EncodeItem ...
func (d IronDoor) EncodeItem() (name string, meta int16) {
	return "minecraft:iron_door", 0
}

// EncodeBlock ...
func (d IronDoor) EncodeBlock() (name string, properties map[string]any) {
	return "minecraft:iron_door", map[string]any{"minecraft:cardinal_direction": d.Facing.RotateRight().String(), "door_hinge_bit": d.Right, "open_bit": d.Open, "upper_block_bit": d.Top}
}

// allIronDoors returns a list of all iron door block states.
func allIronDoors() (doors []world.Block) {
	for i := cube.Direction(0); i <= 3; i++ {
		doors = append(doors, IronDoor{Facing: i, Open: false, Top: false, Right: false})
		doors = append(doors, IronDoor{Facing: i, Open: false, Top: true, Right: false})
		doors = append(doors, IronDoor{Facing: i, Open: true, Top: true, Right: false})
		doors = append(doors, IronDoor{Facing: i, Open: true, Top: false, Right: false})
		doors = append(doors, IronDoor{Facing: i, Open: false, Top: false, Right: true})
		doors = append(doors, IronDoor{Facing: i, Open: false, Top: true, Right: true})
		doors = append(doors, IronDoor{Facing: i, Open: true, Top: true, Right: true})
		doors = append(doors, IronDoor{Facing: i, Open: true, Top: false, Right: true})
	}
	return
}
