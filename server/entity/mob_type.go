package entity

import (
	"github.com/df-mc/dragonfly/server/world"
)

// MobType is the world.EntityType of a mob that is driven by goals.
type MobType interface {
	world.EntityType
	world.EntityConfig
	// SpawnEggName returns the name of the spawn egg of the mob without the
	// '_spawn_egg' suffix, such as 'pig'. Mobs without a spawn egg return an
	// empty string.
	SpawnEggName() string
}

// NewMob creates a mob of the MobType passed at the position in the options.
func NewMob(t MobType, opts world.EntitySpawnOpts) *world.EntityHandle {
	return opts.New(t, t)
}
