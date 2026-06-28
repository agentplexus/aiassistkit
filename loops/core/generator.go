package core

import (
	"fmt"
	"strings"

	multiagentspec "github.com/plexusone/multi-agent-spec/sdk/go"
)

// LoopAgentConfig describes how to generate agents from a loop.
type LoopAgentConfig struct {
	// IncludeLoopInstructions adds loop-specific instructions to agents.
	IncludeLoopInstructions bool

	// ValidatorSuffix is appended to validator agent instructions.
	ValidatorSuffix string

	// ActorSuffix is appended to actor agent instructions.
	ActorSuffix string
}

// DefaultLoopAgentConfig returns the default configuration.
func DefaultLoopAgentConfig() *LoopAgentConfig {
	return &LoopAgentConfig{
		IncludeLoopInstructions: true,
	}
}

// GenerateLoopInstructions generates loop-aware instructions for an agent.
// This is appended to the agent's existing instructions when it participates in a loop.
func GenerateLoopInstructions(loop *Loop, role string) string {
	var sb strings.Builder

	sb.WriteString("\n\n## Loop Participation\n\n")
	sb.WriteString(fmt.Sprintf("This agent participates in the **%s** loop (%s pattern).\n\n", loop.Name, loop.Type))

	if loop.Description != "" {
		sb.WriteString(fmt.Sprintf("**Purpose:** %s\n\n", loop.Description))
	}

	switch role {
	case "validator":
		sb.WriteString(generateValidatorInstructions(loop))
	case "actor":
		sb.WriteString(generateActorInstructions(loop))
	}

	sb.WriteString(fmt.Sprintf("\n**Max Attempts:** %d\n", loop.EffectiveMaxAttempts()))
	sb.WriteString(fmt.Sprintf("**Escalation Policy:** %s\n", loop.EffectiveEscalation()))

	if loop.SuccessCriteria != "" {
		sb.WriteString(fmt.Sprintf("\n**Success Criteria:**\n%s\n", loop.SuccessCriteria))
	}

	return sb.String()
}

func generateValidatorInstructions(loop *Loop) string {
	var sb strings.Builder

	sb.WriteString("### Validator Role\n\n")
	sb.WriteString("As the **validator** in this loop, your responsibility is to:\n\n")
	sb.WriteString("1. Run all validation checks\n")
	sb.WriteString("2. Report GO/NO-GO status for each check\n")
	sb.WriteString("3. Provide detailed findings for any failures\n")
	sb.WriteString("4. Do NOT modify any files (read-only)\n\n")

	if len(loop.Checks) > 0 {
		sb.WriteString("### Validation Checks\n\n")
		sb.WriteString("| ID | Type | Required | Description |\n")
		sb.WriteString("|----|----|-------|-------------|\n")

		for _, check := range loop.Checks {
			required := "Yes"
			if !check.IsRequired() {
				required := "No"
				_ = required
			}
			desc := check.Description
			if desc == "" {
				desc = check.ID
			}
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n",
				check.ID, check.Type, required, desc))
		}
		sb.WriteString("\n")

		// Detail commands/patterns
		sb.WriteString("### Check Details\n\n")
		for _, check := range loop.Checks {
			sb.WriteString(fmt.Sprintf("**%s**", check.ID))
			if check.Description != "" {
				sb.WriteString(fmt.Sprintf(": %s", check.Description))
			}
			sb.WriteString("\n")

			switch check.Type {
			case CheckTypeCommand:
				sb.WriteString(fmt.Sprintf("- Command: `%s`\n", check.Command))
			case CheckTypePattern:
				sb.WriteString(fmt.Sprintf("- Pattern: `%s`\n", check.Pattern))
				if check.Files != "" {
					sb.WriteString(fmt.Sprintf("- Files: `%s`\n", check.Files))
				}
			case CheckTypeFile:
				sb.WriteString(fmt.Sprintf("- File: `%s`\n", check.File))
			}

			if check.Expected != "" {
				sb.WriteString(fmt.Sprintf("- Expected: %s\n", check.Expected))
			}
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

func generateActorInstructions(loop *Loop) string {
	var sb strings.Builder

	sb.WriteString("### Actor Role\n\n")

	if loop.Type == LoopVEAL {
		sb.WriteString("As the **actor** in this VEAL loop, your responsibility is to:\n\n")
		sb.WriteString("1. Receive findings from the validator\n")
		sb.WriteString("2. Fix identified issues\n")
		sb.WriteString("3. Apply corrections systematically\n")
		sb.WriteString("4. Report what actions were taken\n\n")

		if len(loop.Checks) > 0 {
			sb.WriteString("### Issues to Address\n\n")
			sb.WriteString("The validator may report issues for these checks:\n\n")
			for _, check := range loop.Checks {
				sb.WriteString(fmt.Sprintf("- **%s**: %s\n", check.ID, check.Description))
			}
			sb.WriteString("\n")
		}
	} else if loop.Type == LoopREAL {
		sb.WriteString("As the **actor** in this REAL loop, your responsibility is to:\n\n")
		sb.WriteString("1. Work toward the mission goal\n")
		sb.WriteString("2. Report progress after each iteration\n")
		sb.WriteString("3. Determine when the mission is complete\n")
		sb.WriteString("4. Request escalation if stuck\n\n")

		if loop.Mission != "" {
			sb.WriteString("### Mission\n\n")
			sb.WriteString(loop.Mission)
			sb.WriteString("\n\n")
		}
	}

	return sb.String()
}

// EnrichAgentWithLoop adds loop participation instructions to an agent.
func EnrichAgentWithLoop(agent *multiagentspec.Agent, loop *Loop, role string) *multiagentspec.Agent {
	enriched := *agent // shallow copy

	// Add loop instructions
	loopInstructions := GenerateLoopInstructions(loop, role)
	enriched.Instructions = agent.Instructions + loopInstructions

	return &enriched
}

// AgentsFromLoop returns the agent names involved in a loop.
func AgentsFromLoop(loop *Loop) []string {
	return loop.Agents()
}

// GenerateCoordinatorInstructions generates instructions for a loop coordinator.
// This is used when a coordinator agent orchestrates multiple loops.
func GenerateCoordinatorInstructions(loops []*Loop) string {
	var sb strings.Builder

	sb.WriteString("\n\n## Loop Orchestration\n\n")
	sb.WriteString("You coordinate the following autonomous loops:\n\n")

	// Group by type
	var vealLoops, realLoops []*Loop
	for _, loop := range loops {
		if loop.Type == LoopVEAL {
			vealLoops = append(vealLoops, loop)
		} else {
			realLoops = append(realLoops, loop)
		}
	}

	if len(vealLoops) > 0 {
		sb.WriteString("### VEAL Loops (State-Driven Validation)\n\n")
		sb.WriteString("| Loop | Validator | Actor | Max Attempts |\n")
		sb.WriteString("|------|-----------|-------|-------------|\n")
		for _, loop := range vealLoops {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %d |\n",
				loop.Name, loop.Validator, loop.Actor, loop.EffectiveMaxAttempts()))
		}
		sb.WriteString("\n")
	}

	if len(realLoops) > 0 {
		sb.WriteString("### REAL Loops (Mission-Driven)\n\n")
		sb.WriteString("| Loop | Actor | Mission | Max Attempts |\n")
		sb.WriteString("|------|-------|---------|-------------|\n")
		for _, loop := range realLoops {
			mission := loop.Mission
			if len(mission) > 40 {
				mission = mission[:37] + "..."
			}
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %d |\n",
				loop.Name, loop.Actor, mission, loop.EffectiveMaxAttempts()))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("### Loop Execution Protocol\n\n")
	sb.WriteString("For each **VEAL loop**:\n")
	sb.WriteString("1. Invoke validator agent to check state\n")
	sb.WriteString("2. If GO → proceed to next loop\n")
	sb.WriteString("3. If NO-GO → invoke actor agent to fix\n")
	sb.WriteString("4. Re-invoke validator to verify fixes\n")
	sb.WriteString("5. Repeat until GO or max attempts reached\n")
	sb.WriteString("6. On max attempts → apply escalation policy\n\n")

	sb.WriteString("For each **REAL loop**:\n")
	sb.WriteString("1. Invoke actor agent with mission\n")
	sb.WriteString("2. Actor reports progress/completion\n")
	sb.WriteString("3. If complete → proceed to next loop\n")
	sb.WriteString("4. If not complete → continue loop\n")
	sb.WriteString("5. On max attempts → apply escalation policy\n")

	return sb.String()
}
