// Package capabilities provides composable capabilities for AI assistants.
//
// Capabilities are reusable building blocks that can be converted to both:
//   - Skills: For human-guided use via slash commands or triggers
//   - ValidationArea Checks: For autonomous QA validation by agents
//
// This package solves the problem of defining validation logic once and using
// it in multiple contexts. For example, the golangci-lint capability can be:
//   - Invoked as a /lint skill by a human for interactive guidance
//   - Used by a QA agent as part of automated release validation
//
// # Architecture
//
// The capabilities package follows the same adapter pattern as the rest of
// assistantkit:
//
//	Capability (canonical type)
//	    ↓ ToSkill()
//	Skill (for human-guided use)
//	    ↓ ToCheck()
//	Check (for validation agents)
//
// # Usage
//
// Define a capability once:
//
//	golangciLint := gocaps.GolangciLint()
//
// Convert to a skill for human-guided use:
//
//	skill := golangciLint.ToSkill()
//	// Use with Claude Code skill adapter to generate SKILL.md
//
// Convert to a check for QA validation:
//
//	check := golangciLint.ToCheck()
//	// Add to ValidationArea for automated checking
//
// Use capability sets for grouped operations:
//
//	qaChecks := gocaps.QACapabilities()
//
//	// Generate all skills
//	skills := qaChecks.ToSkills()
//
//	// Create a QA validation area
//	va := qaChecks.ToValidationArea(
//	    "All checks pass",
//	    gocaps.QAValidationInstructions(),
//	)
//
// # Package Structure
//
//	capabilities/
//	├── core/           # Canonical Capability type and converters
//	│   ├── capability.go
//	│   └── converters.go
//	├── go/             # Go-specific capabilities
//	│   ├── golangci_lint.go
//	│   ├── go_test.go
//	│   ├── go_build.go
//	│   └── registry.go
//	└── doc.go          # This file
//
// # Adding New Capabilities
//
// To add a new capability:
//
//  1. Create a function that returns *core.Capability
//  2. Set name, description, command, triggers, dependencies
//  3. Add instructions for AI assistants
//  4. Add common fixes if applicable
//  5. Add to the appropriate registry (All, QACapabilities, etc.)
//
// Example:
//
//	func MyNewCheck() *core.Capability {
//	    c := core.NewCapability("my-check", "Description of my check")
//	    c.WithCommand("my-tool", "run")
//	    c.WithTriggers("check", "my-check")
//	    c.WithDependencies("my-tool")
//	    c.WithInstructions(`# My Check Skill\n...`)
//	    c.AsRequired()
//	    return c
//	}
package capabilities
