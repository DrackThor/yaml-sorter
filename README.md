# yaml-sort

A command-line tool written in Go that sorts YAML files alphabetically by their keys while preserving the structure and hierarchy.

## Features

- Sort YAML files **recursively** by keys (every mapping level, including inside lists)
- Optional **Kubernetes manifest** mode (`-k`): root keys in fixed order (`apiVersion`, `kind`, `metadata`, `spec`, …), rest alphabetical
- **Config file** (`-c`): sort lists of objects by a specific key (e.g. `spec.egress` by `name`) for stable, deterministic order
- In-place sorting option (`-i`)
- Output to a new file option (`-o`)
- Standard output support
- Comprehensive test coverage

## Installation

### From Source

```bash
git clone https://github.com/drackthor/yaml-sort.git
cd yaml-sort
go build -o yaml-sort .
```

### Using Go Install

```bash
go install github.com/drackthor/yaml-sort@latest
```

### Releases

Pushing to `main` runs [semantic-release](https://github.com/go-semantic-release/semantic-release): the next version is derived from **conventional commit messages** since the last tag (`feat:` → minor, `fix:` → patch, `BREAKING CHANGE` → major). A GitHub release and tag are created automatically; the Build workflow then attaches binaries for Linux, macOS, and Windows. Download them from the [Releases](https://github.com/drackthor/yaml-sort/releases) page.

## Usage

### Basic Usage

Sort a YAML file and output to stdout:

```bash
yaml-sort file.yaml
```

### In-place Sorting

Sort a file in-place, replacing the original file:

```bash
yaml-sort -i file.yaml
# or
yaml-sort --inplace file.yaml
```

### Output to File

Sort a file and write the result to a new file:

```bash
yaml-sort -o sorted.yaml file.yaml
# or
yaml-sort --output sorted.yaml file.yaml
```

### Kubernetes manifests (-k)

For Kubernetes-style YAML (e.g. `kind`, `apiVersion`, `metadata`, `spec`), use `-k` so the **root** keys are output in a fixed order instead of A–Z:

- `apiVersion`, `kind`, `metadata`, `spec`, `data`, `status`, then any other root keys alphabetically.

Everything under those keys (e.g. under `metadata` or `spec`) is still sorted **recursively** and alphabetically.

```bash
yaml-sort -k deployment.yaml
yaml-sort -k -o sorted.yaml manifest.yaml
```

### Sort lists of objects by key (config file, `-c`)

For YAML with **lists of objects** (e.g. `spec.egress`, `spec.ingress` in NeuVector CRDs), you can sort each list by a field (e.g. `name`) so the order is stable. Use a **config file** and pass it with `-c`.

**Config format** (YAML):

```yaml
listSortKeys:
  - path: spec.egress   # dot-separated path from document root to the list
    key: name          # field inside each list element to sort by
  - path: spec.ingress
    key: name
  - path: spec.process
    key: name
```

- **path**: Where the list lives (e.g. `spec.egress`, `metadata.labels`).
- **key**: For each item in that list (must be a mapping), sort by this key’s value; missing keys sort as empty string.

Example with NeuVector runtime group and K8s root order:

```bash
cp .yaml-sort.example.yaml .yaml-sort.yaml
yaml-sort -k -c .yaml-sort.yaml -o sorted.yaml neuvector-runtime-group.yaml
```

An example config is in the repo: [.yaml-sort.example.yaml](.yaml-sort.example.yaml). For before/after examples of list sorting, see [EXAMPLES.md](EXAMPLES.md).

### Help

Display help information:

```bash
yaml-sort -h
# or
yaml-sort --help
```

| Flag        | Short | Description                                                  |
|-------------|-------|--------------------------------------------------------------|
| `--inplace` | `-i`  | Write output back to the input file                          |
| `--output`  | `-o`  | Write output to a file                                       |
| `--k8s`     | `-k`  | Use K8s root key order (apiVersion, kind, metadata, spec, …) |
| `--config`  | `-c`  | Config file for list sort keys (path → key)                  |

## Examples

**[EXAMPLES.md](EXAMPLES.md)** has detailed before/after examples, including:

- Simple and nested mapping sort
- Kubernetes manifest root order (`-k`)
- **Sorting lists of objects by key** with a config file (`-c`): `spec.egress`, `spec.ingress` by `name`
- Multiple list paths and combined K8s + list sort (e.g. NeuVector-style manifests)

Quick illustration (default alphabetical sort):

```yaml
# Before
zebra:
  c: value3
  a: value1
  b: value2
apple: value
banana: value

# After (yaml-sort file.yaml)
apple: value
banana: value
zebra:
  a: value1
  b: value2
  c: value3
```

## Development

### Prerequisites

- Go 1.21 or later
- [golangci-lint](https://golangci-lint.run/) (for lint and pre-commit)
- Optional: [pre-commit](https://pre-commit.com/) (for commit-msg and pre-commit hooks)

### Go commands

Run these from the repository root.

| Command                            | What it does                                                                                                                                        |
|------------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------|
| `go mod tidy`                      | Updates `go.mod` and `go.sum`: adds any missing dependencies, removes unused ones, and pins versions. Run after cloning or when you change imports. |
| `go mod download`                  | Downloads all modules listed in `go.mod` into the module cache (optional; `go build` and `go test` do this automatically).                          |
| `go build -o yaml-sort .`          | Builds the current package (`.`) and writes the executable to `yaml-sort`.                                                                          |
| `go test ./...`                    | Runs all tests in the module (unit tests and `test-cases/` integration tests).                                                                      |
| `go test -short ./...`             | Same as above but skips long-running tests if the code uses `testing.Short()`.                                                                      |
| `go test -cover ./...`             | Runs tests and prints per-package coverage.                                                                                                         |
| `go test -race ./...`              | Runs tests with the race detector to find data races.                                                                                               |
| `go fmt ./...`                     | Formats all Go files in the module (standard style).                                                                                                |
| `go vet ./...`                     | Runs the Go vet tool for static checks (e.g. suspicious constructs, unreachable code).                                                              |
| `go run ./scripts/gen_expected.go` | Generates sorted YAML files from `test-cases/inputs/` into `test-cases/expected/`.                                                                  |
| `go run . file.yaml`               | Builds and runs the CLI in one step (e.g. `go run . -o out.yaml file.yaml`).                                                                        |

**First-time setup:** after cloning, run `go mod tidy` so `go.sum` is populated; then `go build -o yaml-sort .` and `go test ./...` should work.

### Running Tests

```bash
go test ./...
```

Integration tests use real-world YAML under `test-cases/inputs/` (e.g. SUSE NeuVector CRDs). See [test-cases/README.md](test-cases/README.md).

Run tests with coverage:

```bash
go test -cover ./...
```

Run tests with race detection:

```bash
go test -race ./...
```

### Linting

This project uses `golangci-lint`. To run linting:

```bash
golangci-lint run
```

### Pre-commit hooks

Pre-commit enforces **conventional commit messages** and **Go best practices** (format, vet, tests, lint).

**Requirements:** Python 3 with `pre-commit`, Go 1.21+, and `golangci-lint` on your PATH.

Install hooks (run from repo root):

```bash
pre-commit install --install-hook-types commit-msg pre-commit
```

After this, each commit will:

- Reject messages that don’t follow [Conventional Commits](https://www.conventionalcommits.org/) (e.g. `feat: add X`, `fix: correct Y`).
- Run `go fmt ./...`, `go vet ./...`, `go test -short ./...`, and `golangci-lint run ./...`.

Run manually on all files:

```bash
pre-commit run --all-files
```

### Building

```bash
go build -o yaml-sort .
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Contribution Guidelines

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for new functionality (including real-world cases in `test-cases/inputs/` if relevant)
5. Ensure all tests pass (`go test ./...`)
6. Ensure linting passes (`golangci-lint run`)
7. Use [Conventional Commits](https://www.conventionalcommits.org/) for commit messages (e.g. `feat: add X`, `fix: correct Y`). Pre-commit will enforce this if hooks are installed.
8. Push to the branch (`git push origin feature/amazing-feature`)
9. Open a Pull Request

### Code Style

- Follow Go best practices and conventions
- Use `gofmt` for code formatting
- Write tests for all new features
- Keep functions small and focused
- Add comments for exported functions and types

### Testing

- Write tests using the standard `testing` package
- Aim for high test coverage
- Use table-driven tests where appropriate
- Test both success and error cases

## YAML terminology used in this project

The code and comments refer to the following YAML and tree concepts (as in [gopkg.in/yaml.v3](https://github.com/go-yaml/yaml)):

| Term               | Meaning                                                                                                                                                        |
|--------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **node**           | A single element in the YAML tree: either a **scalar** (string/number/bool/null), a **mapping**, or a **sequence**. Represented as `*yaml.Node`.               |
| **DocumentNode**   | The root of a YAML document. Its `Content` is a single child (usually a mapping or sequence).                                                                  |
| **MappingNode**    | A YAML **mapping** (key-value structure). In the tree, `Content` is a flat list of alternating key nodes and value nodes: `[key1, value1, key2, value2, ...]`. |
| **SequenceNode**   | A YAML **sequence** (list/array). `Content` is a list of child nodes, one per element.                                                                         |
| **key**            | In a mapping, the first of each pair of nodes (the key name, usually a scalar).                                                                                |
| **value**          | In a mapping, the second of each pair; can be a scalar, another mapping, or a sequence.                                                                        |
| **Content**        | The `Content` field of a `yaml.Node`: the slice of child nodes. For a mapping it’s key/value pairs; for a sequence it’s the list items.                        |
| **recursive sort** | We sort each mapping’s keys and then **recurse** into each value (and into each sequence element), so every nested map and list is sorted too.                 |

So when we say “sort a mapping node”, we mean: take its key-value pairs, sort them by key (alphabetically or, at the root with `-k`, by a fixed K8s order), then recurse into each value.

## Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) for CLI functionality
- Uses [gopkg.in/yaml.v3](https://github.com/go-yaml/yaml) for YAML parsing

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.