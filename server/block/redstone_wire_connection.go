package block

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// RedstoneWireConnection is the way redstone wire is drawn towards one of its horizontal sides.
type RedstoneWireConnection struct {
	redstoneWireConnection
}

// RedstoneWireNone is a side that the wire is not drawn towards.
func RedstoneWireNone() RedstoneWireConnection {
	return RedstoneWireConnection{0}
}

// RedstoneWireSide is a side that the wire runs flat towards.
func RedstoneWireSide() RedstoneWireConnection {
	return RedstoneWireConnection{1}
}

// RedstoneWireUp is a side that the wire climbs up, towards wire on top of the neighbouring block.
func RedstoneWireUp() RedstoneWireConnection {
	return RedstoneWireConnection{2}
}

// RedstoneWireConnections returns all possible RedstoneWireConnections.
func RedstoneWireConnections() []RedstoneWireConnection {
	return []RedstoneWireConnection{RedstoneWireNone(), RedstoneWireSide(), RedstoneWireUp()}
}

type redstoneWireConnection uint8

// Uint8 returns the RedstoneWireConnection as a uint8.
func (c redstoneWireConnection) Uint8() uint8 {
	return uint8(c)
}

// String returns the RedstoneWireConnection as a string.
func (c redstoneWireConnection) String() string {
	switch c {
	case 0:
		return "none"
	case 1:
		return "side"
	case 2:
		return "up"
	}
	panic("unknown redstone wire connection")
}

// redstoneWireConnections calculates how redstone wire at pos is drawn towards each horizontal side. Wire connects
// to wire next to it, one block lower or one block higher, and to power sources. Wire with a single connection
// also runs towards the opposite side, and wire without any runs towards all four.
func redstoneWireConnections(pos cube.Pos, src world.BlockSource) (north, east, south, west RedstoneWireConnection) {
	sides := make(map[cube.Face]RedstoneWireConnection, 4)
	aboveSolid := redstoneWireSolid(src, pos.Side(cube.FaceUp))
	for _, face := range cube.HorizontalFaces() {
		side := pos.Side(face)
		switch src.Block(side).(type) {
		case RedstoneWire, world.RedstonePowerSource:
			sides[face] = RedstoneWireSide()
			continue
		}
		if _, ok := src.Block(side.Side(cube.FaceUp)).(RedstoneWire); ok && !aboveSolid {
			sides[face] = RedstoneWireUp()
			continue
		}
		if _, ok := src.Block(side.Side(cube.FaceDown)).(RedstoneWire); ok && !redstoneWireSolid(src, side) {
			sides[face] = RedstoneWireSide()
		}
	}
	switch len(sides) {
	case 0:
		return RedstoneWireSide(), RedstoneWireSide(), RedstoneWireSide(), RedstoneWireSide()
	case 1:
		for face := range sides {
			sides[face.Opposite()] = RedstoneWireSide()
		}
	}
	get := func(face cube.Face) RedstoneWireConnection {
		if c, ok := sides[face]; ok {
			return c
		}
		return RedstoneWireNone()
	}
	return get(cube.FaceNorth), get(cube.FaceEast), get(cube.FaceSouth), get(cube.FaceWest)
}

// redstoneWireSolid reports whether the block at pos is solid on its top and bottom, which stops wire from
// connecting up or down past it.
func redstoneWireSolid(src world.BlockSource, pos cube.Pos) bool {
	m := src.Block(pos).Model()
	return m.FaceSolid(pos, cube.FaceUp, src) && m.FaceSolid(pos, cube.FaceDown, src)
}
