# Test cases

Real-world YAML samples used for integration tests.

## Layout

- **`inputs/`** – Sample YAML files (e.g. SUSE NeuVector Kubernetes CRDs). These are sorted by the test to verify the sorter handles real-world structure.
- **`expected/`** – (Optional) Canonical sorted output for each input. Generated with `go run ./scripts/gen_expected.go` or by running the CLI. Tests can optionally diff against these.

## Inputs

| File | Description |
|------|-------------|
| `grp-ui-old.yaml` | NvSecurityRule (NeuVector) – older variant, Monitor mode |
| `grp-ui-new.yaml` | NvSecurityRule – newer variant, Protect mode, more rules |
| `cfgGroupsExport.yaml` | NvSecurityRule export snippet with comments |
| `cfgAdmissionRules.yaml` | NvAdmissionControlSecurityRule with deny rules |

## Running tests

From repository root:

```bash
go test ./test-cases/...
```

To run only integration tests:

```bash
go test -run RealWorld ./test-cases/...
```

## Adding new cases

1. Add `inputs/<name>.yaml`.
2. Integration tests will pick it up automatically.
3. Optionally add `expected/<name>.yaml` for strict output comparison.
