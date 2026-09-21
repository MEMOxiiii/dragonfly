package entity

import (
	"math"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

// Mob is a living entity driven by a set of Goal implementations.
type Mob struct {
	*Ent
}

func (m *Mob) mob() *mobBehaviour {
	return m.Behaviour().(*mobBehaviour)
}

func (m *Mob) Tx() *world.Tx { return m.tx }

func (m *Mob) Health() float64    { return m.mob().health.Health() }
func (m *Mob) MaxHealth() float64 { return m.mob().health.MaxHealth() }

func (m *Mob) SetMaxHealth(v float64) {
	m.mob().health.SetMaxHealth(v)
	m.mob().sendHealth(m)
}

func (m *Mob) Dead() bool { return m.Health() <= mgl64.Epsilon }

func (m *Mob) Hurt(damage float64, src world.DamageSource) (float64, bool) {
	return m.mob().hurt(m, damage, src)
}

func (m *Mob) Heal(health float64, src world.HealingSource) float64 {
	b := m.mob()
	if m.Dead() || health < 0 {
		return 0
	}
	before := b.health.Health()
	b.health.AddHealth(health)
	b.sendHealth(m)
	return b.health.Health() - before
}

func (m *Mob) KnockBack(src mgl64.Vec3, force, height float64) {
	if m.Dead() {
		return
	}
	vel := m.Position().Sub(src)
	vel[1] = 0
	if vel.Len() != 0 {
		vel = vel.Normalize().Mul(force)
	}
	vel[1] = height
	m.SetVelocity(vel)
}

func (m *Mob) AddEffect(e effect.Effect) {
	m.mob().effects.Add(e, m)
	m.updateState()
}

func (m *Mob) RemoveEffect(e effect.Type) {
	m.mob().effects.Remove(e, m)
	m.updateState()
}

func (m *Mob) Effects() []effect.Effect { return m.mob().effects.Effects() }

func (m *Mob) Speed() float64     { return m.mob().speed }
func (m *Mob) SetSpeed(v float64) { m.mob().speed = v }

// Target returns the entity the mob is attacking, if it has one.
func (m *Mob) Target() (world.Entity, bool) {
	b := m.mob()
	if b.target == nil {
		return nil, false
	}
	e, ok := b.target.Entity(m.tx)
	if !ok {
		b.target = nil
	}
	return e, ok
}

// SetTarget sets the entity the mob attacks. Nil clears the current target.
func (m *Mob) SetTarget(e world.Entity) {
	if e == nil {
		m.mob().target = nil
		return
	}
	m.mob().target = e.H()
}

// MoveTo searches a path to the position passed and starts following it, moving
// at the mob's speed multiplied by the multiplier passed.
func (m *Mob) MoveTo(pos cube.Pos, speedMultiplier float64) bool {
	return m.mob().nav.moveTo(m, pos, speedMultiplier)
}

// Navigating reports whether the mob is currently following a path.
func (m *Mob) Navigating() bool { return m.mob().nav.navigating() }

// StopNavigating clears the path the mob is following.
func (m *Mob) StopNavigating() { m.mob().nav.stop() }

// Move makes the mob walk in the direction passed for a single tick, at the
// mob's speed multiplied by the multiplier passed.
func (m *Mob) Move(dir mgl64.Vec3, speedMultiplier float64) {
	dir[1] = 0
	if l := dir.Len(); l > 1e-4 {
		b := m.mob()
		b.moveInput, b.moveSpeed = dir.Mul(1/l), speedMultiplier
	}
}

// MoveTowards moves the mob towards the position passed in all three axes. It
// is used by mobs that fly or swim.
func (m *Mob) MoveTowards(pos mgl64.Vec3, speedMultiplier float64) {
	dir := pos.Sub(m.Position())
	if l := dir.Len(); l > 1e-4 {
		b := m.mob()
		b.moveInput, b.moveSpeed = dir.Mul(1/l), speedMultiplier
		m.LookAt(pos)
	}
}

// Explode makes the mob explode with the power passed, removing it.
func (m *Mob) Explode(power float64) {
	b := m.mob()
	b.exploding = true
	pos, tx := m.Position(), m.tx
	_ = m.Close()
	block.ExplosionConfig{SuppressUnderwaterImpact: true}.Explode(tx, explosionSource{pos: pos, size: power})
}

// Swelling reports whether the mob is a creeper about to explode.
func (m *Mob) Swelling() bool { return m.mob().swelling }

// SetSwelling makes a creeper swell up or stop swelling again.
func (m *Mob) SetSwelling(swelling bool) {
	m.mob().swelling = swelling
	m.updateState()
}

// InLiquid reports whether the mob is inside a liquid.
func (m *Mob) InLiquid() bool {
	_, ok := m.tx.Liquid(cube.PosFromVec3(m.Position()))
	return ok
}

// Jump makes the mob jump on the next movement tick if it is on the ground.
func (m *Mob) Jump() { m.mob().jumping = true }

func (m *Mob) OnGround() bool { return m.mob().mc.OnGround() }

// LookAt turns the mob's head and body towards the position passed.
func (m *Mob) LookAt(pos mgl64.Vec3) {
	diff := pos.Sub(m.EyePosition())
	horizontal := math.Hypot(diff[0], diff[2])
	m.SetRotation(cube.Rotation{
		mgl64.RadToDeg(math.Atan2(diff[2], diff[0])) - 90,
		-mgl64.RadToDeg(math.Atan2(diff[1], horizontal)),
	})
}

// SetRotation sets the rotation of the mob, which is sent to viewers on the
// next movement tick.
func (m *Mob) SetRotation(rot cube.Rotation) {
	m.data.Rot = rot
}

func (m *Mob) EyePosition() mgl64.Vec3 {
	return m.Position().Add(mgl64.Vec3{0, m.eyeHeight(), 0})
}

func (m *Mob) eyeHeight() float64 {
	return m.H().Type().BBox(m).Height() * 0.85
}

type explosionSource struct {
	pos  mgl64.Vec3
	size float64
}

func (s explosionSource) Position() mgl64.Vec3 { return s.pos }
func (s explosionSource) Size() float64        { return s.size }

func (m *Mob) updateState() {
	for _, v := range m.tx.Viewers(m.Position()) {
		v.ViewEntityState(m)
	}
}
