# yaml-sort

A command-line tool written in Go that sorts YAML files alphabetically by their keys while preserving the structure and hierarchy.

## Features

- Sort YAML files alphabetically by keys
- Preserve nested structure and hierarchy
- In-place sorting option
- Output to a new file option
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

### Help

Display help information:

```bash
yaml-sort -h
# or
yaml-sort --help
```

## Examples

### Example Input

```yaml
zebra:
  c: value3
  a: value1
  b: value2
apple: value
banana: value
```

### Example Output

```yaml
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
- Make (optional, for convenience commands)

### Running Tests

```bash
go test ./...
```

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
4. Add tests for new functionality
5. Ensure all tests pass (`go test ./...`)
6. Ensure linting passes (`golangci-lint run`)
7. Commit your changes (`git commit -m 'Add some amazing feature'`)
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

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) for CLI functionality
- Uses [gopkg.in/yaml.v3](https://github.com/go-yaml/yaml) for YAML parsing
