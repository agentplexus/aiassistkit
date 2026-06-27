// Package gocaps provides Go-specific capabilities for validation and QA.
package gocaps

import (
	"github.com/plexusone/assistantkit/capabilities/core"
)

// GolangciLint returns the golangci-lint capability.
// This capability runs golangci-lint and provides guidance for fixing
// lint errors following established patterns.
func GolangciLint() *core.Capability {
	c := core.NewCapability("golangci-lint", "Run golangci-lint and fix Go lint errors with proper error handling patterns")

	c.WithCommand("golangci-lint", "run", "./...")
	c.WithTriggers("lint", "golangci-lint", "linter", "go lint")
	c.WithDependencies("golangci-lint")
	c.WithLanguage("go")
	c.AsRequired()

	c.ErrorHandling = &core.ErrorHandlingGuide{
		Priorities: []core.ErrorPriority{
			{
				Level:       1,
				Name:        "Panic",
				Description: "Use panic() when an error should never happen (invariant violation, programming error)",
				Example: `// json.Marshal of simple struct should never fail
data, err := json.Marshal(simpleStruct)
if err != nil {
    panic(fmt.Sprintf("json.Marshal failed: %v", err))
}`,
			},
			{
				Level:       2,
				Name:        "Return Error",
				Description: "If the function signature can return an error, return it",
				Example: `func processData(data []byte) error {
    result, err := parseData(data)
    if err != nil {
        return fmt.Errorf("parse failed: %w", err)
    }
    return nil
}`,
			},
			{
				Level:       3,
				Name:        "Modify Function",
				Description: "If the function doesn't return error but can be modified, add error return",
				Example: `// Before
func loadConfig() Config { ... }

// After
func loadConfig() (Config, error) {
    data, err := os.ReadFile("config.json")
    if err != nil {
        return Config{}, fmt.Errorf("read config: %w", err)
    }
    return cfg, nil
}`,
			},
			{
				Level:       4,
				Name:        "Log Error",
				Description: "When function must fulfill an interface without error return, use logging",
				Example: `import "github.com/grokify/mogo/log/slogutil"

logger := slogutil.LoggerFromContext(ctx, nil)
if err != nil {
    slogutil.LogOrNot(ctx, logger, slog.LevelError,
        "operation failed",
        slog.String("error", err.Error()),
        slog.String("operation", "write"))
}`,
			},
			{
				Level:       5,
				Name:        "Raise to Human",
				Description: "If none of the above approaches work, raise to the human with explanation",
				Example: `// Cannot automatically fix because:
// - Function implements interface X which has no error return
// - There is no context available for logging
// - Panic is not appropriate for this recoverable error`,
			},
		},
		References: []core.PackageReference{
			{
				Package:     "github.com/grokify/mogo/log/slogutil",
				Description: "Logger utilities for context-based logging",
				Usage:       "slogutil.LoggerFromContext(ctx, nil)",
			},
			{
				Package:     "github.com/grokify/mogo/lintfix",
				Description: "Remediation database for common lint errors",
				Usage:       "db := lintfix.MustLoadRemediations()",
			},
			{
				Package:     "github.com/grokify/mogo/lintfix/gosec",
				Description: "Gosec nolint helpers",
				Usage:       "gosec.NolintG117(gosec.CommonReasons.OAuthTokenResponse)",
			},
		},
	}

	// Add common fixes
	c.AddCommonFix(core.CommonFix{
		Issue:       "errcheck",
		Description: "Unhandled error return value",
		Before: `file.Close()
writer.Flush()`,
		After: `defer func() {
    if cerr := file.Close(); cerr != nil && err == nil {
        err = cerr
    }
}()

if err := writer.Flush(); err != nil {
    return fmt.Errorf("flush failed: %w", err)
}`,
	})

	c.AddCommonFix(core.CommonFix{
		Issue:       "G306",
		Description: "File permissions too permissive",
		Before:      `os.WriteFile("file.txt", data, 0644)`,
		After:       `os.WriteFile("file.txt", data, 0o600)`,
	})

	c.AddCommonFix(core.CommonFix{
		Issue:       "G115",
		Description: "Integer overflow on conversion",
		Before:      `count := int(uint64Value)`,
		After: `if uint64Value > math.MaxInt {
    return fmt.Errorf("value %d exceeds max int", uint64Value)
}
count := int(uint64Value)`,
	})

	c.AddCommonFix(core.CommonFix{
		Issue:       "G118",
		Description: "Context in goroutines",
		Before: `go func() {
    doWork(context.Background())
}()`,
		After: `go func(ctx context.Context) {
    doWork(ctx)
}(ctx)`,
	})

	c.AddCommonFix(core.CommonFix{
		Issue:       "G120",
		Description: "HTTP form parsing without body limit",
		Before:      `r.ParseForm()`,
		After: `import "github.com/grokify/mogo/lintfix/gosec"

if err := gosec.LimitAndParseForm(w, r, gosec.G120MaxBytes.Webhook); err != nil {
    http.Error(w, "Bad Request", http.StatusBadRequest)
    return
}
value := r.Form.Get("key")  // NOT r.FormValue("key")`,
	})

	c.WithInstructions(generateGolangciLintInstructions())

	return c
}

func generateGolangciLintInstructions() string {
	return `# Go Lint Skill

Run ` + "`golangci-lint run`" + ` on Go projects and fix lint errors following established patterns.

## Running the Linter

` + "```bash" + `
# Run in current directory
golangci-lint run

# Run in specific directory
golangci-lint run ./path/to/dir/...

# Run with verbose output
golangci-lint run -v

# Run specific linters
golangci-lint run --enable=errcheck,gosec
` + "```" + `

## Error Handling Priority Order

**CRITICAL: All errors must be handled. Never assign errors to ` + "`_`" + `.**

When fixing unhandled error lint violations, follow this priority order:

### Priority 1: Panic (Programming Errors)

Use ` + "`panic()`" + ` when an error should never happen (invariant violation, programming error):

` + "```go" + `
// json.Marshal of simple struct should never fail
data, err := json.Marshal(simpleStruct)
if err != nil {
    panic(fmt.Sprintf("json.Marshal failed: %v", err))
}
` + "```" + `

### Priority 2: Return Error

If the function signature can return an error, return it:

` + "```go" + `
func processData(data []byte) error {
    result, err := parseData(data)
    if err != nil {
        return fmt.Errorf("parse failed: %w", err)
    }
    return nil
}
` + "```" + `

### Priority 3: Modify Function to Return Error

If the function doesn't return error but can be modified:

` + "```go" + `
// Before
func loadConfig() Config { ... }

// After
func loadConfig() (Config, error) {
    data, err := os.ReadFile("config.json")
    if err != nil {
        return Config{}, fmt.Errorf("read config: %w", err)
    }
    return cfg, nil
}
` + "```" + `

### Priority 4: Log Error (Interface Compliance)

When function must fulfill an interface without error return, use logging:

` + "```go" + `
import "github.com/grokify/mogo/log/slogutil"

logger := slogutil.LoggerFromContext(ctx, nil)
if err != nil {
    slogutil.LogOrNot(ctx, logger, slog.LevelError,
        "operation failed",
        slog.String("error", err.Error()),
        slog.String("operation", "write"))
}
` + "```" + `

### Priority 5: Raise to Human

If none of the above approaches work, raise to the human with explanation.

## Common Lint Fixes

### File Permissions (G306)
Always use ` + "`0o600`" + ` for file permissions, not ` + "`0o644`" + ` or ` + "`0644`" + `.

### Unhandled Errors (errcheck)
Handle all errors - use defer with named return or inline check.

### HTTP Form Parsing (G120)
Use body limit before parsing forms with ` + "`gosec.LimitAndParseForm`" + `.

### Integer Overflow (G115)
Add bounds check before integer conversion.

### Context in Goroutines (G118)
Pass context explicitly to goroutines.

## Workflow

1. **Run linter**: ` + "`golangci-lint run ./...`" + `
2. **Review errors**: Group by type (errcheck, gosec, etc.)
3. **Apply fixes**: Follow priority order above
4. **Verify**: Re-run linter to confirm fixes
5. **Commit**: Use conventional commit: ` + "`fix: resolve lint errors`" + `

## Resources

- Error logging: ` + "`github.com/grokify/mogo/log/slogutil`" + `
- Lint fixes: ` + "`github.com/grokify/mogo/lintfix`" + `
- Gosec helpers: ` + "`github.com/grokify/mogo/lintfix/gosec`" + `
`
}
