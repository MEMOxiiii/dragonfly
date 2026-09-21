package entity

import (
	"math/rand/v2"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NearestAttackableTargetGoal makes the mob target the closest entity it may
// attack, the vanilla minecraft:behavior.nearest_attackable_target component.
type NearestAttackableTargetGoal struct {
	GoalBase
	// WithinRadius is the distance at which a target is searched for.
	WithinRadius float64
	// MaxDistance is the distance at which a target is forgotten again.
	MaxDistance float64
	// ReselectTargets makes the mob look for a closer target while it already
	// has one.
	ReselectTargets bool
	// Targets are the encoded names of the entities the mob attacks. Players
	// are targeted if it is empty.
	Targets []string
	// ScanInterval is the number of ticks between two searches for a target.
	ScanInterval int

	scanIn int
}

func (g *NearestAttackableTargetGoal) CanStart(m *Mob) bool {
	if !g.scan() {
		return false
	}
	if _, ok := m.Target(); ok && !g.ReselectTargets {
		return false
	}
	e, ok := g.nearest(m)
	if !ok {
		return false
	}
	m.SetTarget(e)
	return true
}

func (g *NearestAttackableTargetGoal) CanContinue(m *Mob) bool {
	e, ok := m.Target()
	if !ok {
		return false
	}
	if l, living := e.(Living); living && l.Dead() {
		m.SetTarget(nil)
		return false
	}
	if e.Position().Sub(m.Position()).Len() > g.maxDistance() {
		m.SetTarget(nil)
		return false
	}
	return true
}

func (g *NearestAttackableTargetGoal) Tick(m *Mob) {
	if !g.ReselectTargets || !g.scan() {
		return
	}
	if e, ok := g.nearest(m); ok {
		m.SetTarget(e)
	}
}

// scan reports whether the goal searches for a target this tick.
func (g *NearestAttackableTargetGoal) scan() bool {
	if g.scanIn--; g.scanIn > 0 {
		return false
	}
	g.scanIn = g.ScanInterval
	if g.scanIn <= 0 {
		g.scanIn = 10
	}
	return true
}

func (g *NearestAttackableTargetGoal) maxDistance() float64 {
	if g.MaxDistance == 0 {
		return g.WithinRadius + 10
	}
	return g.MaxDistance
}

func (g *NearestAttackableTargetGoal) nearest(m *Mob) (world.Entity, bool) {
	r := g.WithinRadius
	pos, best := m.Position(), r*r
	var closest world.Entity
	for e := range m.tx.EntitiesWithin(cube.Box(-r, -r, -r, r, r, r).Translate(pos)) {
		if !g.attackable(e) {
			continue
		}
		if l, ok := e.(Living); !ok || l.Dead() {
			continue
		}
		if d := e.Position().Sub(pos).LenSqr(); d < best {
			best, closest = d, e
		}
	}
	return closest, closest != nil
}

func (g *NearestAttackableTargetGoal) attackable(e world.Entity) bool {
	if len(g.Targets) == 0 {
		return survivalPlayer(e)
	}
	name := e.H().Type().EncodeEntity()
	if name == "minecraft:player" {
		return containsName(g.Targets, name) && survivalPlayer(e)
	}
	return containsName(g.Targets, name)
}

// HurtByTargetGoal makes the mob attack back whatever hurt it, the vanilla
// minecraft:behavior.hurt_by_target component.
type HurtByTargetGoal struct {
	GoalBase
}

func (HurtByTargetGoal) CanStart(m *Mob) bool {
	b := m.mob()
	if b.lastAttacker == nil {
		return false
	}
	attacker, ok := b.lastAttacker.Entity(m.tx)
	b.lastAttacker = nil
	if !ok {
		return false
	}
	m.SetTarget(attacker)
	return false
}

func (HurtByTargetGoal) CanContinue(*Mob) bool { return false }

// MeleeAttackGoal walks the mob to its target and hits it, the vanilla
// minecraft:behavior.melee_attack component.
type MeleeAttackGoal struct {
	GoalBase
	SpeedMultiplier float64
	// Damage is the damage of a single hit, the vanilla minecraft:attack
	// component.
	Damage float64
	// Reach is added to the width of both entities to find the distance at
	// which the mob may hit its target.
	Reach float64

	cooldown int
	repathIn int
}

func (g *MeleeAttackGoal) CanStart(m *Mob) bool {
	_, ok := m.Target()
	return ok
}

func (g *MeleeAttackGoal) CanContinue(m *Mob) bool {
	e, ok := m.Target()
	if !ok {
		return false
	}
	l, living := e.(Living)
	return living && !l.Dead()
}

func (g *MeleeAttackGoal) Start(m *Mob) { g.cooldown, g.repathIn = 0, 0 }

func (g *MeleeAttackGoal) Stop(m *Mob) { m.StopNavigating() }

func (g *MeleeAttackGoal) Tick(m *Mob) {
	e, ok := m.Target()
	if !ok {
		return
	}
	m.LookAt(EyePosition(e))
	if g.cooldown > 0 {
		g.cooldown--
	}
	if g.repathIn--; g.repathIn <= 0 {
		g.repathIn = 4 + rand.IntN(7)
		m.MoveTo(cube.PosFromVec3(e.Position()), g.speed())
	}
	if g.cooldown == 0 && g.inReach(m, e) {
		g.cooldown = 20
		m.StopNavigating()
		if l, living := e.(Living); living {
			l.Hurt(g.Damage, AttackDamageSource{Attacker: m})
		}
		for _, v := range m.tx.Viewers(m.Position()) {
			v.ViewEntityAction(m, SwingArmAction{})
		}
	}
}

func (g *MeleeAttackGoal) inReach(m *Mob, e world.Entity) bool {
	reach := g.Reach
	if reach == 0 {
		reach = 1
	}
	width := (m.H().Type().BBox(m).Width() + e.H().Type().BBox(e).Width()) / 2
	diff := e.Position().Sub(m.Position())
	diff[1] = 0
	return diff.Len() <= width+reach
}

func (g *MeleeAttackGoal) speed() float64 {
	if g.SpeedMultiplier == 0 {
		return 1
	}
	return g.SpeedMultiplier
}

// PanicGoal makes the mob run away after it is hurt, the vanilla
// minecraft:behavior.panic component.
type PanicGoal struct {
	GoalBase
	SpeedMultiplier float64

	ticks int
}

func (g *PanicGoal) CanStart(m *Mob) bool {
	if !m.mob().panicking {
		return false
	}
	pos, ok := g.destination(m)
	if !ok {
		return false
	}
	g.ticks = 100
	return m.MoveTo(pos, g.speed())
}

func (g *PanicGoal) CanContinue(m *Mob) bool {
	g.ticks--
	return g.ticks > 0 && m.Navigating()
}

func (g *PanicGoal) Stop(m *Mob) {
	m.StopNavigating()
	m.mob().panicking = false
}

func (g *PanicGoal) speed() float64 {
	if g.SpeedMultiplier == 0 {
		return 1.25
	}
	return g.SpeedMultiplier
}

func (g *PanicGoal) destination(m *Mob) (cube.Pos, bool) {
	pos := cube.PosFromVec3(m.Position())
	for i := 0; i < 10; i++ {
		try := pos.Add(cube.Pos{rand.IntN(21) - 10, rand.IntN(5) - 2, rand.IntN(21) - 10})
		if walkable(m, try) {
			return try, true
		}
	}
	return cube.Pos{}, false
}

func survivalPlayer(e world.Entity) bool {
	if e.H().Type().EncodeEntity() != "minecraft:player" {
		return false
	}
	g, ok := e.(interface{ GameMode() world.GameMode })
	return !ok || g.GameMode().AllowsTakingDamage()
}

func (PanicGoal) Controls() Control { return ControlMove | ControlLook }

func (HurtByTargetGoal) Controls() Control { return ControlTarget }

func (NearestAttackableTargetGoal) Controls() Control { return ControlTarget }

func (MeleeAttackGoal) Controls() Control { return ControlMove | ControlLook }

func (g *PanicGoal) Clone() Goal { c := *g; return &c }

func (g HurtByTargetGoal) Clone() Goal { return g }

func (g *NearestAttackableTargetGoal) Clone() Goal { c := *g; return &c }

func (g *MeleeAttackGoal) Clone() Goal { c := *g; return &c }
