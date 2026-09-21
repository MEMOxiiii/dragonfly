package entity

import (
	"github.com/df-mc/dragonfly/server/internal/nbtconv"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/google/uuid"
)

// decodeMob reads the state of a mob from the NBT map passed into the data of
// an entity that the MobConfig of its type was applied to.
func decodeMob(m map[string]any, data *world.EntityData) {
	b, ok := data.Data.(*mobBehaviour)
	if !ok {
		return
	}
	if health := nbtconv.Float64(m, "Health"); health > 0 {
		b.health = NewHealthManager(health, b.conf.MaxHealth)
	}
	b.baby = nbtconv.Bool(m, "IsBaby")
	b.age = int(nbtconv.Int32(m, "Age"))
	b.love = int(nbtconv.Int32(m, "InLove"))
	b.breedCooldown = int(nbtconv.Int32(m, "BreedCooldown"))
	b.fireTicks = int(nbtconv.Int32(m, "Fire"))
	b.sheared = nbtconv.Bool(m, "Sheared")
	b.saddled = nbtconv.Bool(m, "Saddled")
	b.sitting = nbtconv.Bool(m, "Sitting")
	if colour, ok := m["Color"]; ok {
		b.colour = nbtconv.Uint8(map[string]any{"v": colour}, "v")
	}
	if variant, ok := m["Variant"]; ok {
		b.variant = nbtconv.Int32(map[string]any{"v": variant}, "v")
	}
	if id, err := uuid.Parse(nbtconv.String(m, "OwnerUUID")); err == nil {
		b.ownerID = id
	}
}

// encodeMob writes the state of a mob to an NBT map.
func encodeMob(data *world.EntityData) map[string]any {
	b, ok := data.Data.(*mobBehaviour)
	if !ok {
		return nil
	}
	m := map[string]any{
		"Health":        b.health.Health(),
		"IsBaby":        boolByte(b.baby),
		"Age":           int32(b.age),
		"InLove":        int32(b.love),
		"BreedCooldown": int32(b.breedCooldown),
		"Fire":          int32(b.fireTicks),
		"Color":         b.colour,
		"Variant":       b.variant,
	}
	if b.sheared {
		m["Sheared"] = boolByte(true)
	}
	if b.saddled {
		m["Saddled"] = boolByte(true)
	}
	if b.sitting {
		m["Sitting"] = boolByte(true)
	}
	if owner := b.ownerUUID(); owner != uuid.Nil {
		m["OwnerUUID"] = owner.String()
	}
	return m
}
