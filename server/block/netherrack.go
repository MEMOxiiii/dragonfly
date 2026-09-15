package block

import (
	"math/rand/v2"
	"slices"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
)

// Netherrack is a block found in The Nether.
type Netherrack struct {
	solid
	bassDrum
}

// SoilFor ...
func (n Netherrack) SoilFor(block world.Block) bool {
	flower, ok := block.(Flower)
	return ok && flower.Type == WitherRose()
}

// BoneMeal converts the netherrack into nylium of a type found on one of its neighbouring blocks, picking one
// of the two types at random if both are present.
func (n Netherrack) BoneMeal(pos cube.Pos, tx *world.Tx) item.BoneMealResult {
	var types []Nylium
	for _, f := range cube.Faces() {
		if nylium, ok := tx.Block(pos.Side(f)).(Nylium); ok && !slices.Contains(types, nylium) {
			types = append(types, nylium)
		}
	}
	if len(types) == 0 {
		return item.BoneMealResultNone
	}
	tx.SetBlock(pos, types[rand.IntN(len(types))], nil)
	return item.BoneMealResultSmall
}

// BreakInfo ...
func (n Netherrack) BreakInfo() BreakInfo {
	return newBreakInfo(0.4, pickaxeHarvestable, pickaxeEffective, oneOf(n))
}

// SmeltInfo ...
func (Netherrack) SmeltInfo() item.SmeltInfo {
	return newSmeltInfo(item.NewStack(item.NetherBrick{}, 1), 0.1)
}

// EncodeItem ...
func (Netherrack) EncodeItem() (name string, meta int16) {
	return "minecraft:netherrack", 0
}

// EncodeBlock ...
func (Netherrack) EncodeBlock() (string, map[string]any) {
	return "minecraft:netherrack", nil
}
