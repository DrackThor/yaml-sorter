# Agent Guidelines for yaml-sort

This document provides guidelines for AI agents working on this Go project.

## Project Structure

```
yaml-sort/
├── cmd/              # CLI command implementations
│   └── root.go      # Root command with sort functionality
├── internal/         # Internal packages (not exported)
│   └── sorter/      # YAML sorting logic
│       ├── sorter.go
│       └── sorter_test.go
├── test-cases/       # Real-world YAML test data and integration tests
│   ├── inputs/      # Sample inputs (e.g. SUSE NeuVector CRDs)
│   ├── expected/    # Optional: canonical sorted output (generate with scripts/gen_expected.go)
│   ├── integration_test.go
│   └── README.md
├── scripts/         # Helper scripts (e.g. gen_expected.go)
├── main.go          # Application entry point
├── go.mod           # Go module definition
├── .pre-commit-config.yaml  # Conventional commits + Go hooks
└── .github/         # GitHub Actions workflows
```

## Go Best Practices

### Code Organization

1. **Package Structure**: Follow standard Go project layout
    - `cmd/` for application entry points
    - `internal/` for private application code
    - `pkg/` for library code (if needed)

2. **Naming Conventions**:
    - Use camelCase for unexported functions and variables
    - Use PascalCase for exported functions, types, and constants
    - Use short, descriptive names
    - Avoid abbreviations unless widely understood
    - follow semantic linebreaks, aka only one sentence per line

3. **Error Handling**:
    - Always check and handle errors explicitly
    - Use `fmt.Errorf` with `%w` verb for error wrapping
    - Return errors, don't ignore them
    - Provide context in error messages

4. **Function Design**:
    - Keep functions small and focused (single responsibility)
    - Prefer pure functions when possible
    - Limit function parameters (max 3-4)
    - Return errors as the last return value

### Testing

1. **Test Files**: Place test files next to source files with `_test.go` suffix
2. **Test Functions**: Use `Test` prefix for unit tests, `Benchmark` for benchmarks
3. **Table-Driven Tests**: Use table-driven tests for multiple test cases
4. **Test Coverage**: Aim for high test coverage, especially for core logic
5. **Test Naming**: Use descriptive test names: `TestFunctionName_Scenario`

Example:

```go
func TestSortYAML_SimpleMapping(t *testing.T) {
    // test implementation
}
```

### Code Style

1. **Formatting**: Always use `gofmt` or `goimports`
2. **Comments**:
    - Add comments for exported functions, types, and constants
    - Use complete sentences
    - Start with the name of the thing being described
3. **Imports**: Group imports: stdlib, third-party, local
4. **Line Length**: Keep lines under 100 characters when possible

### Dependencies

1. **Minimal Dependencies**: Keep dependencies minimal and well-justified
2. **Version Pinning**: Use specific versions in `go.mod`
3. **Vendor**: Consider vendoring dependencies for production builds

### CLI Design

1. **Cobra Commands**: Use Cobra for CLI structure
2. **Flag Naming**: Use short (`-i`) and long (`--inplace`) flag forms
3. **Error Messages**: Provide clear, actionable error messages
4. **Help Text**: Include comprehensive help text for commands

### Performance

1. **Allocations**: Minimize allocations in hot paths
2. **String Operations**: Use `strings.Builder` for string concatenation
3. **Slices**: Pre-allocate slices when size is known
4. **Profiling**: Use `go test -bench` and `pprof` for performance analysis

### Security

1. **File Operations**: Validate file paths and permissions
2. **Input Validation**: Validate all user inputs
3. **Error Messages**: Don't leak sensitive information in errors

## Common Patterns

### Error Wrapping

```go
if err != nil {
    return fmt.Errorf("failed to read file: %w", err)
}
```

### Table-Driven Tests

```go
tests := []struct {
    name     string
    input    string
    expected string
    wantErr  bool
}{
    {
        name:     "simple case",
        input:    "test",
        expected: "result",
        wantErr:  false,
    },
}
```

### Context Usage

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
```

## Workflow Integration

### Before Committing

1. Run tests: `go test ./...` (includes `test-cases/` integration tests with real-world YAML).
2. Run linter: `golangci-lint run`
3. Format code: `go fmt ./...`
4. Check for unused imports: `goimports -l .`
5. Use [Conventional Commits](https://www.conventionalcommits.org/) for commit messages. Pre-commit hooks (see `.pre-commit-config.yaml`) enforce this and the above steps when installed.

### CI/CD

- Tests run on every push and PR
- Linting runs on every push and PR
- Builds run on multiple platforms
- Coverage reports are generated

## When Adding Features

1. **Start with Tests**: Write tests first (TDD approach)
2. **Implement Feature**: Write minimal code to pass tests
3. **Refactor**: Improve code while keeping tests green
4. **Document**: Update README.md and add code comments
5. **Lint**: Ensure code passes linting

## When Fixing Bugs

1. **Reproduce**: Write a test that reproduces the bug
2. **Fix**: Implement the fix
3. **Verify**: Ensure the test passes and no regressions

## Code Review Checklist

- [ ] Code follows Go conventions
- [ ] Tests are included and passing
- [ ] Error handling is proper
- [ ] Documentation is updated
- [ ] Linting passes
- [ ] No unnecessary dependencies added
- [ ] Performance considerations addressed

## Resources

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Go Best Practices](https://github.com/golang-standards/project-layout)
