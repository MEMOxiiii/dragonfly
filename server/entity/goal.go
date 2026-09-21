package entity

import (
	"slices"
)

// Control is a bit mask of the parts of a Mob that a Goal takes over while it
// runs. Two goals that share a Control may never run at the same time.
type Control uint8

const (
	ControlMove Control = 1 << iota
	ControlLook
	ControlJump
	ControlTarget
)

// Goal is a single piece of behaviour of a Mob, such as strolling around or
// attacking a target. Goals of a lower priority value run first and stop goals
// of a higher value that use the same Control.
type Goal interface {
	Priority() int
	Controls() Control
	// Clone returns a copy of the Goal, so that every mob runs its own.
	Clone() Goal
	CanStart(m *Mob) bool
	CanContinue(m *Mob) bool
	Start(m *Mob)
	Stop(m *Mob)
	Tick(m *Mob)
}

// GoalBase implements the parts of Goal that most goals do not need to change.
type GoalBase struct {
	Prio int
}

func (g GoalBase) Priority() int { return g.Prio }
func (GoalBase) Start(*Mob)      {}
func (GoalBase) Stop(*Mob)       {}
func (GoalBase) Tick(*Mob)       {}

type goalSelector struct {
	goals   []Goal
	running []Goal
}

func newGoalSelector(goals []Goal) *goalSelector {
	slices.SortStableFunc(goals, func(a, b Goal) int { return a.Priority() - b.Priority() })
	return &goalSelector{goals: goals}
}

func (s *goalSelector) tick(m *Mob) {
	for i := 0; i < len(s.running); {
		if g := s.running[i]; !g.CanContinue(m) {
			g.Stop(m)
			s.running = slices.Delete(s.running, i, i+1)
			continue
		}
		i++
	}
	for _, g := range s.goals {
		if slices.Contains(s.running, g) || s.blocked(g) || !g.CanStart(m) {
			continue
		}
		s.preempt(m, g)
		g.Start(m)
		s.running = append(s.running, g)
	}
	for _, g := range s.running {
		g.Tick(m)
	}
}

// blocked reports whether a goal of an equal or higher priority is running with
// one of the controls the goal passed needs.
func (s *goalSelector) blocked(g Goal) bool {
	for _, running := range s.running {
		if running.Controls()&g.Controls() != 0 && running.Priority() <= g.Priority() {
			return true
		}
	}
	return false
}

// preempt stops the running goals of a lower priority that hold a control the
// goal passed needs.
func (s *goalSelector) preempt(m *Mob, g Goal) {
	for i := 0; i < len(s.running); {
		if running := s.running[i]; running.Controls()&g.Controls() != 0 {
			running.Stop(m)
			s.running = slices.Delete(s.running, i, i+1)
			continue
		}
		i++
	}
}

func (s *goalSelector) stop(m *Mob) {
	for _, g := range s.running {
		g.Stop(m)
	}
	s.running = s.running[:0]
}
