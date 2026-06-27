// Package core provides canonical types for composable capabilities.
// Capabilities are reusable building blocks that can be converted to
// both Skills (for human-guided use) and ValidationArea Checks (for
// autonomous QA validation).
package core

// Capability represents a composable capability that can be used by
// both skills and subagents. It serves as the single source of truth
// for validation/QA capabilities like linting, testing, and building.
type Capability struct {
	// Metadata
	Name        string `json:"name"`        // Unique identifier (e.g., "golangci-lint")
	Description string `json:"description"` // Brief description

	// Content
	Instructions string `json:"instructions"` // Full instructions/prompt for the capability

	// Execution
	Command     string   `json:"command,omitempty"`      // CLI command to execute
	CommandArgs []string `json:"command_args,omitempty"` // Default command arguments

	// Triggers (for skill invocation)
	Triggers []string `json:"triggers,omitempty"` // Keywords that invoke this capability

	// Dependencies
	Dependencies []string `json:"dependencies,omitempty"` // Required CLI tools

	// Error Handling Guidance
	ErrorHandling *ErrorHandlingGuide `json:"error_handling,omitempty"` // How to handle errors

	// Common Fixes
	CommonFixes []CommonFix `json:"common_fixes,omitempty"` // Known issues and fixes

	// Language/Platform
	Language string `json:"language,omitempty"` // Target language (e.g., "go", "typescript")
	Platform string `json:"platform,omitempty"` // Target platform (e.g., "all", "darwin")

	// Validation
	Required        bool   `json:"required"`                    // If true, failure blocks release
	SuccessPattern  string `json:"success_pattern,omitempty"`   // Regex pattern indicating success
	FailurePattern  string `json:"failure_pattern,omitempty"`   // Regex pattern indicating failure
	ExitCodeSuccess []int  `json:"exit_code_success,omitempty"` // Exit codes indicating success (default: [0])
}

// ErrorHandlingGuide provides instructions for handling errors.
type ErrorHandlingGuide struct {
	// Priority order for handling errors
	Priorities []ErrorPriority `json:"priorities"`

	// Reference packages for error handling
	References []PackageReference `json:"references,omitempty"`
}

// ErrorPriority represents a priority level for error handling.
type ErrorPriority struct {
	Level       int    `json:"level"`       // Priority level (1 = highest)
	Name        string `json:"name"`        // Priority name (e.g., "Panic", "Return Error")
	Description string `json:"description"` // When to use this priority
	Example     string `json:"example"`     // Code example
}

// PackageReference is a reference to a package for error handling.
type PackageReference struct {
	Package     string `json:"package"`     // Import path
	Description string `json:"description"` // What this package provides
	Usage       string `json:"usage"`       // Example usage
}

// CommonFix represents a common issue and its fix.
type CommonFix struct {
	Issue       string `json:"issue"`       // Lint rule or error pattern (e.g., "G306", "errcheck")
	Description string `json:"description"` // What the issue is
	Before      string `json:"before"`      // Code before fix
	After       string `json:"after"`       // Code after fix
}

// NewCapability creates a new Capability with the given name and description.
func NewCapability(name, description string) *Capability {
	return &Capability{
		Name:        name,
		Description: description,
	}
}

// WithCommand sets the CLI command for the capability.
func (c *Capability) WithCommand(cmd string, args ...string) *Capability {
	c.Command = cmd
	c.CommandArgs = args
	return c
}

// WithTriggers sets the trigger keywords.
func (c *Capability) WithTriggers(triggers ...string) *Capability {
	c.Triggers = triggers
	return c
}

// WithDependencies sets the required CLI tools.
func (c *Capability) WithDependencies(deps ...string) *Capability {
	c.Dependencies = deps
	return c
}

// WithLanguage sets the target language.
func (c *Capability) WithLanguage(lang string) *Capability {
	c.Language = lang
	return c
}

// WithInstructions sets the full instructions.
func (c *Capability) WithInstructions(instructions string) *Capability {
	c.Instructions = instructions
	return c
}

// AsRequired marks this capability as required (failure blocks release).
func (c *Capability) AsRequired() *Capability {
	c.Required = true
	return c
}

// AddCommonFix adds a common fix to the capability.
func (c *Capability) AddCommonFix(fix CommonFix) *Capability {
	c.CommonFixes = append(c.CommonFixes, fix)
	return c
}
