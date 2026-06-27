package gocaps

import (
	"github.com/plexusone/assistantkit/capabilities/core"
)

// GoBuild returns the go build capability.
// This capability compiles Go code and verifies it builds successfully.
func GoBuild() *core.Capability {
	c := core.NewCapability("go-build", "Build Go code and verify compilation")

	c.WithCommand("go", "build", "./...")
	c.WithTriggers("build", "go build", "compile")
	c.WithDependencies("go")
	c.WithLanguage("go")
	c.AsRequired()

	c.WithInstructions(`# Go Build Skill

Build Go code to verify compilation and catch type errors.

## Running Build

` + "```bash" + `
# Build all packages
go build ./...

# Build specific package
go build ./pkg/...

# Build with output binary
go build -o myapp ./cmd/myapp

# Build for different OS/arch
GOOS=linux GOARCH=amd64 go build -o myapp-linux ./cmd/myapp

# Build with race detection
go build -race ./...

# Build with debug info stripped (smaller binary)
go build -ldflags="-s -w" -o myapp ./cmd/myapp
` + "```" + `

## Common Build Errors

### Import Errors
- Missing dependency: Run ` + "`go mod tidy`" + `
- Cyclic import: Restructure packages
- Wrong import path: Check module name in go.mod

### Type Errors
- Type mismatch: Check function signatures
- Undefined: Check spelling and imports
- Cannot convert: Add explicit conversion

### Syntax Errors
- Missing semicolon: Usually missing brace or paren
- Unexpected token: Check for typos

## Build Flags

| Flag | Description |
|------|-------------|
| ` + "`-v`" + ` | Print package names as compiled |
| ` + "`-race`" + ` | Enable race detection |
| ` + "`-o`" + ` | Output binary name |
| ` + "`-ldflags`" + ` | Linker flags (version, strip) |
| ` + "`-tags`" + ` | Build tags to include |

## Workflow

1. **Build**: ` + "`go build ./...`" + `
2. **Fix errors**: Address compilation errors
3. **Verify**: Re-run build
4. **Test**: Run tests to verify behavior
`)

	return c
}

// GoModTidy returns the go mod tidy capability.
// This capability ensures go.mod and go.sum are in sync.
func GoModTidy() *core.Capability {
	c := core.NewCapability("go-mod-tidy", "Tidy Go module dependencies")

	c.WithCommand("go", "mod", "tidy")
	c.WithTriggers("tidy", "go mod tidy", "dependencies")
	c.WithDependencies("go")
	c.WithLanguage("go")
	c.AsRequired()

	c.WithInstructions(`# Go Mod Tidy Skill

Ensure Go module dependencies are correct and go.mod/go.sum are in sync.

## Running Mod Tidy

` + "```bash" + `
# Tidy dependencies
go mod tidy

# Verify dependencies without modifying files
go mod tidy -diff

# Verbose output
go mod tidy -v
` + "```" + `

## What Mod Tidy Does

1. **Adds** missing module requirements
2. **Removes** unused dependencies
3. **Updates** go.sum with required hashes
4. **Validates** go.mod and go.sum consistency

## Common Issues

### Missing go.sum Entry
Run ` + "`go mod tidy`" + ` to regenerate go.sum.

### Module Not Found
Check import path matches module name in remote repo's go.mod.

### Version Conflict
Use ` + "`go mod graph`" + ` to identify conflicting requirements.

### Local Replace Directive
Remove ` + "`replace`" + ` directives pointing to local paths before pushing:

` + "```go" + `
// BAD (for push)
replace github.com/example/pkg => ../pkg

// GOOD (for push)
require github.com/example/pkg v1.2.3
` + "```" + `

## Verification

After tidy, verify no uncommitted changes:

` + "```bash" + `
git diff go.mod go.sum
` + "```" + `

If there are changes, commit them with the related code changes.
`)

	return c
}

// GoVet returns the go vet capability.
// This capability runs go vet for static analysis.
func GoVet() *core.Capability {
	c := core.NewCapability("go-vet", "Run go vet for static analysis")

	c.WithCommand("go", "vet", "./...")
	c.WithTriggers("vet", "go vet", "static analysis")
	c.WithDependencies("go")
	c.WithLanguage("go")
	c.AsRequired()

	c.WithInstructions(`# Go Vet Skill

Run go vet to find suspicious code that may indicate bugs.

## Running Go Vet

` + "```bash" + `
# Vet all packages
go vet ./...

# Vet specific package
go vet ./pkg/...

# Enable all checks
go vet -all ./...
` + "```" + `

## Common Vet Findings

### Printf Format Errors
Wrong format specifier for argument type.

### Unreachable Code
Code after return/panic that will never execute.

### Shadowed Variables
Variable declaration shadows outer scope variable.

### Unused Results
Result of a function call is ignored.

### Struct Tag Errors
Invalid or malformed struct tags.

## Workflow

1. **Run vet**: ` + "`go vet ./...`" + `
2. **Review findings**: Each finding indicates potential bug
3. **Fix issues**: Address the root cause
4. **Verify**: Re-run vet to confirm fixes
`)

	return c
}

// GoFmt returns the go fmt capability.
// This capability ensures code is properly formatted.
func GoFmt() *core.Capability {
	c := core.NewCapability("go-fmt", "Format Go code with gofmt")

	c.WithCommand("gofmt", "-l", "-d", ".")
	c.WithTriggers("fmt", "gofmt", "format")
	c.WithDependencies("gofmt")
	c.WithLanguage("go")
	c.AsRequired()

	c.WithInstructions(`# Go Fmt Skill

Format Go code to follow standard style.

## Running Go Fmt

` + "```bash" + `
# Check formatting (show diff)
gofmt -l -d .

# Format files in place
gofmt -w .

# Using go fmt (formats all packages)
go fmt ./...

# Check if files need formatting (for CI)
test -z "$(gofmt -l .)"
` + "```" + `

## What Go Fmt Does

- Standardizes indentation (tabs)
- Aligns struct fields
- Formats import blocks
- Normalizes whitespace

## Best Practices

1. **Run on save**: Configure editor to format on save
2. **Pre-commit hook**: Format before committing
3. **CI check**: Fail CI if unformatted code is pushed

## Integration

Most editors and IDEs auto-format Go code:
- VS Code: gopls extension
- GoLand: Built-in formatter
- vim-go: ` + "`:GoFmt`" + ` command
`)

	return c
}
