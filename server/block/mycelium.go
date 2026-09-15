package block

import (
	"math/rand/v2"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// Mycelium is a variant of dirt that generates naturally in mushroom fields biomes.
type Mycelium struct {
	solid
}

// SoilFor ...
func (m Mycelium) SoilFor(block world.Block) bool {
	switch block.(type) {
	case ShortGrass, Fern, DoubleTallGrass, Flower, DoubleFlower, NetherSprouts, PinkPetals, SugarCane, DeadBush, BambooSapling, Bamboo:
		return true
	}
	return false
}

// RandomTick handles the ticking of mycelium, which may or may not result in the spreading of mycelium onto
// dirt.
func (m Mycelium) RandomTick(pos cube.Pos, tx *world.Tx, r *rand.Rand) {
	spreadDirt(m, pos, tx, r)
}

// Shovel ...
func (Mycelium) Shovel() (world.Block, bool) {
	return DirtPath{}, true
}

// BreakInfo ...
func (m Mycelium) BreakInfo() BreakInfo {
	return newBreakInfo(0.6, alwaysHarvestable, shovelEffective, silkTouchOneOf(Dirt{}, m))
}

// EncodeItem ...
func (Mycelium) EncodeItem() (name string, meta int16) {
	return "minecraft:mycelium", 0
}

// EncodeBlock ...
func (Mycelium) EncodeBlock() (string, map[string]any) {
	return "minecraft:mycelium", nil
}
