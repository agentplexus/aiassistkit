// Package core provides core loop types and interfaces for REAL/VEAL loop patterns.
package core

import (
	multiagentspec "github.com/plexusone/multi-agent-spec/sdk/go"
)

// Loop is the canonical loop type from multi-agent-spec.
type Loop = multiagentspec.Loop

// LoopType is the loop execution pattern.
type LoopType = multiagentspec.LoopType

// LoopCheck is a validation check within a loop.
type LoopCheck = multiagentspec.LoopCheck

// CheckType is how a check is executed.
type CheckType = multiagentspec.CheckType

// EscalationPolicy defines what to do when max attempts are reached.
type EscalationPolicy = multiagentspec.EscalationPolicy

// Loop type constants.
const (
	LoopREAL = multiagentspec.LoopREAL
	LoopVEAL = multiagentspec.LoopVEAL
)

// Check type constants.
const (
	CheckTypeCommand = multiagentspec.CheckTypeCommand
	CheckTypePattern = multiagentspec.CheckTypePattern
	CheckTypeFile    = multiagentspec.CheckTypeFile
	CheckTypeManual  = multiagentspec.CheckTypeManual
)

// Escalation policy constants.
const (
	EscalationHuman    = multiagentspec.EscalationHuman
	EscalationAbort    = multiagentspec.EscalationAbort
	EscalationContinue = multiagentspec.EscalationContinue
	EscalationFallback = multiagentspec.EscalationFallback
)

// NewLoop creates a new loop with the given name and type.
func NewLoop(name string, loopType LoopType) *Loop {
	return multiagentspec.NewLoop(name, loopType)
}

// NewVEALLoop creates a new VEAL loop with validator and actor agents.
func NewVEALLoop(name, validator, actor string) *Loop {
	return multiagentspec.NewVEALLoop(name, validator, actor)
}

// NewREALLoop creates a new REAL loop with actor agent and mission.
func NewREALLoop(name, actor, mission string) *Loop {
	return multiagentspec.NewREALLoop(name, actor, mission)
}

// LoadLoopFromFile loads a loop from a JSON or YAML file.
func LoadLoopFromFile(path string) (*Loop, error) {
	return multiagentspec.LoadLoopFromFile(path)
}

// LoadLoopsFromDir loads all loops from a directory.
func LoadLoopsFromDir(dir string) ([]*Loop, error) {
	return multiagentspec.LoadLoopsFromDir(dir)
}

// LoopSet is a collection of loops indexed by name.
type LoopSet struct {
	loops map[string]*Loop
}

// NewLoopSet creates a new empty LoopSet.
func NewLoopSet() *LoopSet {
	return &LoopSet{
		loops: make(map[string]*Loop),
	}
}

// Add adds a loop to the set.
func (ls *LoopSet) Add(loop *Loop) {
	ls.loops[loop.Name] = loop
}

// Get returns a loop by name.
func (ls *LoopSet) Get(name string) (*Loop, bool) {
	loop, ok := ls.loops[name]
	return loop, ok
}

// All returns all loops in the set.
func (ls *LoopSet) All() []*Loop {
	result := make([]*Loop, 0, len(ls.loops))
	for _, loop := range ls.loops {
		result = append(result, loop)
	}
	return result
}

// Names returns all loop names.
func (ls *LoopSet) Names() []string {
	names := make([]string, 0, len(ls.loops))
	for name := range ls.loops {
		names = append(names, name)
	}
	return names
}

// Len returns the number of loops.
func (ls *LoopSet) Len() int {
	return len(ls.loops)
}

// VEALLoops returns all VEAL loops.
func (ls *LoopSet) VEALLoops() []*Loop {
	var result []*Loop
	for _, loop := range ls.loops {
		if loop.Type == LoopVEAL {
			result = append(result, loop)
		}
	}
	return result
}

// REALLoops returns all REAL loops.
func (ls *LoopSet) REALLoops() []*Loop {
	var result []*Loop
	for _, loop := range ls.loops {
		if loop.Type == LoopREAL {
			result = append(result, loop)
		}
	}
	return result
}

// LoadLoopSet loads a LoopSet from a directory.
func LoadLoopSet(dir string) (*LoopSet, error) {
	loops, err := LoadLoopsFromDir(dir)
	if err != nil {
		return nil, err
	}

	set := NewLoopSet()
	for _, loop := range loops {
		set.Add(loop)
	}

	return set, nil
}
