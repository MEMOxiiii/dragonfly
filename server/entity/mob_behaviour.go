package entity

import (
	"math"
	"math/rand/v2"
	"time"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"

	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/google/uuid"
)

const (
	mobImmunity   = time.Millisecond * 500
	mobDeathTicks = 20
	mobJumpHeight = 0.42
	// mobSpeedFactor converts the vanilla minecraft:movement value into the
	// speed a mob keeps while walking, measured on a vanilla server.
	mobSpeedFactor = 0.54
	// The factors below do the same for the other ways a mob travels, each
	// measured on a vanilla server: a cod for swimming, a parrot for flying, a
	// bee for hovering and a bat for floating.
	mobSwimFactor  = 0.73
	mobFlyFactor   = 0.42
	mobHoverFactor = 1.5
	mobFloatFactor = 6.1
	// mobTravelDrag is the velocity a swimming or flying mob loses every tick.
	mobTravelDrag = 0.09
)

// Travel is the way a mob moves through the world, the vanilla
// minecraft:navigation components.
type Travel uint8

const (
	// TravelWalk is a mob that walks over blocks, minecraft:navigation.walk.
	TravelWalk Travel = iota
	// TravelSwim is a mob that swims through liquid, minecraft:navigation.generic.
	TravelSwim
	// TravelFly is a mob that flies to a target, minecraft:navigation.fly.
	TravelFly
	// TravelHover is a mob that hovers in place, minecraft:navigation.hover.
	TravelHover
	// TravelFloat is a mob that floats around erratically, the vanilla
	// minecraft:navigation.float.
	TravelFloat
	// TravelClimb is a mob that walks and climbs walls, minecraft:navigation.climb.
	TravelClimb
)

// airborne reports whether a mob of this Travel moves through the air.
func (t Travel) airborne() bool {
	return t == TravelFly || t == TravelHover || t == TravelFloat
}

// factor is the number the movement value of a mob is multiplied with to get
// the speed it travels at.
func (t Travel) factor() float64 {
	switch t {
	case TravelSwim:
		return mobSwimFactor
	case TravelFly:
		return mobFlyFactor
	case TravelHover:
		return mobHoverFactor
	case TravelFloat:
		return mobFloatFactor
	}
	return mobSpeedFactor
}

// MobConfig holds the parameters of a Mob, matching the components of
// its vanilla entity definition.
type MobConfig struct {
	MaxHealth float64
	// Speed is the movement speed of the mob in blocks per tick, the vanilla
	// minecraft:movement component.
	Speed float64
	// Travel is the way the mob moves through the world.
	Travel Travel
	// KnockBackResistance ranges from 0 to 1 and reduces the knock back the mob
	// takes by that fraction.
	KnockBackResistance float64
	Gravity, Drag       float64
	// Goals are the goals the mob runs, the vanilla minecraft:behavior
	// components.
	Goals []Goal
	// Drops is the loot table rolled when the mob dies.
	Drops LootTable
	// BreedItems are the encoded names of the items the mob breeds with, the
	// vanilla minecraft:breedable component.
	BreedItems []string
	// BurnsInDaylight sets the mob on fire while it stands in sunlight, the
	// vanilla minecraft:burns_in_daylight component.
	BurnsInDaylight bool
	// FireImmune stops the mob from taking fire damage, the vanilla
	// minecraft:fire_immune component.
	FireImmune bool
	// LavaDamage is the damage the mob takes every tick in lava, the vanilla
	// minecraft:hurt_on_condition component.
	LavaDamage float64
	// Breathable holds what the mob breathes, the vanilla minecraft:breathable
	// component.
	Breathable Breathable
	// Healable holds the items the mob is healed by.
	Healable Healable
	// Tameable holds the items the mob is tamed with.
	Tameable Tameable
	// Colours are the colours the mob spawns with, such as the wool colour of
	// a sheep, with the weights vanilla picks them with.
	Colours []VariantChoice
	// Variants are the variants the mob spawns as, the vanilla
	// minecraft:variant component.
	Variants []VariantChoice
	// ShearDrops is the loot table rolled when the mob is sheared.
	ShearDrops LootTable
	// ExperienceDrops is the range of experience the mob drops when killed by a
	// player.
	ExperienceDrops [2]int
}

func (conf MobConfig) Apply(data *world.EntityData) {
	data.Data = conf.New()
}

func (conf MobConfig) New() *mobBehaviour {
	if conf.MaxHealth == 0 {
		conf.MaxHealth = 20
	}
	if conf.Breathable.TotalSupply == 0 && conf.Breathable.Air {
		conf.Breathable.TotalSupply = 15
	}
	if conf.Gravity == 0 && !conf.Travel.airborne() {
		conf.Gravity = 0.08
	}
	if conf.Drag == 0 {
		conf.Drag = 0.02
	}
	if conf.Travel.airborne() {
		conf.Gravity = 0
	}
	b := &mobBehaviour{
		conf:     conf,
		health:   NewHealthManager(conf.MaxHealth, conf.MaxHealth),
		effects:  NewEffectManager(),
		speed:    conf.Speed,
		mc:       &MovementComputer{Gravity: conf.Gravity, Drag: conf.Drag, DragBeforeGravity: true},
		travelMC: &MovementComputer{Drag: mobTravelDrag},
		nav:      &navigator{},
	}
	goals := make([]Goal, len(conf.Goals))
	for i, g := range conf.Goals {
		goals[i] = g.Clone()
	}
	b.goals = newGoalSelector(goals)
	return b
}

type mobBehaviour struct {
	conf     MobConfig
	health   *HealthManager
	effects  *EffectManager
	goals    *goalSelector
	nav      *navigator
	mc       *MovementComputer
	travelMC *MovementComputer

	speed        float64
	target       *world.EntityHandle
	lastAttacker *world.EntityHandle
	panicking    bool

	immuneTill time.Time
	lastDamage float64

	swelling      bool
	exploding     bool
	baby          bool
	age           int
	love          int
	breedCooldown int
	partner       *world.EntityHandle

	moveInput            mgl64.Vec3
	moveSpeed            float64
	jumping              bool
	horizontallyCollided bool

	fireTicks    int
	airSupply    int
	breathTicks  int
	owner        *world.EntityHandle
	ownerID      uuid.UUID
	sitting      bool
	colour       uint8
	variant      int32
	sheared      bool
	saddled      bool
	fallDistance float64
	deathTicks   int

	self       *Mob
	pathHeight int
}

func (b *mobBehaviour) Tick(e *Ent, tx *world.Tx) *Movement {
	m := b.wrap(e)
	if b.health.Health() <= mgl64.Epsilon {
		if b.deathTicks++; b.deathTicks >= mobDeathTicks {
			_ = e.Close()
		}
		return nil
	}
	b.effects.Tick(m, tx)
	b.tickFire(m)
	b.tickBreathing(m)
	b.tickBlockDamage(m)
	b.tickBreeding(m)
	b.goals.tick(m)
	b.nav.tick(m)
	return b.move(m)
}

// wrap returns the Mob of the Ent passed, reusing the same value every tick.
func (b *mobBehaviour) wrap(e *Ent) *Mob {
	if b.self == nil {
		b.self = &Mob{Ent: e}
		b.pathHeight = max(int(math.Ceil(e.H().Type().BBox(b.self).Height())), 1)
		return b.self
	}
	b.self.Ent = e
	return b.self
}

func (b *mobBehaviour) move(m *Mob) *Movement {
	if b.conf.Travel.airborne() || (b.conf.Travel == TravelSwim && b.inLiquid(m)) {
		return b.travel(m)
	}
	vel := m.Velocity()
	if target := b.moveSpeed * b.speed * mobSpeedFactor * b.effectSpeed(); target > 0 {
		friction := b.friction(m)
		accel := target * (1 - friction) / friction
		if !b.mc.OnGround() {
			accel *= 0.2
		}
		vel[0] += b.moveInput[0] * accel
		vel[2] += b.moveInput[2] * accel
	}
	if b.jumping && b.mc.OnGround() {
		vel[1] = mobJumpHeight
	}
	idle := b.moveInput == (mgl64.Vec3{}) && !b.jumping
	b.moveInput, b.moveSpeed, b.jumping = mgl64.Vec3{}, 0, false
	if idle && b.resting(m, vel) {
		return nil
	}

	before := m.Position()
	mv := b.mc.TickMovement(m, before, vel, m.Rotation(), m.tx)
	b.horizontallyCollided = math.Abs(mv.vel[0]) < 1e-5 && math.Abs(mv.vel[2]) < 1e-5 && (math.Abs(vel[0]) > 1e-5 || math.Abs(vel[2]) > 1e-5)
	b.tickFall(m, before[1]-mv.pos[1])

	m.data.Pos, m.data.Vel = mv.pos, mv.vel
	return mv
}

func (b *mobBehaviour) friction(m *Mob) float64 {
	friction := 1 - b.conf.Drag
	if b.mc.OnGround() {
		if f, ok := m.tx.Block(cube.PosFromVec3(m.Position()).Side(cube.FaceDown)).(interface {
			Friction() float64
		}); ok {
			return friction * f.Friction()
		}
		return friction * 0.6
	}
	return friction
}

// travel moves a mob that is not bound to the ground, such as a fish in water
// or a flying mob, towards its move input in all three axes.
func (b *mobBehaviour) travel(m *Mob) *Movement {
	factor := b.conf.Travel.factor()
	target := b.moveInput.Mul(b.moveSpeed * b.speed * factor * b.effectSpeed())
	b.moveInput, b.moveSpeed, b.jumping = mgl64.Vec3{}, 0, false

	vel := m.Velocity()
	vel = vel.Add(target.Sub(vel).Mul(0.2))
	if vel.LenSqr() < 1e-8 {
		return nil
	}
	mv := b.travelMC.TickMovement(m, m.Position(), vel, m.Rotation(), m.tx)
	m.data.Pos, m.data.Vel = mv.pos, mv.vel
	return mv
}

func (b *mobBehaviour) inLiquid(m *Mob) bool {
	_, ok := m.tx.Liquid(cube.PosFromVec3(m.Position()))
	return ok
}

// resting reports whether the mob stands still on solid ground, in which case
// its collision does not have to be computed at all.
func (b *mobBehaviour) resting(m *Mob, vel mgl64.Vec3) bool {
	if !b.mc.OnGround() || vel.LenSqr() > 1e-8 {
		return false
	}
	pos := m.Position()
	below := cube.PosFromVec3(pos).Side(cube.FaceDown)
	for _, box := range m.tx.Block(below).Model().BBox(below, m.tx) {
		if box.Max()[1] >= pos[1]-float64(below[1])-1e-4 {
			return true
		}
	}
	return false
}

func (b *mobBehaviour) effectSpeed() float64 {
	multiplier := 1.0
	if e, ok := b.effects.Effect(effect.Speed); ok {
		multiplier += 0.2 * float64(e.Level())
	}
	if e, ok := b.effects.Effect(effect.Slowness); ok {
		multiplier -= 0.15 * float64(e.Level())
	}
	return max(multiplier, 0)
}

// OnFireDuration returns how long the mob keeps burning.
func (m *Mob) OnFireDuration() time.Duration {
	return time.Duration(m.mob().fireTicks) * time.Second / 20
}

// SetOnFire sets the mob on fire for the duration passed.
func (m *Mob) SetOnFire(duration time.Duration) {
	ticks := int(duration.Seconds() * 20)
	if ticks > m.mob().fireTicks {
		m.mob().fireTicks = ticks
		m.updateState()
	}
}

func (b *mobBehaviour) tickFire(m *Mob) {
	if b.conf.FireImmune {
		b.fireTicks = 0
		return
	}
	if b.conf.BurnsInDaylight && b.fireTicks <= 0 && inSunlight(m) {
		m.SetOnFire(8 * time.Second)
	}
	if b.fireTicks <= 0 {
		return
	}
	if _, ok := m.tx.Liquid(cube.PosFromVec3(m.Position())); ok {
		b.fireTicks = 0
		m.updateState()
		return
	}
	if b.fireTicks--; b.fireTicks%20 == 0 {
		b.hurt(m, 1, block.FireDamageSource{})
	}
	if b.fireTicks == 0 {
		m.updateState()
	}
}

func (b *mobBehaviour) tickFall(m *Mob, dy float64) {
	if b.mc.OnGround() {
		if b.fallDistance > 0 {
			if damage := math.Floor(b.fallDistance - 3); damage > 0 {
				b.hurt(m, damage, FallDamageSource{})
			}
			b.fallDistance = 0
		}
		return
	}
	if dy > 0 {
		b.fallDistance += dy
	} else {
		b.fallDistance = 0
	}
}

func (b *mobBehaviour) hurt(m *Mob, damage float64, src world.DamageSource) (float64, bool) {
	if m.Dead() || damage < 0 {
		return 0, false
	}
	immune := time.Now().Before(b.immuneTill)
	if immune && damage <= b.lastDamage {
		return 0, false
	}
	if immune {
		damage -= b.lastDamage
	}
	if damage = b.finalDamage(damage, src); damage <= 0 {
		return 0, true
	}
	b.health.AddHealth(-damage)
	b.sendHealth(m)

	if !immune {
		b.immuneTill = time.Now().Add(mobImmunity)
		b.lastDamage = damage
	}
	for _, v := range m.tx.Viewers(m.Position()) {
		v.ViewEntityAction(m, HurtAction{})
	}
	b.panicking = true
	if attack, ok := src.(AttackDamageSource); ok && attack.Attacker != nil {
		m.KnockBack(attack.Attacker.Position(), 0.4*(1-b.conf.KnockBackResistance), 0.3608*(1-b.conf.KnockBackResistance))
		b.lastAttacker = attack.Attacker.H()
	}
	if m.Dead() {
		b.kill(m, src)
	}
	return damage, true
}

func (b *mobBehaviour) finalDamage(damage float64, src world.DamageSource) float64 {
	if b.conf.FireImmune && src.Fire() {
		return 0
	}
	if res, ok := b.effects.Effect(effect.Resistance); ok {
		damage *= effect.Resistance.Multiplier(src, res.Level())
	}
	if _, ok := b.effects.Effect(effect.FireResistance); ok && src.Fire() {
		return 0
	}
	return damage
}

func (b *mobBehaviour) kill(m *Mob, src world.DamageSource) {
	b.goals.stop(m)
	b.nav.stop()
	b.target = nil
	for _, v := range m.tx.Viewers(m.Position()) {
		v.ViewEntityAction(m, DeathAction{})
	}
	attack, byAttack := src.(AttackDamageSource)
	for _, drop := range b.conf.Drops.roll(b.fireTicks > 0, lootingLevel(src)) {
		if name, _ := drop.Item().EncodeItem(); containsName(woolItems, name) && len(b.conf.Colours) > 0 {
			drop = colouredWool(m.Colour(), drop.Count())
		}
		opts := world.EntitySpawnOpts{Position: m.Position().Add(dropOffset), Velocity: dropVelocity()}
		m.tx.AddEntity(NewItem(opts, drop))
	}
	if xp := b.conf.ExperienceDrops; xp[1] > 0 && byAttack && attack.Attacker != nil {
		amount := xp[0]
		if xp[1] > xp[0] {
			amount += rand.IntN(xp[1] - xp[0] + 1)
		}
		for _, orb := range NewExperienceOrbs(m.Position(), amount) {
			m.tx.AddEntity(orb)
		}
	}
}

func (b *mobBehaviour) sendHealth(m *Mob) {
	m.updateState()
}

func lootingLevel(src world.DamageSource) int {
	attack, ok := src.(AttackDamageSource)
	if !ok || attack.Attacker == nil {
		return 0
	}
	holder, ok := attack.Attacker.(interface {
		HeldItems() (item.Stack, item.Stack)
	})
	if !ok {
		return 0
	}
	held, _ := holder.HeldItems()
	if e, ok := held.Enchantment(enchantment.Looting); ok {
		return e.Level()
	}
	return 0
}

var dropOffset = mgl64.Vec3{0, 0.25, 0}

func dropVelocity() mgl64.Vec3 {
	return mgl64.Vec3{rand.Float64()*0.2 - 0.1, 0.2, rand.Float64()*0.2 - 0.1}
}

// ownerUUID returns the UUID of the entity that tamed the mob.
func (b *mobBehaviour) ownerUUID() uuid.UUID {
	if b.owner != nil {
		return b.owner.UUID()
	}
	return b.ownerID
}
