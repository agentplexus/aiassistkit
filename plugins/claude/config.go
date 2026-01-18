package claude

import (
	"github.com/agentplexus/assistantkit/plugins/core"
)

// ClaudePlugin represents the Claude Code plugin.json format.
// See: https://docs.anthropic.com/en/docs/claude-code/plugins
type ClaudePlugin struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`

	// Optional metadata
	Author     string `json:"author,omitempty"`
	License    string `json:"license,omitempty"`
	Repository string `json:"repository,omitempty"`
	Homepage   string `json:"homepage,omitempty"`

	// Component paths (relative to plugin root)
	Commands string `json:"commands,omitempty"` // e.g., "./commands/"
	Skills   string `json:"skills,omitempty"`   // e.g., "./skills/"
	Agents   string `json:"agents,omitempty"`   // e.g., "./agents/"
	Hooks    string `json:"hooks,omitempty"`    // e.g., "./hooks/hooks.json"
}

// ToCanonical converts ClaudePlugin to canonical Plugin.
func (cp *ClaudePlugin) ToCanonical() *core.Plugin {
	return &core.Plugin{
		Name:        cp.Name,
		Version:     cp.Version,
		Description: cp.Description,
		Author:      cp.Author,
		License:     cp.License,
		Repository:  cp.Repository,
		Homepage:    cp.Homepage,
		Commands:    cp.Commands,
		Skills:      cp.Skills,
		Agents:      cp.Agents,
		Hooks:       cp.Hooks,
	}
}

// FromCanonical creates a ClaudePlugin from canonical Plugin.
func FromCanonical(p *core.Plugin) *ClaudePlugin {
	cp := &ClaudePlugin{
		Name:        p.Name,
		Version:     p.Version,
		Description: p.Description,
		Author:      p.Author,
		License:     p.License,
		Repository:  p.Repository,
		Homepage:    p.Homepage,
	}

	// Set default paths if components are specified
	if p.Commands != "" {
		cp.Commands = "./commands/"
	}
	if p.Skills != "" {
		cp.Skills = "./skills/"
	}
	if p.Agents != "" {
		cp.Agents = "./agents/"
	}
	if p.Hooks != "" {
		cp.Hooks = "./hooks/hooks.json"
	}

	return cp
}
