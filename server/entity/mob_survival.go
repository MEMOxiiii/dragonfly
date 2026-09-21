package entity

import (
	"math/rand/v2"
	"time"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
)

// Breathable holds the parameters of the vanilla minecraft:breathable
// component.
type Breathable struct {
	// Air, Water and Lava report which of them the mob can breathe.
	Air, Water, Lava bool
	// TotalSupply is the number of seconds the mob can go without breathing.
	TotalSupply int
	// SuffocateTime is the number of seconds between two points of damage once
	// the air supply runs out. A value of 0 damages the mob every second.
	SuffocateTime int
}

// tickBreathing damages a mob that cannot breathe where it is, such as a land
// mob under water or a fish on land.
func (b *mobBehaviour) tickBreathing(m *Mob) {
	br := b.conf.Breathable
	if br.TotalSupply == 0 && !br.Water && !br.Lava {
		return
	}
	pos := cube.PosFromVec3(m.EyePosition())
	liquid, hasLiquid := m.tx.Liquid(pos)
	breathing := br.Air && !hasLiquid
	if hasLiquid {
		if _, lava := liquid.(block.Lava); lava {
			breathing = br.Lava
		} else {
			breathing = br.Water
		}
	}
	if breathing {
		b.airSupply = br.TotalSupply * 20
		return
	}
	if b.airSupply > 0 {
		b.airSupply--
		return
	}
	interval := max(br.SuffocateTime*20, 20)
	if b.breathTicks++; b.breathTicks >= interval {
		b.breathTicks = 0
		b.hurt(m, 2, DrowningDamageSource{})
	}
}

// tickBlockDamage hurts the mob standing in a block that damages it, the
// vanilla minecraft:hurt_on_condition component.
func (b *mobBehaviour) tickBlockDamage(m *Mob) {
	if b.conf.LavaDamage == 0 {
		return
	}
	pos := cube.PosFromVec3(m.Position())
	if liquid, ok := m.tx.Liquid(pos); ok {
		if _, lava := liquid.(block.Lava); lava {
			b.hurt(m, b.conf.LavaDamage, block.LavaDamageSource{})
			m.SetOnFire(15 * time.Second)
		}
	}
}

// Healable holds the items a mob is healed by, the vanilla minecraft:healable
// component.
type Healable struct {
	// Items maps the encoded name of an item to the health it heals.
	Items map[string]float64
	// TamedOnly only heals the mob if it has been tamed.
	TamedOnly bool
}

// Tameable holds the parameters of the vanilla minecraft:tameable component.
type Tameable struct {
	// Items are the encoded names of the items the mob is tamed with.
	Items []string
	// Probability is the chance of a single item taming the mob.
	Probability float64
}

// Tamed reports whether the mob has been tamed by a player.
func (m *Mob) Tamed() bool { return m.mob().owner != nil }

// Owner returns the player that tamed the mob, if it has been tamed.
func (m *Mob) Owner(tx *world.Tx) (world.Entity, bool) {
	if owner := m.mob().owner; owner != nil {
		return owner.Entity(tx)
	}
	return nil, false
}

// SetOwner tames the mob for the entity passed.
func (m *Mob) SetOwner(e world.Entity) {
	b := m.mob()
	if e == nil {
		b.owner = nil
	} else {
		b.owner = e.H()
	}
	m.updateState()
}

// Sitting reports whether a tamed mob is sitting down.
func (m *Mob) Sitting() bool { return m.mob().sitting }

// SetSitting makes a tamed mob sit down or stand up again.
func (m *Mob) SetSitting(sitting bool) {
	m.mob().sitting = sitting
	m.updateState()
}

func (b *mobBehaviour) heal(m *Mob, user world.Entity, held item.Stack) bool {
	healable := b.conf.Healable
	if len(healable.Items) == 0 || m.Health() >= m.MaxHealth() {
		return false
	}
	if healable.TamedOnly && !m.Tamed() {
		return false
	}
	name, _ := held.Item().EncodeItem()
	health, ok := healable.Items[name]
	if !ok {
		return false
	}
	m.Heal(health, healingSource{})
	return true
}

func (b *mobBehaviour) tame(m *Mob, user world.Entity, held item.Stack) bool {
	tameable := b.conf.Tameable
	if len(tameable.Items) == 0 || m.Tamed() {
		return false
	}
	name, _ := held.Item().EncodeItem()
	if !containsName(tameable.Items, name) {
		return false
	}
	if rand.Float64() < tameable.Probability {
		m.SetOwner(user)
		for _, v := range m.tx.Viewers(m.Position()) {
			v.ViewEntityAction(m, TameAction{})
		}
	} else {
		for _, v := range m.tx.Viewers(m.Position()) {
			v.ViewEntityAction(m, TameFailAction{})
		}
	}
	return true
}

type healingSource struct{}

func (healingSource) HealingSource() {}
