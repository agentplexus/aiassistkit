// Package core provides canonical types for AI assistant plugin/extension definitions.
package core

// Plugin represents a canonical plugin/extension definition that can be
// converted to tool-specific formats (Claude, Gemini, etc.).
type Plugin struct {
	// Metadata
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Author      string `json:"author,omitempty"`
	License     string `json:"license,omitempty"`
	Repository  string `json:"repository,omitempty"`
	Homepage    string `json:"homepage,omitempty"`

	// Components - paths to spec files
	Commands string `json:"commands,omitempty"` // Directory containing command specs
	Skills   string `json:"skills,omitempty"`   // Directory containing skill specs
	Agents   string `json:"agents,omitempty"`   // Directory containing agent specs
	Hooks    string `json:"hooks,omitempty"`    // Path to hooks spec file

	// Context file content
	Context string `json:"context,omitempty"` // System prompt / context content

	// Dependencies
	Dependencies []Dependency `json:"dependencies,omitempty"`

	// MCP Servers (used by Gemini extensions)
	MCPServers map[string]MCPServer `json:"mcp_servers,omitempty"`
}

// Dependency represents a required or optional dependency.
type Dependency struct {
	Name     string `json:"name"`
	Command  string `json:"command,omitempty"`  // CLI command to check availability
	Optional bool   `json:"optional,omitempty"` // If true, missing dependency is a warning
}

// MCPServer represents an MCP server configuration.
type MCPServer struct {
	Command string            `json:"command"`
	Args    []string          `json:"args,omitempty"`
	Cwd     string            `json:"cwd,omitempty"`
	Env     map[string]string `json:"env,omitempty"`
}

// NewPlugin creates a new Plugin with the given name, version, and description.
func NewPlugin(name, version, description string) *Plugin {
	return &Plugin{
		Name:        name,
		Version:     version,
		Description: description,
	}
}

// AddDependency adds a required dependency to the plugin.
func (p *Plugin) AddDependency(name, command string) {
	p.Dependencies = append(p.Dependencies, Dependency{
		Name:    name,
		Command: command,
	})
}

// AddOptionalDependency adds an optional dependency to the plugin.
func (p *Plugin) AddOptionalDependency(name, command string) {
	p.Dependencies = append(p.Dependencies, Dependency{
		Name:     name,
		Command:  command,
		Optional: true,
	})
}

// AddMCPServer adds an MCP server configuration to the plugin.
func (p *Plugin) AddMCPServer(name string, server MCPServer) {
	if p.MCPServers == nil {
		p.MCPServers = make(map[string]MCPServer)
	}
	p.MCPServers[name] = server
}
