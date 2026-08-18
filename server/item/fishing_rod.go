package item

import (
	"math/rand/v2"
	"time"

	"github.com/df-mc/dragonfly/server/item/potion"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
)

// FishingRod is a tool used to catch fish and other items out of bodies of water, and to pull entities
// towards the user.
type FishingRod struct{}

// MaxCount ...
func (FishingRod) MaxCount() int {
	return 1
}

// DurabilityInfo ...
func (FishingRod) DurabilityInfo() DurabilityInfo {
	return DurabilityInfo{
		MaxDurability: 385,
		BrokenItem:    simpleItem(Stack{}),
	}
}

// FuelInfo ...
func (FishingRod) FuelInfo() FuelInfo {
	return newFuelInfo(time.Second * 15)
}

// EnchantmentValue ...
func (FishingRod) EnchantmentValue() int {
	return 1
}

// RepairableBy ...
func (FishingRod) RepairableBy(i Stack) bool {
	name, _ := i.Item().EncodeItem()
	return name == "minecraft:string"
}

// Use casts the bobber out if the user has none out yet, and reels it back in if they do.
func (f FishingRod) Use(tx *world.Tx, user User, ctx *UseContext) bool {
	if bobber, ok := castBobber(tx, user); ok {
		f.reel(bobber, tx, user, ctx)
		return true
	}
	f.cast(tx, user, ctx)
	return true
}

// cast spawns a bobber travelling away from the user.
func (FishingRod) cast(tx *world.Tx, user User, ctx *UseContext) {
	held, _ := user.HeldItems()
	lure := 0
	for _, e := range held.Enchantments() {
		if _, ok := e.Type().(interface {
			WaitTimeReduction(int) time.Duration
		}); ok {
			lure = e.Level()
		}
	}
	luck := 0
	for _, e := range held.Enchantments() {
		if _, ok := e.Type().(interface {
			FishingChances(int) (float64, float64)
		}); ok {
			luck = e.Level()
		}
	}

	rot := user.Rotation()
	dir := rot.Vec3().Mul(1.5)
	opts := world.EntitySpawnOpts{
		Position: eyePosition(user),
		Velocity: mgl64.Vec3{dir[0], dir[1] + 0.1, dir[2]},
		Rotation: rot,
	}
	tx.AddEntity(tx.World().EntityRegistry().Config().FishingBobber(opts, user, lure, luck))
	tx.PlaySound(user.Position(), sound.Custom{Name: "random.bow", Volume: 0.5, Pitch: 0.4})
	ctx.DamageItem(0)
}

// reel pulls the bobber back in, dragging a hooked entity along with it or spawning a catch when one is
// waiting to be reeled in.
func (f FishingRod) reel(bobber fishingBobber, tx *world.Tx, user User, ctx *UseContext) {
	switch hooked, ok := bobber.Hooked(); {
	case ok:
		if v, ok := hooked.(interface{ SetVelocity(mgl64.Vec3) }); ok {
			// A reeled in entity is pulled towards the user at a tenth of the distance between the two.
			v.SetVelocity(user.Position().Sub(hooked.Position()).Mul(0.1))
		}
		ctx.DamageItem(5)
	case bobber.Biting():
		f.catch(bobber, tx, user)
		ctx.DamageItem(1)
	default:
		ctx.DamageItem(0)
	}
	tx.PlaySound(user.Position(), sound.Custom{Name: "random.bow", Volume: 0.5, Pitch: 0.4})
	_ = bobber.Close()
}

// catch spawns the caught item flying towards the user, along with the experience the catch is worth.
func (FishingRod) catch(bobber fishingBobber, tx *world.Tx, user User) {
	conf := tx.World().EntityRegistry().Config()
	pos := bobber.Position()

	opts := world.EntitySpawnOpts{Position: pos, Velocity: user.Position().Sub(pos).Mul(0.1)}
	tx.AddEntity(conf.Item(opts, fishingLoot(bobber)))
	for _, orb := range conf.ExperienceOrbs(user.Position(), rand.IntN(6)+1) {
		tx.AddEntity(orb)
	}
}

// fishingBobber is the hook cast out by a fishing rod, implemented by entity.FishingBobber.
type fishingBobber interface {
	world.Entity
	// OwnedBy returns whether the bobber was cast out by the entity passed.
	OwnedBy(*world.EntityHandle) bool
	// Biting returns whether a catch is waiting to be reeled in.
	Biting() bool
	// Hooked returns the entity the bobber has attached itself to, if any.
	Hooked() (world.Entity, bool)
	// LuckOfTheSea returns the level of Luck of the Sea on the rod that cast the bobber.
	LuckOfTheSea() int
	// OpenWater returns whether the bobber rests in open water.
	OpenWater() bool
}

// castBobber returns the bobber the user currently has cast out, if any.
func castBobber(tx *world.Tx, user User) (fishingBobber, bool) {
	for e := range tx.Entities() {
		if bobber, ok := e.(fishingBobber); ok && bobber.OwnedBy(user.H()) {
			return bobber, true
		}
	}
	return nil, false
}

// fishingLoot rolls a catch from the fish, junk or treasure table. Treasure is only reachable out of open
// water, and Luck of the Sea shifts the roll away from junk and towards treasure.
func fishingLoot(bobber fishingBobber) Stack {
	treasure, junk := 5.0, 10.0
	if !bobber.OpenWater() {
		// Treasure requires open water. Its share falls to the junk table instead.
		treasure, junk = 0, junk+5
	}
	if luck := bobber.LuckOfTheSea(); luck > 0 {
		t, j := LuckOfTheSeaChances(luck)
		treasure, junk = treasure+t, max(junk-j, 0)
	}
	switch roll := rand.Float64() * 100; {
	case roll < treasure:
		return weightedLoot(treasureLoot())
	case roll < treasure+junk:
		return weightedLoot(junkLoot())
	default:
		return weightedLoot(fishLoot())
	}
}

// LuckOfTheSeaChances returns the percentage points a Luck of the Sea level adds to the treasure chance and
// subtracts from the junk chance. It is overwritten by the enchantment package on registration.
var LuckOfTheSeaChances = func(int) (treasure, junk float64) { return 0, 0 }

// weightedLoot picks one entry out of the loot table passed, proportionally to the weights of its entries.
func weightedLoot(table []lootEntry) Stack {
	total := 0
	for _, e := range table {
		total += e.weight
	}
	roll := rand.IntN(total)
	for _, e := range table {
		if roll -= e.weight; roll < 0 {
			return e.stack
		}
	}
	return table[len(table)-1].stack
}

// lootEntry is a single possible catch along with its weight within its table.
type lootEntry struct {
	stack  Stack
	weight int
}

// fishLoot returns the fish that may be caught, with the weights they are caught at.
func fishLoot() []lootEntry {
	return []lootEntry{
		{NewStack(Cod{}, 1), 60},
		{NewStack(Salmon{}, 1), 25},
		{NewStack(Pufferfish{}, 1), 13},
		{NewStack(TropicalFish{}, 1), 2},
	}
}

// junkLoot returns the junk that may be caught, with the weights it is caught at.
func junkLoot() []lootEntry {
	entries := []lootEntry{
		{NewStack(Bone{}, 1), 10},
		{NewStack(Bowl{}, 1), 10},
		{NewStack(Leather{}, 1), 10},
		{NewStack(Boots{Tier: ArmourTierLeather{}}, 1), 10},
		{NewStack(RottenFlesh{}, 1), 10},
		{NewStack(Potion{Type: potion.Water()}, 1), 10},
		{NewStack(Stick{}, 1), 5},
		{NewStack(FishingRod{}, 1), 2},
		{NewStack(InkSac{}, 10), 2},
	}
	if lily, ok := world.ItemByName("minecraft:waterlily", 0); ok {
		entries = append(entries, lootEntry{NewStack(lily, 1), 17})
	}
	if str, ok := world.ItemByName("minecraft:string", 0); ok {
		entries = append(entries, lootEntry{NewStack(str, 1), 5})
	}
	return entries
}

// treasureLoot returns the treasure that may be caught, with the weights it is caught at.
func treasureLoot() []lootEntry {
	return []lootEntry{
		{NewStack(EnchantedBook{}, 1), 6},
		{NewStack(Bow{}, 1), 5},
		{NewStack(FishingRod{}, 1), 5},
		{NewStack(NautilusShell{}, 1), 5},
	}
}

// EncodeItem ...
func (FishingRod) EncodeItem() (name string, meta int16) {
	return "minecraft:fishing_rod", 0
}
