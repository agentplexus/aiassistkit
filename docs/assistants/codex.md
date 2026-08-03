# OpenAI Codex CLI

[OpenAI Codex CLI](https://github.com/openai/codex) is OpenAI's command-line coding assistant.

## Plugin Structure

```
my-plugin/
├── agents.yaml          # Agent definitions (optional)
└── commands/            # Slash commands
    └── build.md
```

## Commands

Commands are Markdown files with YAML frontmatter:

```markdown
---
name: build
description: Build the project
---

Build the project using the appropriate build system.
Detect the project type and run the correct command.
```

## Agents (Optional)

If you need agent-like behavior, define in `agents.yaml`:

```yaml
agents:
  - name: code-reviewer
    description: Reviews code for quality
    instructions: |
      You review code for best practices...
```

## Installation

Copy your plugin to the Codex plugins directory:

```bash
cp -r my-plugin ~/.codex/plugins/
```

## Hooks

AssistantKit includes a Codex `SessionStart` hook command that reads Codex hook
JSON from stdin and injects commit-message guidance with the active model:

```bash
assistantkit codex hooks generated-with
```

Add it to `~/.codex/config.toml`:

```toml
commit_attribution = "Co-authored-by: OpenAI Codex <codex@openai.com>"

[[hooks.SessionStart]]
matcher = "startup|resume|clear|compact"

[[hooks.SessionStart.hooks]]
type = "command"
command = "/Users/johnwang/go/bin/assistantkit codex hooks generated-with"
timeout = 10
statusMessage = "Loading Codex commit attribution context"
```

The hook emits a `Generated-with: OpenAI Codex <model>` instruction using the
`model` field from Codex hook input.

## Available Tools

| Tool | Description |
|------|-------------|
| read | Read file contents |
| write | Write files |
| shell | Execute shell commands |
| search | Search files |

## MCP Server Configuration

Codex CLI supports MCP servers in configuration:

```json
{
  "mcpServers": {
    "filesystem": {
      "command": "npx",
      "args": ["-y", "@anthropic-ai/mcp-server-filesystem"]
    }
  }
}
```

## Converting from Canonical

```go
import (
    "github.com/plexusone/assistantkit/commands/core"
    "github.com/plexusone/assistantkit/commands/codex"
)

canonical := core.Command{
    Name:        "build",
    Description: "Build the project",
    Prompt:      "Build the project...",
}

codexCmd := codex.FromCanonical(canonical)
```

## Limitations

- **No Skills**: Codex CLI does not have a dedicated skills system
- **Limited Agent Support**: Basic agent configuration only
- **Simpler Structure**: Fewer configuration options than Claude Code

## Tool Mapping

| Canonical | Codex |
|-----------|-------|
| Read | read |
| Write | write |
| Edit | edit |
| Bash | shell |
| Glob | search |
| Grep | search |
