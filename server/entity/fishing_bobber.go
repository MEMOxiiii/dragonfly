package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// NewFishingBobber creates a FishingBobber entity. FishingBobber is the hook cast out by a fishing rod. It
// falls until it lands in water, floats there while a catch is waited for and is reeled back in by its owner.
func NewFishingBobber(opts world.EntitySpawnOpts, owner world.Entity, lure int, luckOfTheSea int) *world.EntityHandle {
	return opts.New(FishingBobberType, FishingBobberBehaviourConfig{
		Owner:        owner.H(),
		Lure:         lure,
		LuckOfTheSea: luckOfTheSea,
	})
}

// FishingBobber is the hook cast out by a fishing rod.
type FishingBobber struct {
	*Ent
}

// behaviour returns the bobber's behaviour.
func (f *FishingBobber) behaviour() *FishingBobberBehaviour {
	return f.Behaviour().(*FishingBobberBehaviour)
}

// OwnedBy returns whether the bobber was cast out by the entity passed.
func (f *FishingBobber) OwnedBy(h *world.EntityHandle) bool {
	return f.behaviour().conf.Owner == h
}

// Biting returns whether a catch is waiting to be reeled in.
func (f *FishingBobber) Biting() bool {
	return f.behaviour().Biting()
}

// Hooked returns the entity the bobber has attached itself to, if any.
func (f *FishingBobber) Hooked() (world.Entity, bool) {
	h := f.behaviour().hooked
	if h == nil {
		return nil, false
	}
	return h.Entity(f.tx)
}

// LuckOfTheSea returns the level of Luck of the Sea on the rod that cast the bobber.
func (f *FishingBobber) LuckOfTheSea() int {
	return f.behaviour().conf.LuckOfTheSea
}

// OpenWater returns whether the bobber rests in open water, which is required to catch treasure.
func (f *FishingBobber) OpenWater() bool {
	return openWater(cube.PosFromVec3(f.Position()), f.tx)
}

// FishingBobberType is a world.EntityType implementation for FishingBobber.
var FishingBobberType fishingBobberType

type fishingBobberType struct{}

func (fishingBobberType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &FishingBobber{Ent: &Ent{tx: tx, handle: handle, data: data}}
}

func (fishingBobberType) EncodeEntity() string { return "minecraft:fishing_hook" }

func (fishingBobberType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.125, 0, -0.125, 0.125, 0.25, 0.125)
}

func (fishingBobberType) DecodeNBT(_ map[string]any, data *world.EntityData) {
	data.Data = FishingBobberBehaviourConfig{}.New()
}

func (fishingBobberType) EncodeNBT(*world.EntityData) map[string]any { return nil }
