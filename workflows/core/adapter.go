// Package core provides the canonical workflow adapter interface.
package core

import (
	"sync"
)

// Adapter converts between canonical Workflow and platform-specific formats.
type Adapter interface {
	// Name returns the adapter identifier (e.g., "claude", "kiro")
	Name() string

	// EntryPointFile returns the entry point filename for this platform
	// (e.g., "CLAUDE.md" for Claude Code, "core-workflow.md" for Kiro steering)
	EntryPointFile() string

	// EntryPointDir returns the directory for the entry point
	// (e.g., "" for root, ".kiro/steering" for Kiro)
	EntryPointDir() string

	// RuleDetailsDir returns the directory for rule details
	// (e.g., ".spec-workflows" for Claude, ".kiro/spec-workflow-details" for Kiro)
	RuleDetailsDir() string

	// Generate creates platform-specific files from a Workflow
	Generate(workflow *Workflow, outputDir string) error

	// GenerateEntryPoint generates just the entry point file content
	GenerateEntryPoint(workflow *Workflow) ([]byte, error)
}

// registry holds all registered adapters.
var (
	registry   = make(map[string]Adapter)
	registryMu sync.RWMutex
)

// Register adds an adapter to the registry.
func Register(adapter Adapter) {
	registryMu.Lock()
	defer registryMu.Unlock()
	registry[adapter.Name()] = adapter
}

// Get returns the adapter with the given name.
func Get(name string) (Adapter, error) {
	registryMu.RLock()
	defer registryMu.RUnlock()

	adapter, ok := registry[name]
	if !ok {
		return nil, &AdapterNotFoundError{Name: name}
	}
	return adapter, nil
}

// List returns all registered adapter names.
func List() []string {
	registryMu.RLock()
	defer registryMu.RUnlock()

	var names []string
	for name := range registry {
		names = append(names, name)
	}
	return names
}

// GenerateConfig specifies how to generate workflow files.
type GenerateConfig struct {
	// Workflow is the source workflow to generate from
	Workflow *Workflow

	// SourceRepo is the spec-workflows repository
	SourceRepo *SourceRepo

	// OutputDir is the target project directory
	OutputDir string

	// Platform is the target platform (e.g., "claude", "kiro")
	Platform string

	// Extension is an optional extension to include
	Extension string

	// CopyRuleDetails copies rule-details to output (vs symlink/reference)
	CopyRuleDetails bool

	// CopyTemplates copies templates to output
	CopyTemplates bool

	// CopyRubrics copies rubrics to output
	CopyRubrics bool
}

// Generate generates platform-specific workflow files.
func Generate(cfg *GenerateConfig) error {
	adapter, err := Get(cfg.Platform)
	if err != nil {
		return err
	}

	return adapter.Generate(cfg.Workflow, cfg.OutputDir)
}
