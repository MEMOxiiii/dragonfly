package enchantment

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
)

// LuckOfTheSea is an enchantment that increases the chance of catching treasure while fishing, at the cost of
// the chance of catching junk.
var LuckOfTheSea luckOfTheSea

type luckOfTheSea struct{}

// Name ...
func (luckOfTheSea) Name() string {
	return "Luck of the Sea"
}

// MaxLevel ...
func (luckOfTheSea) MaxLevel() int {
	return 3
}

// Cost ...
func (luckOfTheSea) Cost(level int) (int, int) {
	minCost := 15 + (level-1)*9
	return minCost, minCost + 50
}

// FishingChances returns the percentage points added to the chance of catching treasure and subtracted from
// the chance of catching junk.
func (luckOfTheSea) FishingChances(level int) (treasure, junk float64) {
	return 2.1 * float64(level), 1.95 * float64(level)
}

// Rarity ...
func (luckOfTheSea) Rarity() item.EnchantmentRarity {
	return item.EnchantmentRarityRare
}

// CompatibleWithEnchantment ...
func (luckOfTheSea) CompatibleWithEnchantment(item.EnchantmentType) bool {
	return true
}

// CompatibleWithItem ...
func (luckOfTheSea) CompatibleWithItem(i world.Item) bool {
	_, ok := i.(item.FishingRod)
	return ok
}
