package gocaps

import (
	"github.com/plexusone/assistantkit/capabilities/core"
)

// GoTest returns the go test capability.
// This capability runs go test and provides guidance for analyzing
// and fixing test failures.
func GoTest() *core.Capability {
	c := core.NewCapability("go-test", "Run Go tests and analyze failures")

	c.WithCommand("go", "test", "-v", "./...")
	c.WithTriggers("test", "go test", "tests", "testing")
	c.WithDependencies("go")
	c.WithLanguage("go")
	c.AsRequired()

	c.FailurePattern = `--- FAIL:`
	c.SuccessPattern = `PASS`

	c.WithInstructions(`# Go Test Skill

Run Go tests and analyze failures to provide actionable fixes.

## Running Tests

` + "```bash" + `
# Run all tests
go test -v ./...

# Run tests in specific package
go test -v ./pkg/...

# Run specific test by name
go test -v -run TestFunctionName ./...

# Run tests with coverage
go test -v -cover ./...

# Run tests with race detection
go test -v -race ./...

# Run short tests only
go test -v -short ./...

# Run tests with timeout
go test -v -timeout 30s ./...
` + "```" + `

## Analyzing Failures

When tests fail, analyze:

1. **Test Name**: Which test(s) failed
2. **Error Message**: What was expected vs actual
3. **Stack Trace**: Where the failure occurred
4. **Test Setup**: Are fixtures/mocks correct

## Common Failure Patterns

### Assertion Failures
- Check expected vs actual values
- Verify test data is correct
- Ensure mocks return expected values

### Timeout Failures
- Check for infinite loops
- Verify external dependencies are mocked
- Consider increasing timeout for slow tests

### Race Conditions
- Use ` + "`-race`" + ` flag to detect races
- Add proper synchronization (mutex, channels)
- Review concurrent access patterns

### Setup/Teardown Issues
- Ensure proper cleanup in ` + "`defer`" + ` or ` + "`t.Cleanup()`" + `
- Check test isolation (no shared state)
- Verify test fixtures exist

## Best Practices

1. **Table-driven tests**: Use subtests for multiple cases
2. **Parallel tests**: Use ` + "`t.Parallel()`" + ` when safe
3. **Test helpers**: Use ` + "`t.Helper()`" + ` for better error reporting
4. **Cleanup**: Always clean up resources with ` + "`t.Cleanup()`" + `

## Workflow

1. **Run tests**: ` + "`go test -v ./...`" + `
2. **Identify failures**: Note failed test names and errors
3. **Analyze root cause**: Read test code and failure message
4. **Fix implementation**: Update code to pass test
5. **Verify fix**: Re-run failed tests
6. **Run full suite**: Ensure no regressions
`)

	return c
}

// GoTestCoverage returns the go test coverage capability.
// This capability runs go test with coverage reporting.
func GoTestCoverage() *core.Capability {
	c := core.NewCapability("go-test-coverage", "Run Go tests with coverage analysis")

	c.WithCommand("go", "test", "-v", "-cover", "-coverprofile=coverage.out", "./...")
	c.WithTriggers("coverage", "go coverage", "test coverage")
	c.WithDependencies("go")
	c.WithLanguage("go")

	c.WithInstructions(`# Go Test Coverage Skill

Run Go tests with coverage analysis to identify untested code.

## Running Coverage

` + "```bash" + `
# Generate coverage profile
go test -v -cover -coverprofile=coverage.out ./...

# View coverage in terminal
go tool cover -func=coverage.out

# Generate HTML report
go tool cover -html=coverage.out -o coverage.html

# View coverage by function
go tool cover -func=coverage.out | grep -v "100.0%"
` + "```" + `

## Coverage Targets

| Level | Target | Description |
|-------|--------|-------------|
| Good | 80%+ | Comprehensive test coverage |
| Acceptable | 60-80% | Most code paths tested |
| Needs Work | < 60% | Significant gaps in testing |

## Improving Coverage

1. **Identify gaps**: Use HTML report to find uncovered lines
2. **Prioritize**: Focus on critical paths and error handling
3. **Add tests**: Write tests for uncovered code
4. **Verify**: Re-run coverage to confirm improvement

## Excluding Files

Some files may be excluded from coverage:

- Generated code (` + "`*_gen.go`" + `)
- Test utilities (` + "`*_test.go`" + `)
- Main packages (` + "`cmd/`" + `)
`)

	return c
}

// GoTestRace returns the go test race detection capability.
func GoTestRace() *core.Capability {
	c := core.NewCapability("go-test-race", "Run Go tests with race detection")

	c.WithCommand("go", "test", "-v", "-race", "./...")
	c.WithTriggers("race", "race detection", "data race")
	c.WithDependencies("go")
	c.WithLanguage("go")
	c.AsRequired()

	c.FailurePattern = `DATA RACE`

	c.WithInstructions(`# Go Race Detection Skill

Run Go tests with the race detector to find data races.

## Running Race Detection

` + "```bash" + `
# Run all tests with race detection
go test -v -race ./...

# Run specific test with race detection
go test -v -race -run TestConcurrent ./...
` + "```" + `

## Understanding Data Races

A data race occurs when:
- Two goroutines access the same memory location
- At least one access is a write
- There is no synchronization between accesses

## Common Race Patterns

### Shared State Without Lock
` + "```go" + `
// BAD
var counter int
go func() { counter++ }()
go func() { counter++ }()

// GOOD
var mu sync.Mutex
var counter int
go func() { mu.Lock(); counter++; mu.Unlock() }()
go func() { mu.Lock(); counter++; mu.Unlock() }()
` + "```" + `

### Loop Variable Capture
` + "```go" + `
// BAD
for _, v := range values {
    go func() { process(v) }()  // v captured by reference
}

// GOOD
for _, v := range values {
    go func(v Value) { process(v) }(v)  // v passed by value
}
` + "```" + `

### Channel vs Mutex
- Use channels for communication between goroutines
- Use mutexes for protecting shared state

## Fixing Races

1. **Add synchronization**: Use ` + "`sync.Mutex`" + ` or ` + "`sync.RWMutex`" + `
2. **Use channels**: Send data between goroutines
3. **Use atomic operations**: For simple counters
4. **Redesign**: Eliminate shared state if possible
`)

	return c
}
