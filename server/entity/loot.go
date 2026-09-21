package entity

import (
	"math/rand/v2"

	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
)

// LootPool is a pool of a loot table. Every roll of a pool drops one of its
// entries, picked by weight.
type LootPool struct {
	Rolls   int
	Entries []LootEntry
}

// Always returns a LootPool that always drops between min and max of the item
// passed.
func Always(name string, minCount, maxCount int, opts ...LootOption) LootPool {
	return LootPool{Rolls: 1, Entries: []LootEntry{Entry(name, minCount, maxCount, opts...)}}
}

// OneOf returns a LootPool that drops one of the entries passed, picked by
// weight.
func OneOf(entries ...LootEntry) LootPool {
	return LootPool{Rolls: 1, Entries: entries}
}

// Entry returns a LootEntry for the item passed, dropping between min and max
// of it.
func Entry(name string, minCount, maxCount int, opts ...LootOption) LootEntry {
	e := LootEntry{Name: name, Min: minCount, Max: maxCount, Weight: 1}
	for _, opt := range opts {
		opt(&e)
	}
	return e
}

// LootOption changes a LootEntry, such as the item it drops when the mob dies
// on fire.
type LootOption func(e *LootEntry)

// Smelted returns a LootOption that drops the item passed instead if the mob
// dies on fire.
func Smelted(name string) LootOption {
	return func(e *LootEntry) { e.Smelted = name }
}

// Looting returns a LootOption that adds up to n items for every level of
// looting of the weapon that killed the mob.
func Looting(n int) LootOption {
	return func(e *LootEntry) { e.LootingMax = n }
}

// Weight returns a LootOption that sets how often an entry is picked from a
// pool relative to its other entries.
func Weight(n int) LootOption {
	return func(e *LootEntry) { e.Weight = n }
}

// LootEntry is an item a LootPool may drop.
type LootEntry struct {
	// Name is the encoded name of the item, such as 'minecraft:porkchop'.
	Name string
	// Meta is the metadata value of the item.
	Meta int16
	// Min and Max are the bounds of the amount dropped.
	Min, Max int
	// Weight is the chance of this entry being picked relative to the other
	// entries of the pool.
	Weight int
	// Smelted is the item dropped instead when the mob dies on fire.
	Smelted string
	// LootingMax is the extra amount a single level of looting may add.
	LootingMax int
}

// LootTable is the set of pools rolled when a mob dies.
type LootTable []LootPool

func (t LootTable) roll(onFire bool, looting int) []item.Stack {
	var stacks []item.Stack
	for _, pool := range t {
		for i := 0; i < max(pool.Rolls, 1); i++ {
			if e, ok := pool.pick(); ok {
				if st, ok := e.stack(onFire, looting); ok {
					stacks = append(stacks, st)
				}
			}
		}
	}
	return stacks
}

func (p LootPool) pick() (LootEntry, bool) {
	total := 0
	for _, e := range p.Entries {
		total += max(e.Weight, 1)
	}
	if total == 0 {
		return LootEntry{}, false
	}
	n := rand.IntN(total)
	for _, e := range p.Entries {
		if n -= max(e.Weight, 1); n < 0 {
			return e, true
		}
	}
	return LootEntry{}, false
}

func (e LootEntry) stack(onFire bool, looting int) (item.Stack, bool) {
	name := e.Name
	if onFire && e.Smelted != "" {
		name = e.Smelted
	}
	it, ok := world.ItemByName(name, e.Meta)
	if !ok {
		return item.Stack{}, false
	}
	count := e.Min
	if e.Max > e.Min {
		count += rand.IntN(e.Max - e.Min + 1)
	}
	if looting > 0 && e.LootingMax > 0 {
		count += rand.IntN(looting*e.LootingMax + 1)
	}
	if count <= 0 {
		return item.Stack{}, false
	}
	return item.NewStack(it, count), true
}
