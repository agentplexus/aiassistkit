---
marp: true
theme: vibeminds
paginate: true
header: "AI Assist Kit: Unified AI Assistant Configuration"
---

# AI Assist Kit

## Unified Configuration Management for AI Coding Assistants

A Go library for reading, writing, and converting between tool-specific configuration formats.

---

# The Problem

Each AI coding assistant has its own configuration format:

| Tool | MCP Config | Format |
|------|------------|--------|
| Claude Code | `.mcp.json` | JSON |
| Cursor | `~/.cursor/mcp.json` | JSON |
| VS Code | `settings.json` | JSON |
| Windsurf | `~/.codeium/windsurf/mcp_config.json` | JSON |
| Codex | `codex.toml` | TOML |
| Kiro | `.kiro/settings/mcp.json` | JSON |

**Result:** N tools = N different configs to maintain

---

# The Solution

## Adapter Pattern with Canonical Model

```
Tool A Format ──► Adapter A ──► Canonical Model ──► Adapter B ──► Tool B Format
```

- **N adapters** instead of **N² direct conversions**
- Single source of truth
- Automatic format translation

---

# Supported Tools

| Tool | Description |
|------|-------------|
| **Claude** | Claude Code / Claude Desktop |
| **Cursor** | Cursor IDE |
| **Windsurf** | Windsurf (Codeium) |
| **VS Code** | VS Code / GitHub Copilot |
| **Codex** | OpenAI Codex CLI |
| **Cline** | Cline VS Code extension |
| **Roo** | Roo Code VS Code extension |
| **Kiro** | AWS Kiro CLI |
| **Gemini** | Google Gemini CLI |

---

# Configuration Types

## Available Now

| Type | Description |
|------|-------------|
| **MCP** | Model Context Protocol server configurations |
| **Agents** | AI assistant agent definitions |
| **Commands** | Slash command definitions |
| **Plugins** | Plugin/extension configurations |
| **Skills** | Reusable skill definitions |
| **Teams** | Multi-agent team configurations |
| **Validation** | Configuration validators |

---

# Architecture Overview

```
aiassistkit/
├── mcp/           # MCP server configurations
├── agents/        # Agent definitions
├── commands/      # Slash commands
├── plugins/       # Plugin configs
├── skills/        # Skill definitions
├── teams/         # Team configurations
├── validation/    # Config validators
└── context/       # Project context converters
```

---

# Package Structure

Each package follows the adapter pattern:

```
mcp/
├── core/          # Canonical types + registry
├── claude/        # Claude adapter
├── cursor/        # Cursor adapter
├── vscode/        # VS Code adapter
├── codex/         # Codex adapter (TOML)
├── kiro/          # Kiro adapter
└── ...            # Other adapters
```

---

# The Adapter Interface

```go
type Adapter interface {
    // Identity
    Name() string
    DefaultPaths() []string

    // Parsing & Marshaling
    Parse(data []byte) (*Config, error)
    Marshal(cfg *Config) ([]byte, error)

    // File I/O
    ReadFile(path string) (*Config, error)
    WriteFile(cfg *Config, path string) error
}
```

---

# Canonical MCP Config

```go
type Config struct {
    // Map of server names to configurations
    Servers map[string]Server `json:"servers"`

    // Input variables for sensitive data (VS Code)
    Inputs []InputVariable `json:"inputs,omitempty"`
}

type Server struct {
    Transport Transport         // stdio, http, sse
    Command   string            // For stdio
    Args      []string
    Env       map[string]string
    URL       string            // For http/sse
    Headers   map[string]string
}
```

---

# Usage Example: Reading Config

```go
import (
    "github.com/grokify/aiassistkit/mcp/claude"
)

// Read Claude's project config
cfg, err := claude.ReadProjectConfig()
if err != nil {
    log.Fatal(err)
}

// Access servers
for name, server := range cfg.Servers {
    fmt.Printf("Server: %s, Command: %s\n", name, server.Command)
}
```

---

# Usage Example: Converting Formats

```go
import (
    "github.com/grokify/aiassistkit/mcp/core"
    _ "github.com/grokify/aiassistkit/mcp/claude"
    _ "github.com/grokify/aiassistkit/mcp/vscode"
)

// Convert Claude config to VS Code format
vsCodeData, err := core.Convert(claudeJSON, "claude", "vscode")
if err != nil {
    log.Fatal(err)
}
```

Adapters auto-register via `init()` functions.

---

# Usage Example: Writing Config

```go
import (
    "github.com/grokify/aiassistkit/mcp/core"
    "github.com/grokify/aiassistkit/mcp/vscode"
)

// Create a new config
cfg := core.NewConfig()
cfg.AddServer("my-mcp", core.Server{
    Transport: core.TransportStdio,
    Command:   "npx",
    Args:      []string{"-y", "@my/mcp-server"},
})

// Write to VS Code format
vscode.NewAdapter().WriteFile(cfg, ".vscode/settings.json")
```

---

# Agents Package

Define AI assistant agents with tools and skills:

```go
import "github.com/grokify/aiassistkit/agents"

agent := agents.NewAgent("release-coordinator",
    "Orchestrates software releases")
agent.SetModel("sonnet")
agent.AddTools("Read", "Write", "Bash", "Glob", "Grep")
agent.AddSkills("version-analysis", "commit-classification")

// Write to Claude format
claudeAdapter, _ := agents.GetAdapter("claude")
claudeAdapter.WriteFile(agent, "./agents/release-coordinator.md")
```

---

# Commands Package

Define slash commands across tools:

```go
import "github.com/grokify/aiassistkit/commands"

cmd := commands.NewCommand("deploy", "Deploy to production")
cmd.AddRequiredArgument("environment", "Target environment")
cmd.AddOptionalArgument("--dry-run", "Preview without changes")

// Write to Claude format (Markdown)
claudeAdapter, _ := commands.GetAdapter("claude")
claudeAdapter.WriteFile(cmd, "./commands/deploy.md")

// Write to Gemini format (TOML)
geminiAdapter, _ := commands.GetAdapter("gemini")
geminiAdapter.WriteFile(cmd, "./commands/deploy.toml")
```

---

# Skills Package

Define reusable skills:

```go
import "github.com/grokify/aiassistkit/skills"

skill := skills.NewSkill("code-review", "Reviews code for issues")
skill.SetPrompt("Review this code for bugs and improvements...")
skill.AddDependency("Read")
skill.AddDependency("Grep")

// Write to Claude format
claudeAdapter, _ := skills.GetAdapter("claude")
claudeAdapter.WriteSkillDir(skill, "./skills/")
```

---

# Key Design Decisions

## Tri-state Booleans

```go
// *bool allows: true, false, or nil (default)
type Package struct {
    Public *bool `json:"public,omitempty"`
}

func (p *Package) IsPublic() bool {
    if p.Public == nil {
        return true  // Default to true
    }
    return *p.Public
}
```

---

# Key Design Decisions

## Error Handling

```go
// Custom errors implement Unwrap() for error chains
type ParseError struct {
    Format string
    Path   string
    Err    error
}

func (e *ParseError) Unwrap() error {
    return e.Err
}

// Usage with errors.Is/As
if errors.Is(err, os.ErrNotExist) {
    // Handle missing file
}
```

---

# Testing Strategy

## Patterns Used

- **Table-driven tests** with subtests
- **Round-trip tests**: marshal → parse → compare
- **Adapter conversion tests** between formats
- **Event mapping validation** tests

```bash
# Run all tests
go test ./...

# With coverage
go test ./... -cover

# Verbose
go test ./... -v
```

---

# Project Structure

Each adapter package follows a consistent pattern:

```
mcp/claude/
├── adapter.go       # Adapter implementation
├── config.go        # Tool-specific types
└── adapter_test.go  # Tests
```

**Conventions:**
- `Parse()` / `Marshal()` work with `[]byte`
- `ReadFile()` / `WriteFile()` work with paths
- File mode `0600` for security

---

# Version History

## v0.1.0
- 6 new packages: agents, commands, plugins, skills, teams, validation
- Expanded tool support with Gemini adapters
- Full documentation (PRD, TRD, ROADMAP)
- MCP configuration support
- 8 tool adapters

---

# Getting Started

```bash
# Install
go get github.com/grokify/aiassistkit

# Import adapters you need
import (
    "github.com/grokify/aiassistkit/mcp/core"
    _ "github.com/grokify/aiassistkit/mcp/claude"
    _ "github.com/grokify/aiassistkit/mcp/cursor"
)
```

---

<!-- _class: lead -->

# Thank You

**GitHub:** github.com/grokify/aiassistkit

**License:** MIT

Questions?
