package block

import (
	"math/rand/v2"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
)

// CopperBulb is a light source block that toggles its lit state whenever it receives a redstone pulse. It
// does not need continuous power to stay lit.
type CopperBulb struct {
	solid
	bassDrum

	// Oxidation is the level of oxidation of the copper bulb.
	Oxidation OxidationType
	// Waxed bool is whether the copper bulb has been waxed with honeycomb.
	Waxed bool
	// Lit is whether the copper bulb is currently emitting light.
	Lit bool
	// Powered is whether the copper bulb was receiving redstone power during its last redstone update.
	Powered bool
}

func (c CopperBulb) Strip() (world.Block, world.Sound, bool) {
	if c.Waxed {
		c.Waxed = false
		return c, sound.WaxRemoved{}, true
	} else if ot, ok := c.Oxidation.Decrease(); ok {
		c.Oxidation = ot
		return c, sound.CopperScraped{}, true
	}
	return c, nil, false
}

// Wax waxes the copper bulb to stop it from oxidising further.
func (c CopperBulb) Wax(cube.Pos, mgl64.Vec3) (world.Block, bool) {
	if c.Waxed {
		return c, false
	}
	c.Waxed = true
	return c, true
}

func (c CopperBulb) CanOxidate() bool {
	return !c.Waxed
}

func (c CopperBulb) OxidationLevel() OxidationType {
	return c.Oxidation
}

func (c CopperBulb) WithOxidationLevel(o OxidationType) Oxidisable {
	c.Oxidation = o
	return c
}

func (c CopperBulb) RandomTick(pos cube.Pos, tx *world.Tx, r *rand.Rand) {
	attemptOxidation(pos, tx, r, c)
}

// LightEmissionLevel returns the light emitted by the copper bulb while lit. Brighter oxidation stages emit
// less light, matching the copper bulb's darkening appearance as it oxidises.
func (c CopperBulb) LightEmissionLevel() uint8 {
	if !c.Lit {
		return 0
	}
	switch c.Oxidation {
	case ExposedOxidation():
		return 12
	case WeatheredOxidation():
		return 8
	case OxidisedOxidation():
		return 4
	default:
		return 15
	}
}

// RedstoneUpdate toggles the lit state of the copper bulb whenever it detects a rising edge of redstone
// power, i.e. when it starts being powered having not been powered on the previous update. Losing power
// does not toggle the bulb again.
func (c CopperBulb) RedstoneUpdate(pos cube.Pos, tx *world.Tx) {
	powered := c.powered(pos, tx)
	if powered == c.Powered {
		return
	}
	c.Powered = powered
	if powered {
		c.Lit = !c.Lit
		if c.Lit {
			tx.PlaySound(pos.Vec3Centre(), sound.CopperBulbTurnOn{Block: c})
		} else {
			tx.PlaySound(pos.Vec3Centre(), sound.CopperBulbTurnOff{Block: c})
		}
	}
	tx.SetBlock(pos, c, nil)
}

// powered reports whether the copper bulb is currently receiving redstone power.
func (c CopperBulb) powered(pos cube.Pos, tx *world.Tx) bool {
	for _, face := range cube.Faces() {
		if tx.RedstonePower(pos.Side(face), face, true) > 0 {
			return true
		}
	}
	return false
}

// BreakInfo ...
func (c CopperBulb) BreakInfo() BreakInfo {
	return newBreakInfo(3, func(t item.Tool) bool {
		return t.ToolType() == item.TypePickaxe && t.HarvestLevel() >= item.ToolTierStone.HarvestLevel
	}, pickaxeEffective, oneOf(c)).withBlastResistance(6)
}

// EncodeItem ...
func (c CopperBulb) EncodeItem() (name string, meta int16) {
	return copperBlockName("copper_bulb", c.Oxidation, c.Waxed), 0
}

// EncodeBlock ...
func (c CopperBulb) EncodeBlock() (string, map[string]any) {
	return copperBlockName("copper_bulb", c.Oxidation, c.Waxed), map[string]any{"lit": c.Lit, "powered_bit": c.Powered}
}

// allCopperBulbs returns a list of all copper bulb variants.
func allCopperBulbs() (c []world.Block) {
	f := func(waxed bool) {
		for _, o := range OxidationTypes() {
			c = append(c, CopperBulb{Oxidation: o, Waxed: waxed})
			c = append(c, CopperBulb{Oxidation: o, Waxed: waxed, Lit: true})
			c = append(c, CopperBulb{Oxidation: o, Waxed: waxed, Powered: true})
			c = append(c, CopperBulb{Oxidation: o, Waxed: waxed, Lit: true, Powered: true})
		}
	}
	f(true)
	f(false)
	return
}
