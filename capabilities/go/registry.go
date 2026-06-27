package gocaps

import (
	"github.com/plexusone/assistantkit/capabilities/core"
)

// All returns all Go capabilities as a slice.
func All() []*core.Capability {
	return []*core.Capability{
		GolangciLint(),
		GoTest(),
		GoTestCoverage(),
		GoTestRace(),
		GoBuild(),
		GoModTidy(),
		GoVet(),
		GoFmt(),
	}
}

// QACapabilities returns capabilities relevant for QA validation.
// These are the core checks that should pass before release.
func QACapabilities() *core.CapabilitySet {
	cs := core.NewCapabilitySet("go-qa", "Go Quality Assurance checks")
	cs.Language = "go"

	// Add QA-relevant capabilities in execution order
	cs.Add(GoBuild())
	cs.Add(GoModTidy())
	cs.Add(GoFmt())
	cs.Add(GoVet())
	cs.Add(GolangciLint())
	cs.Add(GoTest())
	cs.Add(GoTestRace())

	return cs
}

// LintCapabilities returns capabilities for linting and static analysis.
func LintCapabilities() *core.CapabilitySet {
	cs := core.NewCapabilitySet("go-lint", "Go linting and static analysis")
	cs.Language = "go"

	cs.Add(GoFmt())
	cs.Add(GoVet())
	cs.Add(GolangciLint())

	return cs
}

// TestCapabilities returns capabilities for testing.
func TestCapabilities() *core.CapabilitySet {
	cs := core.NewCapabilitySet("go-test", "Go testing capabilities")
	cs.Language = "go"

	cs.Add(GoTest())
	cs.Add(GoTestCoverage())
	cs.Add(GoTestRace())

	return cs
}

// BuildCapabilities returns capabilities for building.
func BuildCapabilities() *core.CapabilitySet {
	cs := core.NewCapabilitySet("go-build", "Go build capabilities")
	cs.Language = "go"

	cs.Add(GoBuild())
	cs.Add(GoModTidy())

	return cs
}

// QAValidationInstructions returns the instructions for a Go QA validation agent.
func QAValidationInstructions() string {
	return `# Go QA Validation Agent

You are a Quality Assurance validation agent for Go projects. Your role is to
verify that code meets quality standards before release.

## Validation Checks

Run the following checks in order:

1. **Build**: ` + "`go build ./...`" + ` - Verify code compiles
2. **Module Tidy**: ` + "`go mod tidy -diff`" + ` - Ensure dependencies are clean
3. **Format**: ` + "`gofmt -l -d .`" + ` - Check code formatting
4. **Vet**: ` + "`go vet ./...`" + ` - Run static analysis
5. **Lint**: ` + "`golangci-lint run ./...`" + ` - Run comprehensive linting
6. **Test**: ` + "`go test -v ./...`" + ` - Run all tests
7. **Race**: ` + "`go test -v -race ./...`" + ` - Check for data races

## Status Determination

- **GO**: All checks pass
- **NO-GO**: Any required check fails
- **WARN**: Non-required checks fail or warnings present

## Error Handling

When lint errors involve unhandled errors, follow this priority:

1. Panic - For programming errors that should never happen
2. Return Error - If function can return error
3. Modify Function - Add error return if possible
4. Log Error - Use slogutil if context available
5. Raise to Human - If no automated fix possible

## Output Format

Report results in the NASA-style Go/No-Go format:

` + "```" + `
┌──────────────────────────────────────────────────────────────────────────┐
│                        QA VALIDATION REPORT                              │
├──────────────────────────────────────────────────────────────────────────┤
│ Status: GO / NO-GO                                                       │
├──────────────────────────────────────────────────────────────────────────┤
│ Checks:                                                                  │
│   🟢 go-build       - Code compiles successfully                        │
│   🟢 go-mod-tidy    - Dependencies are clean                            │
│   🟢 go-fmt         - Code is formatted                                 │
│   🟢 go-vet         - No static analysis issues                         │
│   🟢 golangci-lint  - No lint errors                                    │
│   🟢 go-test        - All tests pass                                    │
│   🟢 go-test-race   - No data races detected                            │
└──────────────────────────────────────────────────────────────────────────┘
` + "```" + `
`
}
