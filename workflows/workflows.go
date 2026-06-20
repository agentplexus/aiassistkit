// Package workflows provides workflow generation for AI coding assistants.
//
// Workflows define multi-phase spec creation processes that can be deployed
// across different AI coding assistants (Claude Code, Kiro, Cursor, etc.).
//
// The package follows the adapter pattern:
//   - core/ contains canonical workflow types
//   - claude/, kiro/, cursor/, etc. contain platform-specific adapters
//
// Usage:
//
//	// Load a spec-workflows repository
//	repo, err := core.LoadSourceRepo("/path/to/spec-workflows")
//
//	// Get a specific workflow
//	workflow, err := repo.GetWorkflow("aws-working-backwards/product")
//
//	// Generate for a specific platform
//	adapter, err := core.Get("claude")
//	err = adapter.Generate(workflow, "/path/to/project")
//
// Or use the convenience function:
//
//	cfg := &core.GenerateConfig{
//	    Workflow:  workflow,
//	    SourceRepo: repo,
//	    OutputDir: "/path/to/project",
//	    Platform:  "claude",
//	}
//	err := core.Generate(cfg)
//
// Supported platforms:
//   - claude: Claude Code (CLAUDE.md)
//   - kiro: Kiro CLI/IDE (.kiro/steering/)
//   - cursor: Cursor IDE (.cursor/rules/)
//   - amazonq: Amazon Q Developer (.amazonq/rules/)
//   - cline: Cline (.clinerules/)
//   - copilot: GitHub Copilot (.github/copilot-instructions.md)
package workflows

import (
	// Register all adapters
	_ "github.com/plexusone/assistantkit/workflows/amazonq"
	_ "github.com/plexusone/assistantkit/workflows/claude"
	_ "github.com/plexusone/assistantkit/workflows/cline"
	_ "github.com/plexusone/assistantkit/workflows/copilot"
	_ "github.com/plexusone/assistantkit/workflows/cursor"
	_ "github.com/plexusone/assistantkit/workflows/kiro"
)
