package entity

import (
	"math"
	"math/rand/v2"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

// RangedAttackGoal shoots projectiles at the target of the mob, the vanilla
// minecraft:behavior.ranged_attack component.
type RangedAttackGoal struct {
	GoalBase
	SpeedMultiplier float64
	// AttackRadius is the distance at which the mob starts shooting.
	AttackRadius float64
	// MinInterval and MaxInterval are the bounds of the delay between two
	// shots, in seconds.
	MinInterval, MaxInterval float64
	// BurstShots is the number of projectiles fired in one attack, with
	// BurstInterval seconds between them.
	BurstShots    int
	BurstInterval float64
	// Shoot fires a single projectile at the position passed.
	Shoot func(m *Mob, target world.Entity)

	cooldown  int
	burstLeft int
	burstIn   int
	repathIn  int
}

func (g *RangedAttackGoal) CanStart(m *Mob) bool {
	e, ok := m.Target()
	return ok && g.Shoot != nil && e.Position().Sub(m.Position()).Len() <= g.AttackRadius*2
}

func (g *RangedAttackGoal) CanContinue(m *Mob) bool {
	e, ok := m.Target()
	if !ok {
		return false
	}
	l, living := e.(Living)
	return living && !l.Dead() && e.Position().Sub(m.Position()).Len() <= g.AttackRadius*2
}

func (g *RangedAttackGoal) Start(*Mob) { g.cooldown, g.burstLeft, g.repathIn = 0, 0, 0 }

func (g *RangedAttackGoal) Stop(m *Mob) { m.StopNavigating() }

func (g *RangedAttackGoal) Tick(m *Mob) {
	e, ok := m.Target()
	if !ok {
		return
	}
	m.LookAt(EyePosition(e))
	dist := e.Position().Sub(m.Position()).Len()
	if dist > g.AttackRadius {
		if g.repathIn--; g.repathIn <= 0 {
			g.repathIn = 10
			m.MoveTo(cube.PosFromVec3(e.Position()), g.speed())
		}
	} else {
		m.StopNavigating()
	}

	if g.burstLeft > 0 {
		if g.burstIn--; g.burstIn <= 0 {
			g.burstIn = max(int(g.BurstInterval*20), 1)
			g.burstLeft--
			g.Shoot(m, e)
		}
		return
	}
	if g.cooldown--; g.cooldown > 0 || dist > g.AttackRadius {
		return
	}
	g.cooldown = g.interval()
	if g.BurstShots > 1 {
		g.burstLeft, g.burstIn = g.BurstShots, 0
		return
	}
	g.Shoot(m, e)
}

func (g *RangedAttackGoal) interval() int {
	minTicks, maxTicks := g.MinInterval*20, g.MaxInterval*20
	if maxTicks <= minTicks {
		return max(int(minTicks), 20)
	}
	return int(minTicks) + rand.IntN(int(maxTicks-minTicks))
}

func (g *RangedAttackGoal) speed() float64 {
	if g.SpeedMultiplier == 0 {
		return 1
	}
	return g.SpeedMultiplier
}

// ShootArrow fires an arrow at the target, the vanilla minecraft:shooter
// component with an arrow.
func ShootArrow(m *Mob, target world.Entity) {
	shootProjectile(m, target, func(opts world.EntitySpawnOpts) *world.EntityHandle {
		return NewArrowWithDamage(opts, 2, m)
	}, 3.2)
}

// ShootSmallFireball fires a small fireball at the target.
func ShootSmallFireball(m *Mob, target world.Entity) {
	shootProjectile(m, target, func(opts world.EntitySpawnOpts) *world.EntityHandle {
		return NewSnowball(opts, m)
	}, 1.6)
}

func shootProjectile(m *Mob, target world.Entity, create func(world.EntitySpawnOpts) *world.EntityHandle, speed float64) {
	from := m.EyePosition()
	to := EyePosition(target)
	dir := to.Sub(from)
	dir[1] += math.Hypot(dir[0], dir[2]) * 0.2
	if dir.Len() < 1e-4 {
		return
	}
	vel := dir.Normalize().Mul(speed)
	m.tx.AddEntity(create(world.EntitySpawnOpts{Position: from, Velocity: vel, Rotation: m.Rotation()}))
	for _, v := range m.tx.Viewers(m.Position()) {
		v.ViewEntityAction(m, SwingArmAction{})
	}
}

// SwellGoal makes a creeper swell up and explode once a target comes close, the
// vanilla minecraft:behavior.swell component.
type SwellGoal struct {
	GoalBase
	// StartDistance is the distance at which the creeper starts swelling and
	// StopDistance the distance at which it stops again.
	StartDistance, StopDistance float64
	// FuseTicks is the time the creeper takes to explode once it started
	// swelling, the vanilla minecraft:explode fuse_length.
	FuseTicks int
	// Power is the size of the explosion.
	Power float64

	fuse int
}

func (g *SwellGoal) CanStart(m *Mob) bool {
	e, ok := m.Target()
	return ok && e.Position().Sub(m.Position()).Len() <= g.StartDistance
}

func (g *SwellGoal) CanContinue(m *Mob) bool {
	e, ok := m.Target()
	return ok && e.Position().Sub(m.Position()).Len() <= g.StopDistance
}

func (g *SwellGoal) Start(m *Mob) {
	g.fuse = g.fuseTicks()
	m.SetSwelling(true)
}

func (g *SwellGoal) Stop(m *Mob) {
	g.fuse = 0
	m.SetSwelling(false)
}

func (g *SwellGoal) Tick(m *Mob) {
	if g.fuse--; g.fuse > 0 {
		return
	}
	power := g.Power
	if power == 0 {
		power = 3
	}
	m.SetSwelling(false)
	m.Explode(power)
}

func (g *SwellGoal) fuseTicks() int {
	if g.FuseTicks == 0 {
		return 30
	}
	return g.FuseTicks
}

// FleeSunGoal walks the mob into the shade during the day, the vanilla
// minecraft:behavior.flee_sun component.
type FleeSunGoal struct {
	GoalBase
	SpeedMultiplier float64

	scanIn int
}

func (g *FleeSunGoal) CanStart(m *Mob) bool {
	if !scanTick(&g.scanIn) || !inSunlight(m) {
		return false
	}
	pos, ok := g.shade(m)
	return ok && m.MoveTo(pos, g.speed())
}

func (g *FleeSunGoal) CanContinue(m *Mob) bool { return m.Navigating() && inSunlight(m) }

func (g *FleeSunGoal) Stop(m *Mob) { m.StopNavigating() }

func (g *FleeSunGoal) shade(m *Mob) (cube.Pos, bool) {
	pos := cube.PosFromVec3(m.Position())
	for _, offset := range randomOffsets(12) {
		try := pos.Add(offset)
		if walkable(m, try) && m.tx.SkyLight(try) < 15 {
			return try, true
		}
	}
	return cube.Pos{}, false
}

func (g *FleeSunGoal) speed() float64 {
	if g.SpeedMultiplier == 0 {
		return 1
	}
	return g.SpeedMultiplier
}

// AvoidMobTypeGoal runs the mob away from entities it is afraid of, the vanilla
// minecraft:behavior.avoid_mob_type component.
type AvoidMobTypeGoal struct {
	GoalBase
	// Avoided are the encoded names of the entities the mob runs away from.
	Avoided []string
	// MaxDistance is the distance at which the mob notices an entity.
	MaxDistance float64
	// WalkSpeedMultiplier is used while walking away and SprintSpeedMultiplier
	// once the entity is close.
	WalkSpeedMultiplier, SprintSpeedMultiplier float64

	scanIn int
}

func (g *AvoidMobTypeGoal) CanStart(m *Mob) bool {
	if !scanTick(&g.scanIn) {
		return false
	}
	e, ok := g.nearest(m)
	if !ok {
		return false
	}
	away, ok := g.destination(m, e.Position())
	return ok && m.MoveTo(away, g.speed())
}

func (g *AvoidMobTypeGoal) CanContinue(m *Mob) bool {
	if !m.Navigating() {
		return false
	}
	_, ok := g.nearest(m)
	return ok
}

func (g *AvoidMobTypeGoal) Stop(m *Mob) { m.StopNavigating() }

func (g *AvoidMobTypeGoal) destination(m *Mob, from mgl64.Vec3) (cube.Pos, bool) {
	away := m.Position().Sub(from)
	away[1] = 0
	if away.Len() < 1e-4 {
		return cube.Pos{}, false
	}
	target := m.Position().Add(away.Normalize().Mul(12))
	return nearestWalkable(m, cube.PosFromVec3(target))
}

func (g *AvoidMobTypeGoal) nearest(m *Mob) (world.Entity, bool) {
	d := g.MaxDistance
	if d == 0 {
		d = 6
	}
	pos, best := m.Position(), d*d
	var closest world.Entity
	for e := range m.tx.EntitiesWithin(cube.Box(-d, -d, -d, d, d, d).Translate(pos)) {
		name := e.H().Type().EncodeEntity()
		if !containsName(g.Avoided, name) {
			continue
		}
		if dist := e.Position().Sub(pos).LenSqr(); dist < best {
			best, closest = dist, e
		}
	}
	return closest, closest != nil
}

func (g *AvoidMobTypeGoal) speed() float64 {
	if g.SprintSpeedMultiplier > 0 {
		return g.SprintSpeedMultiplier
	}
	if g.WalkSpeedMultiplier > 0 {
		return g.WalkSpeedMultiplier
	}
	return 1
}

func containsName(names []string, name string) bool {
	for _, n := range names {
		if n == name {
			return true
		}
	}
	return false
}

func inSunlight(m *Mob) bool {
	tx := m.tx
	if tx.World().Dimension() != world.Overworld {
		return false
	}
	if day := tx.World().Time() % 24000; day > 12600 {
		return false
	}
	pos := cube.PosFromVec3(m.EyePosition())
	return tx.SkyLight(pos) >= 15
}

func (RangedAttackGoal) Controls() Control { return ControlMove | ControlLook }

func (SwellGoal) Controls() Control { return ControlMove }

func (FleeSunGoal) Controls() Control { return ControlMove | ControlLook }

func (AvoidMobTypeGoal) Controls() Control { return ControlMove | ControlLook }

func (g *RangedAttackGoal) Clone() Goal { c := *g; return &c }

func (g *SwellGoal) Clone() Goal { c := *g; return &c }

func (g *FleeSunGoal) Clone() Goal { c := *g; return &c }

func (g *AvoidMobTypeGoal) Clone() Goal { c := *g; return &c }
