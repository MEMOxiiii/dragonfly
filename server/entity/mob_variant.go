package entity

import (
	"math/rand/v2"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
)

// VariantChoice is one of the variants a mob may spawn as, with the weight
// vanilla picks it with.
type VariantChoice struct {
	Value  int32
	Weight int
}

// Colour returns the colour of a mob such as a sheep.
func (m *Mob) Colour() item.Colour {
	return item.Colours()[m.mob().colour&15]
}

// SetColour sets the colour of a mob such as a sheep.
func (m *Mob) SetColour(c item.Colour) {
	m.mob().colour = c.Uint8()
	m.updateState()
}

// ColourValue returns the raw colour value of the mob, the vanilla
// minecraft:color component.
func (m *Mob) ColourValue() uint8 { return m.mob().colour }

// Variant returns the variant of the mob, the vanilla minecraft:variant
// component.
func (m *Mob) Variant() int32 { return m.mob().variant }

// SetVariant sets the variant of the mob.
func (m *Mob) SetVariant(variant int32) {
	m.mob().variant = variant
	m.updateState()
}

// Sheared reports whether a sheep has been sheared.
func (m *Mob) Sheared() bool { return m.mob().sheared }

// SetSheared shears the mob or grows its wool back.
func (m *Mob) SetSheared(sheared bool) {
	m.mob().sheared = sheared
	m.updateState()
}

// Saddled reports whether the mob carries a saddle.
func (m *Mob) Saddled() bool { return m.mob().saddled }

// SetSaddled puts a saddle on the mob or takes it off.
func (m *Mob) SetSaddled(saddled bool) {
	m.mob().saddled = saddled
	m.updateState()
}

func pickVariant(choices []VariantChoice) int32 {
	total := 0
	for _, c := range choices {
		total += c.Weight
	}
	if total == 0 {
		return 0
	}
	n := rand.IntN(total)
	for _, c := range choices {
		if n -= c.Weight; n < 0 {
			return c.Value
		}
	}
	return 0
}

// shear shears the mob if the item held is a pair of shears, dropping the
// items of its shear loot table.
func (b *mobBehaviour) shear(m *Mob, held item.Stack) bool {
	if len(b.conf.ShearDrops) == 0 || b.sheared {
		return false
	}
	if _, ok := held.Item().(item.Shears); !ok {
		return false
	}
	m.SetSheared(true)
	for _, drop := range b.conf.ShearDrops.roll(false, 0) {
		if name, _ := drop.Item().EncodeItem(); containsName(woolItems, name) {
			drop = colouredWool(m.Colour(), drop.Count())
		}
		m.tx.AddEntity(NewItem(world.EntitySpawnOpts{Position: m.Position().Add(dropOffset), Velocity: dropVelocity()}, drop))
	}
	return true
}

// woolItems are the names the vanilla loot tables use for wool, which is
// dropped in the colour of the mob.
var woolItems = []string{"minecraft:wool", "minecraft:white_wool"}

func colouredWool(c item.Colour, count int) item.Stack {
	return item.NewStack(block.Wool{Colour: c}, count)
}
