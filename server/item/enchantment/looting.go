package enchantment

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
)

// Looting is an enchantment that increases the amount of items dropped by mobs.
var Looting looting

type looting struct{}

// Name ...
func (looting) Name() string {
	return "Looting"
}

// MaxLevel ...
func (looting) MaxLevel() int {
	return 3
}

// Cost ...
func (looting) Cost(level int) (int, int) {
	minCost := 15 + (level-1)*9
	return minCost, minCost + 50
}

// Rarity ...
func (looting) Rarity() item.EnchantmentRarity {
	return item.EnchantmentRarityRare
}

// CompatibleWithEnchantment ...
func (looting) CompatibleWithEnchantment(item.EnchantmentType) bool {
	return true
}

// CompatibleWithItem ...
func (looting) CompatibleWithItem(i world.Item) bool {
	t, ok := i.(item.Tool)
	return ok && t.ToolType() == item.TypeSword
}
