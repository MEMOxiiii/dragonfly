package entity

import (
	"math/rand/v2"
	"time"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/particle"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
)

// FishingBobberBehaviourConfig holds optional parameters for a FishingBobberBehaviour.
type FishingBobberBehaviourConfig struct {
	// Owner is the entity that cast the bobber out. Only the owner may reel it back in.
	Owner *world.EntityHandle
	// Lure is the level of the Lure enchantment on the rod that cast the bobber.
	Lure int
	// LuckOfTheSea is the level of the Luck of the Sea enchantment on the rod that cast the bobber.
	LuckOfTheSea int
	// Gravity is the amount of Y velocity subtracted every tick while the bobber is airborne.
	Gravity float64
	// Drag is used to reduce all axes of the velocity every tick while the bobber is airborne.
	Drag float64
}

func (conf FishingBobberBehaviourConfig) Apply(data *world.EntityData) {
	data.Data = conf.New()
}

// New creates a FishingBobberBehaviour using the parameters in conf.
func (conf FishingBobberBehaviourConfig) New() *FishingBobberBehaviour {
	if conf.Gravity == 0 {
		conf.Gravity = 0.03
	}
	if conf.Drag == 0 {
		conf.Drag = 0.01
	}
	b := &FishingBobberBehaviour{conf: conf}
	b.passive = PassiveBehaviourConfig{
		Gravity: conf.Gravity,
		Drag:    conf.Drag,
		Tick:    b.tick,
	}.New()
	return b
}

// FishingBobberBehaviour implements the behaviour of a fishing rod's bobber. The bobber falls until it lands
// in water, floats there while a catch is waited for and may be reeled back in by its owner at any point.
type FishingBobberBehaviour struct {
	conf FishingBobberBehaviourConfig

	passive *PassiveBehaviour

	floating bool
	waitTime time.Duration
	hooked   *world.EntityHandle
}

// Owner returns the handle of the entity that cast the bobber out.
func (f *FishingBobberBehaviour) Owner() *world.EntityHandle {
	return f.conf.Owner
}

// Hooked returns the entity the bobber has attached itself to, if any.
func (f *FishingBobberBehaviour) Hooked() *world.EntityHandle {
	return f.hooked
}

// Biting returns whether a catch is currently waiting to be reeled in.
func (f *FishingBobberBehaviour) Biting() bool {
	return f.floating && f.waitTime <= 0
}

// LuckOfTheSea returns the level of Luck of the Sea on the rod that cast the bobber.
func (f *FishingBobberBehaviour) LuckOfTheSea() int {
	return f.conf.LuckOfTheSea
}

// Tick ...
func (f *FishingBobberBehaviour) Tick(e *Ent, tx *world.Tx) *Movement {
	return f.passive.Tick(e, tx)
}

// Explode ...
func (f *FishingBobberBehaviour) Explode(e *Ent, src world.ExplosionSource, impact float64) {
	f.passive.Explode(e, src, impact)
}

// PortalTravelComputer ...
func (f *FishingBobberBehaviour) PortalTravelComputer() *PortalTravelComputer {
	return f.passive.PortalTravelComputer()
}

// tick runs the bobber's state machine: it hooks entities it passes through while airborne, settles once it
// reaches water and counts down to a catch while floating.
func (f *FishingBobberBehaviour) tick(e *Ent, tx *world.Tx) {
	owner, ok := f.conf.Owner.Entity(tx)
	if !ok || e.Position().Sub(owner.Position()).Len() > 33 {
		_ = e.Close()
		return
	}
	if f.hooked != nil {
		if hooked, ok := f.hooked.Entity(tx); ok {
			e.Teleport(hooked.Position())
			return
		}
		f.hooked = nil
	}
	if !f.floating {
		if f.hookEntity(e, tx, owner) {
			return
		}
		if _, ok := tx.Liquid(cube.PosFromVec3(e.Position())); !ok {
			return
		}
		f.floating = true
		e.SetVelocity(mgl64.Vec3{})
		f.waitTime = f.newWaitTime(e, tx)
		return
	}
	f.tickFloating(e, tx)
}

// hookEntity attaches the bobber to the first entity other than its owner that it overlaps with, returning
// whether one was found.
func (f *FishingBobberBehaviour) hookEntity(e *Ent, tx *world.Tx, owner world.Entity) bool {
	for other := range tx.EntitiesWithin(e.H().Type().BBox(e).Translate(e.Position()).Grow(0.25)) {
		if other == e || other == owner {
			continue
		}
		if _, ok := other.(Living); !ok {
			continue
		}
		f.hooked = other.H()
		e.SetVelocity(mgl64.Vec3{})
		return true
	}
	return false
}

// tickFloating counts down towards a catch and plays the sound and particles announcing one.
func (f *FishingBobberBehaviour) tickFloating(e *Ent, tx *world.Tx) {
	if f.waitTime <= 0 {
		return
	}
	step := time.Second / 20
	if tx.RainingAt(cube.PosFromVec3(e.Position())) && rand.Float64() < 0.25 {
		// Rain makes each tick count down twice as often a quarter of the time, cutting the wait by ~20%.
		step *= 2
	}
	if f.waitTime -= step; f.waitTime > 0 {
		return
	}
	f.waitTime = 0
	pos := e.Position()
	tx.AddParticle(pos, particle.Splash{})
	tx.PlaySound(pos, sound.Custom{Name: "random.splash", Volume: 0.25, Pitch: 1})
}

// newWaitTime rolls the time until something bites the hook. The base range is 5 to 30 seconds, shortened by
// 5 seconds for every level of Lure and doubled when the bobber cannot see the sky.
func (f *FishingBobberBehaviour) newWaitTime(e *Ent, tx *world.Tx) time.Duration {
	lure := min(f.conf.Lure, 5) * 100
	minTicks, maxTicks := max(100-lure, 0), max(600-lure, 0)
	ticks := minTicks
	if maxTicks > minTicks {
		ticks += rand.IntN(maxTicks - minTicks)
	}
	wait := time.Duration(ticks) * time.Second / 20
	if pos := cube.PosFromVec3(e.Position()); tx.SkyLight(pos) == 0 {
		wait *= 2
	}
	return wait
}

// openWater reports whether the bobber rests in open water, which vanilla defines as a 5x4x5 volume centred
// on it holding nothing but air, water and blocks that do not obstruct either. Catches outside open water are
// limited to the fish table.
func openWater(pos cube.Pos, tx *world.Tx) bool {
	for x := -2; x <= 2; x++ {
		for z := -2; z <= 2; z++ {
			for y := -2; y <= 1; y++ {
				p := pos.Add(cube.Pos{x, y, z})
				if _, ok := tx.Liquid(p); ok {
					continue
				}
				if _, ok := tx.Block(p).(interface{ HasLiquidDrops() bool }); ok {
					continue
				}
				if len(tx.Block(p).Model().BBox(p, tx)) != 0 {
					return false
				}
			}
		}
	}
	return true
}
