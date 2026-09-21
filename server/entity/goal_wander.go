package entity

import (
	"math"
	"math/rand/v2"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// RandomStrollGoal walks the mob to a random position nearby, the vanilla
// minecraft:behavior.random_stroll component.
type RandomStrollGoal struct {
	GoalBase
	SpeedMultiplier float64
	// Range is the horizontal distance a target position may be away.
	Range int
	// Interval is the average number of ticks between two strolls.
	Interval int

	cooldown int
}

func (g *RandomStrollGoal) CanStart(m *Mob) bool {
	if g.cooldown > 0 {
		g.cooldown--
		return false
	}
	if _, ok := m.Target(); ok {
		return false
	}
	pos, ok := g.target(m)
	if !ok || !m.MoveTo(pos, g.speed()) {
		g.cooldown = g.interval()
		return false
	}
	return true
}

func (g *RandomStrollGoal) CanContinue(m *Mob) bool { return m.Navigating() }

func (g *RandomStrollGoal) Stop(m *Mob) {
	m.StopNavigating()
	g.cooldown = g.interval()
}

func (g *RandomStrollGoal) interval() int {
	if g.Interval <= 0 {
		return 40 + rand.IntN(40)
	}
	return rand.IntN(g.Interval) + g.Interval/2
}

func (g *RandomStrollGoal) speed() float64 {
	if g.SpeedMultiplier == 0 {
		return 1
	}
	return g.SpeedMultiplier
}

func (g *RandomStrollGoal) target(m *Mob) (cube.Pos, bool) {
	r := g.Range
	if r <= 0 {
		r = 10
	}
	pos := cube.PosFromVec3(m.Position())
	for i := 0; i < 10; i++ {
		try := pos.Add(cube.Pos{rand.IntN(2*r+1) - r, rand.IntN(3) - 1, rand.IntN(2*r+1) - r})
		if walkable(m, try) {
			return try, true
		}
	}
	return cube.Pos{}, false
}

// LookAtPlayerGoal turns the mob towards the closest player, the vanilla
// minecraft:behavior.look_at_player component.
type LookAtPlayerGoal struct {
	GoalBase
	Distance float64
	// LookTime is the range in seconds the mob keeps looking at a player.
	MinLookTime, MaxLookTime float64

	target *world.EntityHandle
	ticks  int
}

func (g *LookAtPlayerGoal) CanStart(m *Mob) bool {
	if rand.Float64() > 0.02 {
		return false
	}
	e, ok := closestPlayer(m, g.Distance)
	if !ok {
		return false
	}
	g.target, g.ticks = e.H(), g.lookTicks()
	return true
}

func (g *LookAtPlayerGoal) CanContinue(m *Mob) bool {
	if g.ticks <= 0 || g.target == nil {
		return false
	}
	e, ok := g.target.Entity(m.tx)
	return ok && e.Position().Sub(m.Position()).Len() <= g.Distance
}

func (g *LookAtPlayerGoal) Tick(m *Mob) {
	g.ticks--
	if e, ok := g.target.Entity(m.tx); ok {
		m.LookAt(EyePosition(e))
	}
}

func (g *LookAtPlayerGoal) Stop(*Mob) { g.target = nil }

func (g *LookAtPlayerGoal) lookTicks() int {
	minTime, maxTime := g.MinLookTime, g.MaxLookTime
	if maxTime <= minTime {
		maxTime = minTime + 1
	}
	return int((minTime + rand.Float64()*(maxTime-minTime)) * 20)
}

// RandomLookAroundGoal turns the mob's head to a random direction, the vanilla
// minecraft:behavior.random_look_around component.
type RandomLookAroundGoal struct {
	GoalBase

	ticks int
	yaw   float64
}

func (g *RandomLookAroundGoal) CanStart(*Mob) bool {
	if rand.Float64() > 0.02 {
		return false
	}
	g.ticks = 20 + rand.IntN(20)
	g.yaw = rand.Float64()*360 - 180
	return true
}

func (g *RandomLookAroundGoal) CanContinue(*Mob) bool { return g.ticks > 0 }

func (g *RandomLookAroundGoal) Tick(m *Mob) {
	g.ticks--
	rot := m.Rotation()
	m.SetRotation(cube.Rotation{rot.Yaw() + math.Max(math.Min(g.yaw-rot.Yaw(), 10), -10), 0})
}

// FloatGoal keeps the mob swimming at the surface of a liquid, the vanilla
// minecraft:behavior.float component.
type FloatGoal struct {
	GoalBase
}

func (FloatGoal) CanStart(m *Mob) bool {
	pos := cube.PosFromVec3(m.Position())
	_, ok := m.tx.Liquid(pos)
	return ok
}

func (g FloatGoal) CanContinue(m *Mob) bool { return g.CanStart(m) }

func (FloatGoal) Tick(m *Mob) {
	if rand.Float64() < 0.8 {
		m.Jump()
	}
}

func closestPlayer(m *Mob, dist float64) (world.Entity, bool) {
	var closest world.Entity
	best := dist * dist
	pos := m.Position()
	for e := range m.tx.EntitiesWithin(cube.Box(-dist, -dist, -dist, dist, dist, dist).Translate(pos)) {
		if e.H().Type().EncodeEntity() != "minecraft:player" {
			continue
		}
		if d := e.Position().Sub(pos).LenSqr(); d < best {
			best, closest = d, e
		}
	}
	return closest, closest != nil
}

func (FloatGoal) Controls() Control { return ControlJump }

func (RandomStrollGoal) Controls() Control { return ControlMove | ControlLook }

func (LookAtPlayerGoal) Controls() Control { return ControlLook }

func (RandomLookAroundGoal) Controls() Control { return ControlLook }

func (g FloatGoal) Clone() Goal { return g }

func (g *RandomStrollGoal) Clone() Goal { c := *g; return &c }

func (g *LookAtPlayerGoal) Clone() Goal { c := *g; return &c }

func (g *RandomLookAroundGoal) Clone() Goal { c := *g; return &c }
