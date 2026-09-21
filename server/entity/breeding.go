package entity

import (
	"math/rand/v2"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
)

const (
	loveTicks          = 600
	breedCooldownTicks = 6000
	growUpTicks        = 24000
)

// Interactable is an entity that reacts to a player using an item on it.
type Interactable interface {
	// Interact is called when a player uses the item held on the entity. It
	// returns true if one of the items held should be used up.
	Interact(user world.Entity, held item.Stack) bool
}

// Baby reports whether the mob is a baby.
func (m *Mob) Baby() bool { return m.mob().baby }

// SetBaby turns the mob into a baby or an adult.
func (m *Mob) SetBaby(baby bool) {
	b := m.mob()
	b.baby, b.age = baby, 0
	m.updateState()
}

// Scale returns the size of the mob relative to an adult of its type.
func (m *Mob) Scale() float64 {
	if m.mob().baby {
		return 0.5
	}
	return 1
}

// InLove reports whether the mob is looking for a partner to breed with.
func (m *Mob) InLove() bool { return m.mob().love > 0 }

// Interact feeds the mob the item held, putting it in love mode if the item is
// one it breeds with.
func (m *Mob) Interact(user world.Entity, held item.Stack) bool {
	b := m.mob()
	if held.Empty() {
		return false
	}
	if b.heal(m, user, held) || b.tame(m, user, held) || b.shear(m, held) {
		return true
	}
	if len(b.conf.BreedItems) == 0 || !b.breedItem(held) {
		return false
	}
	if b.baby {
		b.age += growUpTicks / 10
		return true
	}
	if b.love > 0 || b.breedCooldown > 0 {
		return false
	}
	b.love, b.partner = loveTicks, nil
	for _, v := range m.tx.Viewers(m.Position()) {
		v.ViewEntityAction(m, LoveAction{})
	}
	return true
}

func (b *mobBehaviour) breedItem(held item.Stack) bool {
	name, _ := held.Item().EncodeItem()
	for _, breedItem := range b.conf.BreedItems {
		if breedItem == name {
			return true
		}
	}
	return false
}

func (b *mobBehaviour) tickBreeding(m *Mob) {
	if b.baby {
		if b.age++; b.age >= growUpTicks {
			b.baby, b.age = false, 0
			m.updateState()
		}
		return
	}
	if b.breedCooldown > 0 {
		b.breedCooldown--
	}
	if b.love > 0 {
		if b.love--; b.love%20 == 0 {
			for _, v := range m.tx.Viewers(m.Position()) {
				v.ViewEntityAction(m, LoveAction{})
			}
		}
	}
}

// BreedGoal walks two mobs in love towards each other and spawns a baby, the
// vanilla minecraft:behavior.breed component.
type BreedGoal struct {
	GoalBase
	SpeedMultiplier float64

	partner *world.EntityHandle
	scanIn  int
}

func (g *BreedGoal) CanStart(m *Mob) bool {
	if !m.InLove() || !scanTick(&g.scanIn) {
		return false
	}
	partner, ok := g.nearestPartner(m)
	if !ok {
		return false
	}
	g.partner = partner.H()
	return true
}

func (g *BreedGoal) CanContinue(m *Mob) bool {
	if !m.InLove() || g.partner == nil {
		return false
	}
	e, ok := g.partner.Entity(m.tx)
	if !ok {
		return false
	}
	partner, ok := e.(*Mob)
	return ok && partner.InLove() && e.Position().Sub(m.Position()).Len() < 16
}

func (g *BreedGoal) Stop(m *Mob) {
	g.partner = nil
	m.StopNavigating()
}

func (g *BreedGoal) Tick(m *Mob) {
	e, ok := g.partner.Entity(m.tx)
	if !ok {
		return
	}
	partner := e.(*Mob)
	if diff := e.Position().Sub(m.Position()); diff.Len() > 2 {
		if !m.Navigating() {
			m.MoveTo(cube.PosFromVec3(e.Position()), g.speed())
		}
		return
	}
	m.StopNavigating()
	g.breed(m, partner)
}

func (g *BreedGoal) breed(m, partner *Mob) {
	b, pb := m.mob(), partner.mob()
	b.love, pb.love = 0, 0
	b.breedCooldown, pb.breedCooldown = breedCooldownTicks, breedCooldownTicks

	baby := NewMob(m.H().Type().(MobType), world.EntitySpawnOpts{Position: m.Position()})
	m.tx.AddEntity(baby).(*Mob).SetBaby(true)
	for _, orb := range NewExperienceOrbs(m.Position(), 1+rand.IntN(7)) {
		m.tx.AddEntity(orb)
	}
}

func (g *BreedGoal) speed() float64 {
	if g.SpeedMultiplier == 0 {
		return 1
	}
	return g.SpeedMultiplier
}

func (g *BreedGoal) nearestPartner(m *Mob) (*Mob, bool) {
	pos, best := m.Position(), 64.0
	var closest *Mob
	for e := range m.tx.EntitiesWithin(cube.Box(-8, -8, -8, 8, 8, 8).Translate(pos)) {
		other, ok := e.(*Mob)
		if !ok || other.H() == m.H() || e.H().Type() != m.H().Type() || !other.InLove() {
			continue
		}
		if d := e.Position().Sub(pos).LenSqr(); d < best {
			best, closest = d, other
		}
	}
	return closest, closest != nil
}

// TemptGoal makes the mob follow a player holding an item it likes, the vanilla
// minecraft:behavior.tempt component.
type TemptGoal struct {
	GoalBase
	SpeedMultiplier float64
	// Items are the encoded names of the items the mob follows.
	Items []string
	// Distance is the range at which the mob notices the item.
	Distance float64

	target *world.EntityHandle
	scanIn int
}

func (g *TemptGoal) CanStart(m *Mob) bool {
	if !scanTick(&g.scanIn) {
		return false
	}
	e, ok := g.holder(m)
	if !ok {
		return false
	}
	g.target = e.H()
	return true
}

func (g *TemptGoal) CanContinue(m *Mob) bool {
	if g.target == nil {
		return false
	}
	e, ok := g.target.Entity(m.tx)
	return ok && g.holding(e) && e.Position().Sub(m.Position()).Len() <= g.distance()
}

func (g *TemptGoal) Stop(m *Mob) {
	g.target = nil
	m.StopNavigating()
}

func (g *TemptGoal) Tick(m *Mob) {
	e, ok := g.target.Entity(m.tx)
	if !ok {
		return
	}
	m.LookAt(EyePosition(e))
	if e.Position().Sub(m.Position()).Len() < 2.5 {
		m.StopNavigating()
		return
	}
	if !m.Navigating() {
		m.MoveTo(cube.PosFromVec3(e.Position()), g.speed())
	}
}

func (g *TemptGoal) speed() float64 {
	if g.SpeedMultiplier == 0 {
		return 1
	}
	return g.SpeedMultiplier
}

func (g *TemptGoal) distance() float64 {
	if g.Distance == 0 {
		return 10
	}
	return g.Distance
}

func (g *TemptGoal) holder(m *Mob) (world.Entity, bool) {
	d := g.distance()
	pos, best := m.Position(), d*d
	var closest world.Entity
	for e := range m.tx.EntitiesWithin(cube.Box(-d, -d, -d, d, d, d).Translate(pos)) {
		if !g.holding(e) {
			continue
		}
		if dist := e.Position().Sub(pos).LenSqr(); dist < best {
			best, closest = dist, e
		}
	}
	return closest, closest != nil
}

func (g *TemptGoal) holding(e world.Entity) bool {
	if e.H().Type().EncodeEntity() != "minecraft:player" {
		return false
	}
	holder, ok := e.(interface {
		HeldItems() (item.Stack, item.Stack)
	})
	if !ok {
		return false
	}
	held, _ := holder.HeldItems()
	if held.Empty() {
		return false
	}
	name, _ := held.Item().EncodeItem()
	for _, tempting := range g.Items {
		if tempting == name {
			return true
		}
	}
	return false
}

// FollowParentGoal makes a baby walk to the closest adult of its own type, the
// vanilla minecraft:behavior.follow_parent component.
type FollowParentGoal struct {
	GoalBase
	SpeedMultiplier float64

	scanIn int
}

func (g *FollowParentGoal) CanStart(m *Mob) bool {
	if !m.Baby() || !scanTick(&g.scanIn) {
		return false
	}
	_, ok := g.parent(m)
	return ok
}

func (g *FollowParentGoal) CanContinue(m *Mob) bool { return m.Baby() }

func (g *FollowParentGoal) Stop(m *Mob) { m.StopNavigating() }

func (g *FollowParentGoal) Tick(m *Mob) {
	parent, ok := g.parent(m)
	if !ok {
		return
	}
	if parent.Position().Sub(m.Position()).Len() < 3 {
		m.StopNavigating()
		return
	}
	if !m.Navigating() {
		m.MoveTo(cube.PosFromVec3(parent.Position()), g.speed())
	}
}

func (g *FollowParentGoal) speed() float64 {
	if g.SpeedMultiplier == 0 {
		return 1.1
	}
	return g.SpeedMultiplier
}

func (g *FollowParentGoal) parent(m *Mob) (*Mob, bool) {
	pos, best := m.Position(), 256.0
	var closest *Mob
	for e := range m.tx.EntitiesWithin(cube.Box(-16, -16, -16, 16, 16, 16).Translate(pos)) {
		other, ok := e.(*Mob)
		if !ok || other.H() == m.H() || e.H().Type() != m.H().Type() || other.Baby() {
			continue
		}
		if d := e.Position().Sub(pos).LenSqr(); d < best {
			best, closest = d, other
		}
	}
	return closest, closest != nil
}

// scanTick counts down the ticks until the next search for a nearby entity.
func scanTick(in *int) bool {
	if *in--; *in > 0 {
		return false
	}
	*in = 10
	return true
}

func (BreedGoal) Controls() Control { return ControlMove | ControlLook }

func (TemptGoal) Controls() Control { return ControlMove | ControlLook }

func (FollowParentGoal) Controls() Control { return ControlMove | ControlLook }

func (g *BreedGoal) Clone() Goal { c := *g; return &c }

func (g *TemptGoal) Clone() Goal { c := *g; return &c }

func (g *FollowParentGoal) Clone() Goal { c := *g; return &c }
