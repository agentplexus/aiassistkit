// Package loops provides REAL/VEAL loop pattern support for multi-agent systems.
//
// Loops are autonomous, bounded execution patterns that run without human
// intervention until they complete or escalate.
//
// REAL (Read Eval Act Loop) - Mission-driven loops for open-ended tasks
// VEAL (Validate Eval Act Loop) - State-driven validation loops
//
// Example usage:
//
//	// Load loops from directory
//	loopSet, err := loops.LoadLoopSet("specs/loops")
//
//	// Get a specific loop
//	qaLoop, _ := loopSet.Get("qa-fix")
//
//	// Generate loop instructions for an agent
//	instructions := loops.GenerateLoopInstructions(qaLoop, "validator")
//
//	// Enrich an agent with loop participation
//	enrichedAgent := loops.EnrichAgentWithLoop(agent, qaLoop, "actor")
package loops

import (
	"github.com/plexusone/assistantkit/loops/core"
	multiagentspec "github.com/plexusone/multi-agent-spec/sdk/go"
)

// Re-export core types
type (
	// Loop is the canonical loop type.
	Loop = core.Loop

	// LoopType is the loop execution pattern (REAL or VEAL).
	LoopType = core.LoopType

	// LoopCheck is a validation check within a loop.
	LoopCheck = core.LoopCheck

	// CheckType is how a check is executed.
	CheckType = core.CheckType

	// EscalationPolicy defines what to do when max attempts are reached.
	EscalationPolicy = core.EscalationPolicy

	// LoopSet is a collection of loops indexed by name.
	LoopSet = core.LoopSet

	// LoopAgentConfig describes how to generate agents from a loop.
	LoopAgentConfig = core.LoopAgentConfig
)

// Loop type constants.
const (
	LoopREAL = core.LoopREAL
	LoopVEAL = core.LoopVEAL
)

// Check type constants.
const (
	CheckTypeCommand = core.CheckTypeCommand
	CheckTypePattern = core.CheckTypePattern
	CheckTypeFile    = core.CheckTypeFile
	CheckTypeManual  = core.CheckTypeManual
)

// Escalation policy constants.
const (
	EscalationHuman    = core.EscalationHuman
	EscalationAbort    = core.EscalationAbort
	EscalationContinue = core.EscalationContinue
	EscalationFallback = core.EscalationFallback
)

// NewLoop creates a new loop with the given name and type.
func NewLoop(name string, loopType LoopType) *Loop {
	return core.NewLoop(name, loopType)
}

// NewVEALLoop creates a new VEAL loop with validator and actor agents.
func NewVEALLoop(name, validator, actor string) *Loop {
	return core.NewVEALLoop(name, validator, actor)
}

// NewREALLoop creates a new REAL loop with actor agent and mission.
func NewREALLoop(name, actor, mission string) *Loop {
	return core.NewREALLoop(name, actor, mission)
}

// NewLoopSet creates a new empty LoopSet.
func NewLoopSet() *LoopSet {
	return core.NewLoopSet()
}

// LoadLoopFromFile loads a loop from a JSON or YAML file.
func LoadLoopFromFile(path string) (*Loop, error) {
	return core.LoadLoopFromFile(path)
}

// LoadLoopsFromDir loads all loops from a directory.
func LoadLoopsFromDir(dir string) ([]*Loop, error) {
	return core.LoadLoopsFromDir(dir)
}

// LoadLoopSet loads a LoopSet from a directory.
func LoadLoopSet(dir string) (*LoopSet, error) {
	return core.LoadLoopSet(dir)
}

// DefaultLoopAgentConfig returns the default configuration.
func DefaultLoopAgentConfig() *LoopAgentConfig {
	return core.DefaultLoopAgentConfig()
}

// GenerateLoopInstructions generates loop-aware instructions for an agent.
func GenerateLoopInstructions(loop *Loop, role string) string {
	return core.GenerateLoopInstructions(loop, role)
}

// EnrichAgentWithLoop adds loop participation instructions to an agent.
func EnrichAgentWithLoop(agent *multiagentspec.Agent, loop *Loop, role string) *multiagentspec.Agent {
	return core.EnrichAgentWithLoop(agent, loop, role)
}

// AgentsFromLoop returns the agent names involved in a loop.
func AgentsFromLoop(loop *Loop) []string {
	return core.AgentsFromLoop(loop)
}

// GenerateCoordinatorInstructions generates instructions for a loop coordinator.
func GenerateCoordinatorInstructions(loops []*Loop) string {
	return core.GenerateCoordinatorInstructions(loops)
}
