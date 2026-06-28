# Loops

The `loops` package provides REAL/VEAL autonomous loop patterns for bounded, self-correcting agent execution.

## Overview

Loops are autonomous execution patterns that run without human intervention until they complete or escalate. They are inspired by the REPL (Read-Eval-Print-Loop) pattern from interactive programming.

## Loop Types

| Type | Category | Pattern | Use Case |
|------|----------|---------|----------|
| **REAL** | Mission-driven | Read → Eval → Act → Loop | Long-running tasks: roadmaps, migrations, feature implementations |
| **VEAL** | State-driven | Validate → Eval → Act → Loop | Fix loops: QA validation, lint fixes, doc updates |

### REAL — Read Eval Act Loop

Mission-driven loops for open-ended tasks where the goal is defined by a mission statement.

```
R - Read (mission, requirements, goals)
E - Eval (assess situation against goals)
A - Act (take action toward goals)
L - Loop (until mission complete or escalate)
```

**Characteristics:**

- Open-ended goals
- Requirements-driven
- May involve discovery/exploration
- Single actor agent

### VEAL — Validate Eval Act Loop

State-driven loops that converge toward a valid state through validation checks.

```
V - Validate (check current state)
E - Eval (compare to expected → GO/NO-GO)
A - Act (correct state if NO-GO)
L - Loop (until valid or escalate)
```

**Characteristics:**

- Convergent (state approaches valid)
- Idempotent checks
- Bounded attempts
- Dual-agent pattern (validator + actor)

## Bounded Execution

Loops run autonomously but are bounded to prevent infinite execution:

| Field | Default | Description |
|-------|---------|-------------|
| `max_attempts` | 3 | Maximum loop iterations before escalation |
| `escalation` | `human` | What to do when max reached |
| `success_criteria` | - | Description of what success looks like |

### Escalation Policies

| Policy | Behavior |
|--------|----------|
| `human` | Stop and request human intervention (default) |
| `abort` | Stop and fail the workflow |
| `continue` | Proceed despite unresolved issues |
| `fallback` | Invoke a fallback agent |

### Choosing max_attempts

| Loop Type | Typical Range | Rationale |
|-----------|---------------|-----------|
| **VEAL** (QA fix) | 3-5 | If code can't pass QA in 3 attempts, there's a deeper issue |
| **VEAL** (docs fix) | 3-5 | Documentation updates should converge quickly |
| **REAL** (feature) | 5-10 | Features may need more iterations |
| **REAL** (roadmap) | 50-100 | Long-running work over many sessions |

## Examples

### VEAL Loop: QA Fix

Short-lived validation loop — if code can't pass QA in 3 attempts, escalate:

```yaml
name: qa-fix
type: VEAL
description: QA validation and fix loop
validator: qa           # Read-only validator agent
actor: code-fixer       # Agent that fixes issues
max_attempts: 3         # Try 3 times max
escalation: human       # Ask human if stuck
checks:
  - id: build
    type: command
    command: go build ./...
    required: true
  - id: lint
    type: command
    command: golangci-lint run
    required: true
  - id: tests
    type: command
    command: go test ./...
    required: true
```

### REAL Loop: Roadmap Implementation

Long-running mission loop — work through a roadmap over many iterations:

```yaml
name: roadmap-impl
type: REAL
description: Implement quarterly roadmap
actor: developer
mission: |
  Implement all items in Q3 roadmap.
  Prioritize by business value.
  Create PRs for each feature.
max_attempts: 100       # Many iterations allowed
escalation: human       # Escalate if truly stuck
success_criteria: |
  All roadmap items completed or explicitly deferred.
  Each feature has passing tests and documentation.
```

## Loop-Aware Agent Generation

When generating agents, loops automatically inject participation instructions into agents referenced as validators or actors.

```go
import "github.com/plexusone/assistantkit/generate"

result, err := generate.Generate("specs", "local", ".")
fmt.Printf("Loaded %d loops\n", result.LoopCount)
// Agents referenced as validator/actor get loop instructions injected
```

The injected instructions include:

- Loop participation section with role description
- Validation checks table (for validators)
- Issue addressing guidance (for actors)
- Max attempts and escalation policy

## Working with Loops Programmatically

### Loading Loops

```go
import "github.com/plexusone/assistantkit/loops"

// Load all loops from specs directory
loopSet, err := loops.LoadLoopSet("specs/loops")
if err != nil {
    log.Fatal(err)
}

// Get specific loop by name
qaLoop, ok := loopSet.Get("qa-fix")

// Filter by type
vealLoops := loopSet.VEALLoops()  // All validation loops
realLoops := loopSet.REALLoops()  // All mission loops

// Get all loop names
names := loopSet.Names()
```

### Generating Instructions

```go
// Generate role-specific instructions for an agent
validatorInstructions := loops.GenerateLoopInstructions(qaLoop, "validator")
actorInstructions := loops.GenerateLoopInstructions(qaLoop, "actor")

// Enrich an agent with loop participation
enrichedAgent := loops.EnrichAgentWithLoop(agent, qaLoop, "validator")

// Generate coordinator instructions for orchestrating multiple loops
coordInstructions := loops.GenerateCoordinatorInstructions(loopSet.All())
```

### Creating Loops Programmatically

```go
// Create a VEAL loop
qaLoop := loops.NewVEALLoop("qa-fix", "qa", "code-fixer").
    WithDescription("QA validation and fix loop").
    WithMaxAttempts(3).
    AddCheck(loops.LoopCheck{
        ID:      "build",
        Type:    loops.CheckTypeCommand,
        Command: "go build ./...",
    })

// Create a REAL loop
featureLoop := loops.NewREALLoop(
    "feature-impl",
    "developer",
    "Implement user authentication with OAuth2",
).WithMaxAttempts(10)
```

## Directory Structure

Loops are defined in the `specs/loops/` directory:

```
specs/
├── agents/
│   ├── qa.md
│   └── code-fixer.md
├── loops/
│   ├── qa-fix.yaml
│   ├── docs-fix.yaml
│   └── security-fix.yaml
└── deployments/
    └── local.json
```

## Using Loops in Workflows

Loops can be referenced as steps in team workflows:

```json
{
  "workflow": {
    "type": "chain",
    "steps": [
      { "name": "validate-scope", "agent": "pm" },
      { "name": "qa-fix", "loop": "qa-fix" },
      { "name": "docs-fix", "loop": "docs-fix" },
      { "name": "tag-release", "agent": "release" }
    ]
  }
}
```

## Best Practices

1. **Keep checks idempotent** — Running the same check multiple times should produce the same result

2. **Use specific checks** — Prefer multiple small checks over one large check

3. **Set reasonable max_attempts** — 3 is usually sufficient for VEAL; more may indicate a deeper issue

4. **Separate concerns** — Validator agents should be read-only; actor agents should fix issues

5. **Document success criteria** — Make it clear what "done" looks like

6. **Use fallback agents** — For complex issues that may need escalation to a specialist

## Related

- [Loop Schema Reference](https://plexusone.github.io/multi-agent-spec/schemas/loop/)
- [Loop Pattern Guide](https://plexusone.github.io/multi-agent-spec/guides/loops/)
