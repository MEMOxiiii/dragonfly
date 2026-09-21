package entity

import (
	"math/rand/v2"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/go-gl/mathgl/mgl64"
)

// RandomSwimGoal swims the mob to a random position in the water around it, the
// vanilla minecraft:behavior.random_swim component.
type RandomSwimGoal struct {
	GoalBase
	SpeedMultiplier float64
	// Range and YRange are the distances a target position may be away, the
	// vanilla xz_dist and y_dist fields.
	Range, YRange int
	// Interval is the number of ticks between two swims.
	Interval int

	target   mgl64.Vec3
	ticks    int
	cooldown int
}

func (g *RandomSwimGoal) CanStart(m *Mob) bool {
	if g.cooldown > 0 {
		g.cooldown--
		return false
	}
	if !m.InLiquid() {
		return false
	}
	pos, ok := g.destination(m)
	if !ok {
		g.cooldown = max(g.Interval, 10)
		return false
	}
	g.target, g.ticks = pos, 100
	return true
}

func (g *RandomSwimGoal) CanContinue(m *Mob) bool {
	g.ticks--
	return g.ticks > 0 && m.InLiquid() && m.Position().Sub(g.target).Len() > 1
}

func (g *RandomSwimGoal) Stop(*Mob) { g.cooldown = g.Interval }

func (g *RandomSwimGoal) Tick(m *Mob) {
	m.MoveTowards(g.target, g.speed())
}

func (g *RandomSwimGoal) destination(m *Mob) (mgl64.Vec3, bool) {
	r, y := float64(g.Range), float64(g.YRange)
	if r <= 0 {
		r = 16
	}
	if y <= 0 {
		y = 4
	}
	pos := m.Position()
	for i := 0; i < 12; i++ {
		if i == 4 {
			r, y = r/2, y/2
		} else if i == 8 {
			r, y = r/2, y/2
		}
		try := pos.Add(mgl64.Vec3{rand.Float64()*2*r - r, rand.Float64()*2*y - y, rand.Float64()*2*r - r})
		if _, ok := m.tx.Liquid(cube.PosFromVec3(try)); ok {
			return try, true
		}
	}
	return mgl64.Vec3{}, false
}

func (g *RandomSwimGoal) speed() float64 {
	if g.SpeedMultiplier == 0 {
		return 1
	}
	return g.SpeedMultiplier
}

// RandomFlyGoal flies the mob to a random position in the air around it, the
// vanilla minecraft:behavior.random_hover, float_wander and random_fly
// components.
type RandomFlyGoal struct {
	GoalBase
	SpeedMultiplier float64
	// Range and YRange are the distances a target position may be away, the
	// vanilla xz_dist and y_dist fields.
	Range, YRange int
	// YOffset is added to the height a target position is picked at.
	YOffset float64
	// MinHeight and MaxHeight bound how far above the ground the mob flies.
	MinHeight, MaxHeight float64

	target mgl64.Vec3
	ticks  int
}

func (g *RandomFlyGoal) CanStart(m *Mob) bool {
	pos, ok := g.destination(m)
	if !ok {
		return false
	}
	g.target, g.ticks = pos, 120
	return true
}

func (g *RandomFlyGoal) CanContinue(m *Mob) bool {
	g.ticks--
	return g.ticks > 0 && m.Position().Sub(g.target).Len() > 1
}

func (g *RandomFlyGoal) Tick(m *Mob) {
	m.MoveTowards(g.target, g.speed())
}

func (g *RandomFlyGoal) speed() float64 {
	if g.SpeedMultiplier == 0 {
		return 1
	}
	return g.SpeedMultiplier
}

func (g *RandomFlyGoal) destination(m *Mob) (mgl64.Vec3, bool) {
	r, y := float64(g.Range), float64(g.YRange)
	if r <= 0 {
		r = 10
	}
	if y <= 0 {
		y = 7
	}
	pos := m.Position()
	height := g.heightAbove(m, pos)
	for i := 0; i < 8; i++ {
		dy := rand.Float64()*2*y - y + g.YOffset
		if height > g.maxHeight() {
			dy = -rand.Float64() * (height - g.maxHeight() + 1)
		} else if height < g.minHeight() {
			dy = rand.Float64() * (g.minHeight() - height + 1)
		}
		try := pos.Add(mgl64.Vec3{rand.Float64()*2*r - r, dy, rand.Float64()*2*r - r})
		if !g.clear(m, try) {
			continue
		}
		if h := g.heightAbove(m, try); h < g.minHeight() || h > g.maxHeight()+2 {
			continue
		}
		return try, true
	}
	return mgl64.Vec3{}, false
}

func (g *RandomFlyGoal) clear(m *Mob, pos mgl64.Vec3) bool {
	at := cube.PosFromVec3(pos)
	return len(m.tx.Block(at).Model().BBox(at, m.tx)) == 0
}

func (g *RandomFlyGoal) heightAbove(m *Mob, pos mgl64.Vec3) float64 {
	at := cube.PosFromVec3(pos)
	for i := 0; i < 16; i++ {
		below := at.Sub(cube.Pos{0, i, 0})
		if len(m.tx.Block(below).Model().BBox(below, m.tx)) != 0 {
			return float64(i)
		}
	}
	return 16
}

func (g *RandomFlyGoal) minHeight() float64 {
	if g.MinHeight == 0 {
		return 1
	}
	return g.MinHeight
}

func (g *RandomFlyGoal) maxHeight() float64 {
	if g.MaxHeight == 0 {
		return 8
	}
	return g.MaxHeight
}

// MoveToWaterGoal walks the mob to the closest water, the vanilla
// minecraft:behavior.move_to_water component.
type MoveToWaterGoal struct {
	GoalBase
	SpeedMultiplier float64
	// SearchRange is the horizontal distance water is searched in.
	SearchRange int

	scanIn int
}

func (g *MoveToWaterGoal) CanStart(m *Mob) bool {
	if m.InLiquid() || !scanTick(&g.scanIn) {
		return false
	}
	pos, ok := g.water(m)
	return ok && m.MoveTo(pos, g.speed())
}

func (g *MoveToWaterGoal) CanContinue(m *Mob) bool { return !m.InLiquid() && m.Navigating() }

func (g *MoveToWaterGoal) Stop(m *Mob) { m.StopNavigating() }

func (g *MoveToWaterGoal) water(m *Mob) (cube.Pos, bool) {
	r := g.SearchRange
	if r <= 0 {
		r = 16
	}
	pos := cube.PosFromVec3(m.Position())
	for _, offset := range randomOffsets(r) {
		try := pos.Add(offset)
		if _, ok := m.tx.Liquid(try); ok {
			return try, true
		}
	}
	return cube.Pos{}, false
}

func (g *MoveToWaterGoal) speed() float64 {
	if g.SpeedMultiplier == 0 {
		return 1
	}
	return g.SpeedMultiplier
}

func randomOffsets(r int) []cube.Pos {
	offsets := make([]cube.Pos, 0, 16)
	for i := 0; i < 16; i++ {
		offsets = append(offsets, cube.Pos{rand.IntN(2*r+1) - r, rand.IntN(5) - 2, rand.IntN(2*r+1) - r})
	}
	return offsets
}

func (RandomSwimGoal) Controls() Control { return ControlMove | ControlLook }

func (RandomFlyGoal) Controls() Control { return ControlMove | ControlLook }

func (MoveToWaterGoal) Controls() Control { return ControlMove | ControlLook }

func (g *RandomSwimGoal) Clone() Goal { c := *g; return &c }

func (g *RandomFlyGoal) Clone() Goal { c := *g; return &c }

func (g *MoveToWaterGoal) Clone() Goal { c := *g; return &c }
