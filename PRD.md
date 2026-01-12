# AI Assist Kit - Product Requirements Document

## Overview

AI Assist Kit is a Go library that provides a unified interface for managing configuration files across multiple AI coding assistants. It solves the fragmentation problem where developers must maintain separate configurations for each tool they use.

## Problem Statement

The AI coding assistant ecosystem is fragmented:

- **Claude Code** uses `.mcp.json`, `.claude-plugin/plugin.json`, `CLAUDE.md`
- **Gemini CLI** uses `gemini-extension.json`, `GEMINI.md`, TOML commands
- **OpenAI Codex** uses `config.toml`, `AGENTS.md`, Markdown prompts/skills
- **Cursor IDE** uses `mcp.json`, `.cursorrules`
- **VS Code Copilot** uses `mcp.json`, `copilot-instructions.md`
- **Windsurf** uses `mcp_config.json`

Developers who use multiple tools must:

1. Maintain duplicate configurations in different formats
2. Manually sync changes across tools
3. Learn each tool's specific syntax and conventions
4. Risk configurations drifting out of sync

## Solution

AI Assist Kit provides:

1. **Canonical Types**: A single source of truth in JSON format
2. **Adapter Pattern**: Tool-specific converters for each AI assistant
3. **Format Conversion**: Automatic transformation between formats
4. **Schema Validation**: JSON Schema for validation and IDE support

## Target Users

1. **Individual Developers**: Using multiple AI coding assistants
2. **Teams**: Standardizing AI tool configurations across projects
3. **Tool Authors**: Building plugins/extensions for multiple platforms
4. **Enterprise**: Managing AI assistant configurations at scale

## Core Features

### v0.1.0 - Initial Release (Completed)

| Feature | Description |
|---------|-------------|
| MCP Configurations | Manage MCP server configs across 8 tools |
| Context | Project context (CONTEXT.json) to CLAUDE.md |
| Plugin Manifests | Canonical plugin definition with Claude/Gemini adapters |
| Commands/Prompts | Command definitions with Claude/Gemini/Codex adapters |
| Skills | Skill definitions with Claude/Codex adapters |
| Agents | Agent definitions with Claude/Codex/Gemini/Kiro adapters |
| Teams | Multi-agent team orchestration types |
| Validation | Configuration validators with Claude/Codex/Gemini adapters |
| JSON Schemas | Validation schemas for all configuration types |

### v0.2.0 - Settings & Rules (Planned)

| Feature | Description |
|---------|-------------|
| Settings | Permission and sandbox configurations |
| Rules | Coding guidelines (.cursorrules, copilot-instructions) |
| Context Converters | GEMINI.md and AGENTS.md generation |

### v1.0.0 - General Availability (Planned)

| Feature | Description |
|---------|-------------|
| CLI Tool | `aiassistkit generate/validate/convert` commands |
| Full Documentation | Complete API docs and examples |
| Stable API | Backward-compatible API guarantees |

## Supported Tools

| Tool | MCP | Context | Plugins | Commands | Skills | Agents | Validation |
|------|-----|---------|---------|----------|--------|--------|------------|
| Claude Code | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| Cursor IDE | ✅ | — | — | — | — | — | — |
| Windsurf (Codeium) | ✅ | — | — | — | — | — | — |
| VS Code / GitHub Copilot | ✅ | — | — | — | — | — | — |
| OpenAI Codex CLI | ✅ | — | — | ✅ | ✅ | ✅ | ✅ |
| Cline | ✅ | — | — | — | — | — | — |
| Roo Code | ✅ | — | — | — | — | — | — |
| AWS Kiro CLI | ✅ | — | — | — | — | ✅ | — |
| Google Gemini CLI | — | — | ✅ | ✅ | — | ✅ | ✅ |

## User Workflows

### Define Once, Generate Everywhere

```bash
# Define canonical plugin spec
plugins/spec/plugin.json
plugins/spec/commands/release.json
plugins/spec/skills/version-analysis.json

# Generate for all platforms
aiassistkit generate plugins/spec/ --output plugins/

# Result:
plugins/claude/.claude-plugin/plugin.json
plugins/claude/commands/release.md
plugins/gemini/gemini-extension.json
plugins/gemini/commands/release.toml
plugins/codex/prompts/release.md
```

### Convert Between Formats

```bash
# Convert Claude MCP config to VS Code format
aiassistkit convert .mcp.json --from claude --to vscode -o .vscode/mcp.json
```

### Validate Configurations

```bash
# Validate plugin spec against JSON Schema
aiassistkit validate plugins/spec/plugin.json --schema plugin

# Validate all specs in directory
aiassistkit validate plugins/spec/ --all
```

## Success Metrics

1. **Adoption**: Number of projects using aiassistkit
2. **Tool Coverage**: Number of AI assistants supported
3. **Configuration Types**: Number of config types supported
4. **Conversion Accuracy**: 100% round-trip fidelity for conversions

## Non-Goals

- Runtime integration with AI assistants
- AI assistant functionality (prompting, completion)
- Cloud-based configuration management
- Configuration encryption or secrets management

## Dependencies

- Go 1.23+
- `github.com/pelletier/go-toml/v2` (for Codex TOML support)

## Related Projects

- [Release Agent](https://github.com/grokify/release-agent) - Uses aiassistkit for plugin generation
- [Structured Changelog](https://github.com/grokify/structured-changelog) - Changelog management
- [Structured Roadmap](https://github.com/grokify/structured-roadmap) - Roadmap management
