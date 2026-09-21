package entity

import (
	"math"
	"math/rand/v2"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/go-gl/mathgl/mgl64"
)

// SquidSwimGoal moves a squid through the water in short bursts, the vanilla
// minecraft:behavior.squid_idle, squid_dive and squid_move_away_from_ground
// components.
type SquidSwimGoal struct {
	GoalBase
	SpeedMultiplier float64

	dir      mgl64.Vec3
	ticks    int
	cooldown int
}

func (g *SquidSwimGoal) CanStart(m *Mob) bool {
	if g.cooldown > 0 {
		g.cooldown--
		return false
	}
	if !m.InLiquid() {
		return false
	}
	g.dir = mgl64.Vec3{rand.Float64()*2 - 1, rand.Float64()*2 - 1, rand.Float64()*2 - 1}.Normalize()
	if g.nearGround(m) {
		g.dir[1] = math.Abs(g.dir[1])
	}
	g.ticks = 10 + rand.IntN(10)
	return true
}

func (g *SquidSwimGoal) CanContinue(m *Mob) bool {
	g.ticks--
	return g.ticks > 0 && m.InLiquid()
}

func (g *SquidSwimGoal) Stop(*Mob) { g.cooldown = 20 + rand.IntN(60) }

func (g *SquidSwimGoal) Tick(m *Mob) {
	m.MoveTowards(m.Position().Add(g.dir), g.speed())
}

func (g *SquidSwimGoal) nearGround(m *Mob) bool {
	below := cube.PosFromVec3(m.Position()).Side(cube.FaceDown)
	return len(m.tx.Block(below).Model().BBox(below, m.tx)) != 0
}

func (g *SquidSwimGoal) speed() float64 {
	if g.SpeedMultiplier == 0 {
		return 1
	}
	return g.SpeedMultiplier
}

func (SquidSwimGoal) Controls() Control { return ControlMove | ControlLook }

func (g *SquidSwimGoal) Clone() Goal { c := *g; return &c }
