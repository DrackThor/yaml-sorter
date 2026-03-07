# Agent Guidelines for ysort

This document provides guidelines for AI agents working on this Go project.

## Project Snapshot

- CLI binary name: `ysort`.
- Go module path: `github.com/drackthor/ysort`.
- Primary behavior: recursively sort YAML mappings, with optional Kubernetes root-key order and configurable list sorting.

## Project Structure

```
ysort/
├── cmd/                         # CLI command implementation (Cobra)
│   ├── command_name.go          # Resolve displayed command name from argv0
│   ├── root.go                  # Root command and flags (-i, -o, -k, -c)
│   └── root_test.go             # Unit tests for command name logic
├── internal/                    # Internal packages (not exported)
│   ├── config/
│   │   └── config.go            # Load list-sort config from YAML file
│   └── sorter/
│       ├── sorter.go            # YAML sorting logic and options
│       └── sorter_test.go       # Unit tests for sorter behavior
├── test-cases/                  # Integration tests with real-world YAML
│   ├── inputs/                  # Sample YAML manifests/CRDs
│   ├── expected/                # Canonical sorted outputs
│   ├── integration_test.go      # Integration coverage for test-cases
│   └── README.md
├── scripts/
│   └── gen_expected.go          # Regenerate expected sorted fixtures
├── .ysort.example.yaml          # Example config for list sort rules
├── .golangci.yml                # Lint configuration
├── .goreleaser.yaml             # Release build configuration
├── EXAMPLES.md                  # Before/after examples
├── README.md
├── main.go                      # Application entry point
├── go.mod                       # Module definition
├── .pre-commit-config.yaml      # Conventional commits + Go hooks
└── .github/workflows/           # CI workflows (PR + release)
```

## Go Best Practices

### Code Organization

1. **Package Structure**: Keep command wiring in `cmd/` and implementation logic in `internal/`.
2. **Single Responsibility**: Keep functions focused and composable.
3. **Error Context**: Wrap errors with `%w` and contextual prefixes.
4. **Flag-to-Option Mapping**: Keep CLI parsing separate from sorting internals.

### Naming Conventions

- Use camelCase for unexported names.
- Use PascalCase for exported names.
- Use short and explicit names.
- Avoid unclear abbreviations.
- Use semantic line breaks (one sentence per line) in documentation when practical.

### Error Handling

- Check and return errors explicitly.
- Wrap errors using `fmt.Errorf("...: %w", err)`.
- Do not swallow errors.
- Write actionable user-facing error messages for CLI failures.

### Function Design

- Keep functions small and testable.
- Prefer pure helper functions where possible.
- Keep parameter counts low.
- Return `error` as the last return value.

### Testing

1. Place tests beside source files with `_test.go` suffix.
2. Use descriptive test names such as `TestSortYAMLK8s_RootOrder`.
3. Prefer table-driven tests for behavior matrices.
4. Cover both recursive map sorting and list-sort rule behavior.
5. Keep integration fixtures in `test-cases/` representative of real manifests.

Example:

```go
func TestSortYAML_SimpleMapping(t *testing.T) {
    // test implementation
}
```

### Code Style

1. **Formatting**: Always run `gofmt`.
2. **Comments**: Follow the `### Comments` section below.
3. **Imports**: Group imports by stdlib, third-party, and local packages.
4. **Line Length**: Keep lines under 100 chars when reasonable.

### Dependencies

1. Keep dependencies minimal and justified.
2. Keep `go.mod`/`go.sum` tidy and reproducible.
3. Avoid adding new dependencies for small tasks that stdlib can handle.

### CLI Design

1. Use Cobra conventions for help and flag descriptions.
2. Keep short and long flags aligned (`-i/--inplace`, `-o/--output`, `-k/--k8s`, `-c/--config`).
3. Reject invalid flag combinations with clear errors.
4. Ensure stdout behavior remains useful for shell pipelines.

### Performance

1. Minimize allocations in recursive sorting paths.
2. Pre-allocate slices/maps when capacity is known.
3. Keep comparisons deterministic for stable output.
4. Benchmark before introducing complexity.

### Security

1. Validate file paths and file IO failures.
2. Treat config and YAML input as untrusted.
3. Avoid leaking sensitive file system details in errors when not needed.

### Comments

1. **Distinguish purpose**:
   - Ordinary comments explain non-obvious intent and constraints, especially **why** code exists.
   - Doc comments describe the API contract for users of a package/symbol.
2. **Doc comment placement**:
   - Put doc comments immediately above top-level `package`, `const`, `var`, `type`, and `func` declarations.
   - Do not insert a blank line between the comment and declaration.
3. **Exported API coverage**:
   - Every exported name should have a doc comment.
   - Add doc comments for unexported names only when they clarify behavior that is otherwise hard to infer.
4. **Opening sentence format**:
   - Use complete sentences.
   - Start with the declared name (`Package sorter ...`, `SortYAML ...`, `Client ...`).
   - For command package comments, start with the program name (`Ysort ...`).
5. **Package comments**:
   - Keep exactly one package comment per package.
   - For long package overviews, prefer a dedicated `doc.go` file.
   - Ensure the first sentence stands alone as a useful summary.
6. **Doc syntax and tooling**:
   - Write doc comments using Go doc syntax (a simplified Markdown subset): paragraphs, headings, lists, links, and indented code/preformatted blocks.
   - Prefer doc links like `[sorter.Options]` or `[path/filepath]` where useful.
   - Run `gofmt`; it canonicalizes doc comment formatting.
7. **Content boundaries**:
   - Keep implementation details (algorithm internals, temporary workarounds) in ordinary inline comments near code.
   - Keep doc comments focused on externally observable behavior and contracts.
8. **Deprecations and directives**:
   - Start deprecation notices with `Deprecated:` in their own paragraph.
   - Keep directives like `//go:generate` separate; they are not doc prose.

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

1. Run formatting: `go fmt ./...`.
2. Run vet: `go vet ./...`.
3. Run tests: `go test ./...`.
4. Run lint: `golangci-lint run ./...`.
5. If fixtures changed, regenerate expected outputs: `go run ./scripts/gen_expected.go`.

### CI/CD

- PR workflow (`.github/workflows/pr.yml`) runs lint, tests, and build.
- Release workflow (`.github/workflows/release.yml`) runs lint + tests, then semantic release with GoReleaser on pushes to `main`.
- Pre-commit hooks enforce conventional commit format and key Go checks when installed.

## When Adding Features

1. Start with a failing test.
2. Implement the minimal fix/feature.
3. Refactor while keeping tests green.
4. Update docs (`README.md`, `EXAMPLES.md`, config examples) when behavior changes.
5. Run format, vet, tests, and lint before finalizing.

## When Fixing Bugs

1. Reproduce with a targeted test.
2. Implement the fix.
3. Verify the new test and full suite pass.
4. Check for regressions in `test-cases/` fixtures.

## Code Review Checklist

- [ ] Code matches current Go conventions and project patterns.
- [ ] Comments/doc comments follow the `### Comments` guidance.
- [ ] Tests cover new behavior and pass.
- [ ] Error handling is explicit and contextual.
- [ ] User-facing documentation is updated where needed.
- [ ] Lint and formatting checks pass.
- [ ] No unnecessary dependencies were introduced.
- [ ] Performance and determinism implications were considered.

## Resources

- [How To Write Comments in Go (DigitalOcean)](https://www.digitalocean.com/community/tutorials/how-to-write-comments-in-go)
- [Godoc: documenting Go code](https://go.dev/blog/godoc)
- [Go Doc Comments](https://go.dev/doc/comment)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Go Project Layout Discussion](https://github.com/golang-standards/project-layout)
