package entity

import "github.com/df-mc/dragonfly/server/world"

// Float returns a FloatGoal with the priority passed.
func Float(priority int) FloatGoal {
	return FloatGoal{GoalBase{priority}}
}

// RandomStroll returns a RandomStrollGoal with the priority and speed
// multiplier passed.
func RandomStroll(priority int, speedMultiplier float64) *RandomStrollGoal {
	return &RandomStrollGoal{GoalBase: GoalBase{priority}, SpeedMultiplier: speedMultiplier}
}

// LookAtPlayer returns a LookAtPlayerGoal that looks at players up to the
// distance passed, for a random time in seconds between minLook and maxLook.
func LookAtPlayer(priority int, distance, minLook, maxLook float64) *LookAtPlayerGoal {
	return &LookAtPlayerGoal{GoalBase: GoalBase{priority}, Distance: distance, MinLookTime: minLook, MaxLookTime: maxLook}
}

// RandomLookAround returns a RandomLookAroundGoal with the priority passed.
func RandomLookAround(priority int) *RandomLookAroundGoal {
	return &RandomLookAroundGoal{GoalBase: GoalBase{priority}}
}

// Panic returns a PanicGoal with the priority and speed multiplier passed.
func Panic(priority int, speedMultiplier float64) *PanicGoal {
	return &PanicGoal{GoalBase: GoalBase{priority}, SpeedMultiplier: speedMultiplier}
}

// HurtByTarget returns a HurtByTargetGoal with the priority passed.
func HurtByTarget(priority int) HurtByTargetGoal {
	return HurtByTargetGoal{GoalBase{priority}}
}

// NearestAttackableTarget returns a NearestAttackableTargetGoal that searches
// for the entities passed within the radius passed and forgets them once they
// are further away than maxDistance.
func NearestAttackableTarget(priority int, radius, maxDistance float64, reselect bool, targets ...string) *NearestAttackableTargetGoal {
	return &NearestAttackableTargetGoal{GoalBase: GoalBase{priority}, WithinRadius: radius,
		MaxDistance: maxDistance, ReselectTargets: reselect, Targets: targets}
}

// MeleeAttack returns a MeleeAttackGoal that deals the damage passed.
func MeleeAttack(priority int, damage, speedMultiplier float64) *MeleeAttackGoal {
	return &MeleeAttackGoal{GoalBase: GoalBase{priority}, Damage: damage, SpeedMultiplier: speedMultiplier}
}

// RangedAttack returns a RangedAttackGoal that shoots a projectile every
// minInterval to maxInterval seconds at targets within the radius passed.
func RangedAttack(priority int, radius, minInterval, maxInterval float64, shoot func(m *Mob, target world.Entity)) *RangedAttackGoal {
	return &RangedAttackGoal{GoalBase: GoalBase{priority}, AttackRadius: radius,
		MinInterval: minInterval, MaxInterval: maxInterval, Shoot: shoot}
}

// RangedBurstAttack returns a RangedAttackGoal that shoots several projectiles
// in a row, with burstInterval seconds between them.
func RangedBurstAttack(priority int, radius, minInterval, maxInterval float64, shots int, burstInterval float64, shoot func(m *Mob, target world.Entity)) *RangedAttackGoal {
	g := RangedAttack(priority, radius, minInterval, maxInterval, shoot)
	g.BurstShots, g.BurstInterval = shots, burstInterval
	return g
}

// Swell returns a SwellGoal that explodes with the power passed after the fuse
// in ticks has run out.
func Swell(priority int, startDistance, stopDistance float64, fuse int, power float64) *SwellGoal {
	return &SwellGoal{GoalBase: GoalBase{priority}, StartDistance: startDistance,
		StopDistance: stopDistance, FuseTicks: fuse, Power: power}
}

// FleeSun returns a FleeSunGoal with the priority and speed multiplier passed.
func FleeSun(priority int, speedMultiplier float64) *FleeSunGoal {
	return &FleeSunGoal{GoalBase: GoalBase{priority}, SpeedMultiplier: speedMultiplier}
}

// AvoidMobType returns an AvoidMobTypeGoal that runs away from the entities
// passed once they are within maxDistance.
func AvoidMobType(priority int, maxDistance, walkSpeed, sprintSpeed float64, avoided ...string) *AvoidMobTypeGoal {
	return &AvoidMobTypeGoal{GoalBase: GoalBase{priority}, MaxDistance: maxDistance,
		WalkSpeedMultiplier: walkSpeed, SprintSpeedMultiplier: sprintSpeed, Avoided: avoided}
}

// RandomSwim returns a RandomSwimGoal that swims to a position up to xzDist
// blocks away and yDist blocks up or down, every interval ticks.
func RandomSwim(priority, xzDist, yDist, interval int, speedMultiplier float64) *RandomSwimGoal {
	return &RandomSwimGoal{GoalBase: GoalBase{priority}, Range: xzDist, YRange: yDist,
		Interval: interval, SpeedMultiplier: speedMultiplier}
}

// RandomFly returns a RandomFlyGoal that flies to a position up to xzDist
// blocks away, staying between minHeight and maxHeight above the ground.
func RandomFly(priority, xzDist, yDist int, yOffset, minHeight, maxHeight, speedMultiplier float64) *RandomFlyGoal {
	return &RandomFlyGoal{GoalBase: GoalBase{priority}, Range: xzDist, YRange: yDist, YOffset: yOffset,
		MinHeight: minHeight, MaxHeight: maxHeight, SpeedMultiplier: speedMultiplier}
}

// MoveToWater returns a MoveToWaterGoal that searches water within the range
// passed.
func MoveToWater(priority, searchRange int, speedMultiplier float64) *MoveToWaterGoal {
	return &MoveToWaterGoal{GoalBase: GoalBase{priority}, SearchRange: searchRange, SpeedMultiplier: speedMultiplier}
}

// SquidSwim returns a SquidSwimGoal with the priority passed.
func SquidSwim(priority int) *SquidSwimGoal {
	return &SquidSwimGoal{GoalBase: GoalBase{priority}}
}

// Breed returns a BreedGoal with the priority and speed multiplier passed.
func Breed(priority int, speedMultiplier float64) *BreedGoal {
	return &BreedGoal{GoalBase: GoalBase{priority}, SpeedMultiplier: speedMultiplier}
}

// Tempt returns a TemptGoal that makes the mob follow players holding one of
// the items passed.
func Tempt(priority int, speedMultiplier, distance float64, items ...string) *TemptGoal {
	return &TemptGoal{GoalBase: GoalBase{priority}, SpeedMultiplier: speedMultiplier, Distance: distance, Items: items}
}

// FollowParent returns a FollowParentGoal with the priority and speed
// multiplier passed.
func FollowParent(priority int, speedMultiplier float64) *FollowParentGoal {
	return &FollowParentGoal{GoalBase: GoalBase{priority}, SpeedMultiplier: speedMultiplier}
}
