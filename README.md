# AI Assist Kit

[![Build Status][build-status-svg]][build-status-url]
[![Lint Status][lint-status-svg]][lint-status-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![License][license-svg]][license-url]

AI Assist Kit is a Go library for managing configuration files across multiple AI coding assistants. It provides a unified interface for reading, writing, and converting between different tool-specific formats.

## Supported Tools

| Tool | MCP | Hooks | Context | Plugins | Commands | Skills | Agents |
|------|-----|-------|---------|---------|----------|--------|--------|
| Claude Code / Claude Desktop | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| Cursor IDE | ✅ | ✅ | — | — | — | — | — |
| Windsurf (Codeium) | ✅ | ✅ | — | — | — | — | — |
| VS Code / GitHub Copilot | ✅ | — | — | — | — | — | — |
| OpenAI Codex CLI | ✅ | — | — | — | ✅ | ✅ | ✅ |
| Cline | ✅ | — | — | — | — | — | — |
| Roo Code | ✅ | — | — | — | — | — | — |
| AWS Kiro CLI | ✅ | — | — | — | — | ✅ | — |
| Google Gemini CLI | — | — | — | ✅ | ✅ | — | ✅ |

## Configuration Types

| Type | Description | Status |
|------|-------------|--------|
| **MCP** | MCP server configurations | ✅ Available |
| **Hooks** | Automation/lifecycle callbacks | ✅ Available |
| **Context** | Project context (CONTEXT.json → CLAUDE.md) | ✅ Available |
| **Plugins** | Plugin/extension configurations | ✅ Available |
| **Commands** | Slash command definitions | ✅ Available |
| **Skills** | Reusable skill definitions | ✅ Available |
| **Agents** | AI assistant agent definitions | ✅ Available |
| **Teams** | Multi-agent team orchestration | ✅ Available |
| **Validation** | Configuration validators | ✅ Available |
| **Settings** | Permissions, sandbox, general settings | 🔜 Coming soon |
| **Rules** | Team rules, coding guidelines | 🔜 Coming soon |
| **Memory** | CLAUDE.md, .cursorrules, etc. | 🔜 Coming soon |

## Installation

```bash
go get github.com/agentplexus/aiassistkit
```

## MCP Configuration

The `mcp` subpackage provides adapters for MCP server configurations.

### Reading and Writing Configs

```go
package main

import (
    "log"

    "github.com/agentplexus/aiassistkit/mcp/claude"
    "github.com/agentplexus/aiassistkit/mcp/vscode"
)

func main() {
    // Read Claude config
    cfg, err := claude.ReadProjectConfig()
    if err != nil {
        log.Fatal(err)
    }

    // Write to VS Code format
    if err := vscode.WriteWorkspaceConfig(cfg); err != nil {
        log.Fatal(err)
    }
}
```

### Creating a New Config

```go
package main

import (
    "github.com/agentplexus/aiassistkit/mcp"
    "github.com/agentplexus/aiassistkit/mcp/claude"
    "github.com/agentplexus/aiassistkit/mcp/core"
)

func main() {
    cfg := mcp.NewConfig()

    // Add a stdio server
    cfg.AddServer("github", core.Server{
        Transport: core.TransportStdio,
        Command:   "npx",
        Args:      []string{"-y", "@modelcontextprotocol/server-github"},
        Env: map[string]string{
            "GITHUB_PERSONAL_ACCESS_TOKEN": "${GITHUB_TOKEN}",
        },
    })

    // Add an HTTP server
    cfg.AddServer("sentry", core.Server{
        Transport: core.TransportHTTP,
        URL:       "https://mcp.sentry.dev/mcp",
        Headers: map[string]string{
            "Authorization": "Bearer ${SENTRY_API_KEY}",
        },
    })

    // Write to Claude format
    claude.WriteProjectConfig(cfg)
}
```

### Converting Between Formats

```go
package main

import (
    "log"
    "os"

    "github.com/agentplexus/aiassistkit/mcp"
)

func main() {
    // Read Claude JSON
    data, _ := os.ReadFile(".mcp.json")

    // Convert to VS Code format
    vscodeData, err := mcp.Convert(data, "claude", "vscode")
    if err != nil {
        log.Fatal(err)
    }

    os.WriteFile(".vscode/mcp.json", vscodeData, 0644)
}
```

### Using Adapters Dynamically

```go
package main

import (
    "log"

    "github.com/agentplexus/aiassistkit/mcp"
)

func main() {
    // Get adapter by name
    adapter, ok := mcp.GetAdapter("claude")
    if !ok {
        log.Fatal("adapter not found")
    }

    // Read config
    cfg, err := adapter.ReadFile(".mcp.json")
    if err != nil {
        log.Fatal(err)
    }

    // Convert to another format
    codexAdapter, _ := mcp.GetAdapter("codex")
    codexAdapter.WriteFile(cfg, "~/.codex/config.toml")
}
```

## MCP Format Differences

### Claude (Reference Format)

Most tools follow Claude's format with `mcpServers` as the root key:

```json
{
  "mcpServers": {
    "server-name": {
      "command": "npx",
      "args": ["-y", "@example/mcp-server"],
      "env": {"API_KEY": "..."}
    }
  }
}
```

### VS Code

VS Code uses `servers` (not `mcpServers`) and supports `inputs` for secrets:

```json
{
  "inputs": [
    {"type": "promptString", "id": "api-key", "description": "API Key", "password": true}
  ],
  "servers": {
    "server-name": {
      "type": "stdio",
      "command": "npx",
      "args": ["-y", "@example/mcp-server"],
      "env": {"API_KEY": "${input:api-key}"}
    }
  }
}
```

### Windsurf

Windsurf uses `serverUrl` instead of `url` for HTTP servers:

```json
{
  "mcpServers": {
    "remote-server": {
      "serverUrl": "https://example.com/mcp"
    }
  }
}
```

### Codex (TOML)

Codex uses TOML format with additional timeout and tool control options:

```toml
[mcp_servers.github]
command = "npx"
args = ["-y", "@modelcontextprotocol/server-github"]
enabled_tools = ["list_repos", "create_issue"]
startup_timeout_sec = 30
tool_timeout_sec = 120
```

### AWS Kiro CLI

Kiro uses a format similar to Claude with support for both local and remote MCP servers. Environment variable substitution uses `${ENV_VAR}` syntax:

```json
{
  "mcpServers": {
    "github": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-github"],
      "env": {
        "GITHUB_TOKEN": "${GITHUB_TOKEN}"
      }
    },
    "remote-api": {
      "url": "https://api.example.com/mcp",
      "headers": {
        "Authorization": "Bearer ${API_TOKEN}"
      }
    },
    "disabled-server": {
      "command": "test",
      "disabled": true
    }
  }
}
```

**File locations:**
- Workspace: `.kiro/settings/mcp.json`
- User: `~/.kiro/settings/mcp.json`

## Hooks Configuration

The `hooks` subpackage provides adapters for automation/lifecycle hooks that execute at defined stages of the agent loop.

### Creating Hooks

```go
package main

import (
    "github.com/agentplexus/aiassistkit/hooks"
    "github.com/agentplexus/aiassistkit/hooks/claude"
)

func main() {
    cfg := hooks.NewConfig()

    // Add a command hook that runs before shell commands
    cfg.AddHookWithMatcher(hooks.BeforeCommand, "Bash",
        hooks.NewCommandHook("echo 'Running command...'"))

    // Add a hook for file writes
    cfg.AddHook(hooks.BeforeFileWrite,
        hooks.NewCommandHook("./scripts/validate-write.sh"))

    // Write to Claude format
    claude.WriteProjectConfig(cfg)
}
```

### Converting Between Formats

```go
package main

import (
    "log"
    "os"

    "github.com/agentplexus/aiassistkit/hooks"
)

func main() {
    // Read Claude hooks JSON
    data, _ := os.ReadFile(".claude/settings.json")

    // Convert to Cursor format
    cursorData, err := hooks.Convert(data, "claude", "cursor")
    if err != nil {
        log.Fatal(err)
    }

    os.WriteFile(".cursor/hooks.json", cursorData, 0644)
}
```

### Supported Events

| Event | Claude | Cursor | Windsurf | Description |
|-------|--------|--------|----------|-------------|
| `before_file_read` | ✅ | ✅ | ✅ | Before reading a file |
| `after_file_read` | ✅ | ✅ | ✅ | After reading a file |
| `before_file_write` | ✅ | ✅ | ✅ | Before writing a file |
| `after_file_write` | ✅ | ✅ | ✅ | After writing a file |
| `before_command` | ✅ | ✅ | ✅ | Before shell command execution |
| `after_command` | ✅ | ✅ | ✅ | After shell command execution |
| `before_mcp` | ✅ | ✅ | ✅ | Before MCP tool call |
| `after_mcp` | ✅ | ✅ | ✅ | After MCP tool call |
| `before_prompt` | ✅ | — | ✅ | Before user prompt processing |
| `on_stop` | ✅ | ✅ | — | When agent stops |
| `on_session_start` | ✅ | — | — | When session starts |
| `on_session_end` | ✅ | — | — | When session ends |
| `after_response` | — | ✅ | — | After AI response (Cursor-only) |
| `after_thought` | — | ✅ | — | After AI thought (Cursor-only) |
| `on_permission` | ✅ | — | — | Permission request (Claude-only) |

### Hook Types

- **Command hooks**: Execute shell commands
- **Prompt hooks**: Run AI prompts (Claude-only)

## Project Structure

```
aiassistkit/
├── aiassistkit.go          # Umbrella package
├── agents/                 # Agent definitions
│   ├── agentkit/           # AWS AgentKit adapter
│   ├── awsagentcore/       # AWS CDK TypeScript generator
│   ├── claude/             # Claude Code adapter
│   ├── codex/              # Codex adapter
│   ├── core/               # Canonical types
│   ├── gemini/             # Gemini adapter
│   └── kiro/               # AWS Kiro CLI adapter
├── cmd/
│   └── genagents/          # Multi-platform agent generator CLI
├── commands/               # Slash command definitions
│   ├── claude/             # Claude adapter
│   ├── codex/              # Codex adapter
│   ├── core/               # Canonical types
│   └── gemini/             # Gemini adapter
├── context/                # Project context (CONTEXT.json → CLAUDE.md)
│   ├── claude/             # CLAUDE.md converter
│   └── core/               # Canonical types
├── hooks/                  # Lifecycle hooks
│   ├── claude/             # Claude adapter
│   ├── core/               # Canonical types
│   ├── cursor/             # Cursor adapter
│   └── windsurf/           # Windsurf adapter
├── mcp/                    # MCP server configurations
│   ├── claude/             # Claude adapter
│   ├── cline/              # Cline adapter
│   ├── codex/              # Codex adapter (TOML)
│   ├── core/               # Canonical types
│   ├── cursor/             # Cursor adapter
│   ├── kiro/               # AWS Kiro CLI adapter
│   ├── roo/                # Roo Code adapter
│   ├── vscode/             # VS Code adapter
│   └── windsurf/           # Windsurf adapter
├── plugins/                # Plugin/extension configurations
│   ├── claude/             # Claude adapter
│   ├── core/               # Canonical types
│   └── gemini/             # Gemini adapter
├── publish/                # Marketplace publishing
│   ├── claude/             # Claude marketplace adapter
│   ├── core/               # Publishing interfaces
│   └── github/             # GitHub API client
├── skills/                 # Reusable skill definitions
│   ├── claude/             # Claude adapter
│   ├── codex/              # Codex adapter
│   └── core/               # Canonical types
├── teams/                  # Multi-agent orchestration
│   └── core/               # Team types and workflows
└── validation/             # Configuration validators
    ├── claude/             # Claude validator
    ├── codex/              # Codex validator
    ├── core/               # Validation interfaces
    └── gemini/             # Gemini validator
```

## Related Projects

AI Assist Kit is part of the AgentPlexus family of Go modules for building AI agents:

- **AI Assist Kit** - AI coding assistant configuration management
- **OmniVault** - Unified secrets management
- **OmniLLM** - Multi-provider LLM abstraction
- **OmniSerp** - Search engine abstraction
- **OmniObserve** - LLM observability abstraction

## License

MIT License - see [LICENSE](LICENSE) for details.

 [build-status-svg]: https://github.com/agentplexus/aiassistkit/actions/workflows/ci.yaml/badge.svg?branch=main
 [build-status-url]: https://github.com/agentplexus/aiassistkit/actions/workflows/ci.yaml
 [lint-status-svg]: https://github.com/agentplexus/aiassistkit/actions/workflows/lint.yaml/badge.svg?branch=main
 [lint-status-url]: https://github.com/agentplexus/aiassistkit/actions/workflows/lint.yaml
 [goreport-svg]: https://goreportcard.com/badge/github.com/agentplexus/aiassistkit
 [goreport-url]: https://goreportcard.com/report/github.com/agentplexus/aiassistkit
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/agentplexus/aiassistkit
 [docs-godoc-url]: https://pkg.go.dev/github.com/agentplexus/aiassistkit
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/agentplexus/aiassistkit/blob/master/LICENSE
 [used-by-svg]: https://sourcegraph.com/github.com/agentplexus/aiassistkit/-/badge.svg
 [used-by-url]: https://sourcegraph.com/github.com/agentplexus/aiassistkit?badge
