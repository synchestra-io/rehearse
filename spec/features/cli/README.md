# Feature: CLI

**Status:** In Progress

## Summary

The Rehearse CLI (`rehearse`) is the primary interface for running and managing markdown-native test scenarios. It provides commands for executing scenarios, listing available tests, and producing structured output for both humans and CI pipelines.

## Design Principles

### Command hierarchy

Commands follow a `rehearse <action>` pattern:

```
rehearse run [path]       — run scenario file or directory
rehearse list             — list available scenarios
rehearse version          — print version information
```

### Exit code contract

All commands share a consistent exit code contract:

| Exit code | Meaning |
|---|---|
| `0` | Success (all scenarios/steps passed) |
| `1` | Test failure (one or more scenarios/steps failed) |
| `2` | Invalid arguments |
| `3` | Resource not found (scenario file or directory does not exist) |
| `10+` | Unexpected errors |

## Behavior

### `rehearse run`

Executes one or more test scenarios. Accepts a file path (single scenario) or directory path (all scenarios in directory, recursive).

```
rehearse run [path]                  — run scenario file or directory
rehearse run --tag e2e               — filter by tag
rehearse run --format json           — machine-readable output
rehearse run --run-manual-tests      — include scenarios tagged 'manual'
rehearse run --spec-root ./my-spec   — override spec root directory
```

| Flag | Default | Description |
|---|---|---|
| `--format` | `text` | Output format: `text` (styled with live progress) or `json` |
| `--spec-root` | `spec` | Override the spec root directory for AC resolution |
| `--tag` | | Filter scenarios by tag (repeatable) |
| `--run-manual-tests` | `false` | Include scenarios tagged `manual` in directory scans |

When no path is given, defaults to scanning the current directory.

Scenarios tagged `manual` are skipped during directory scans unless `--run-manual-tests` is set, but always run when a specific file path is passed directly.

### `rehearse list`

Lists available scenarios without executing them.

```
rehearse list              — list all scenarios in current directory
rehearse list --tag e2e    — list filtered by tag
```

Output includes scenario name, file path, tags, and step count.

### `rehearse version`

Prints the Rehearse version, Go version, and build metadata.

## Interaction with Other Features

| Feature | Interaction |
|---|---|
| [Testing Framework](../testing-framework/README.md) | CLI is the user-facing entry point for the testing framework. `run` and `list` delegate to the test runner. |
| [Test Runner](../testing-framework/test-runner/README.md) | `rehearse run` instantiates the runner, passes configuration, and formats output. |

## Acceptance Criteria

Not defined yet.

## Outstanding Questions

- Acceptance criteria not yet defined for this feature.
- Should `rehearse run` support a `--dry-run` flag that parses and validates scenarios without executing them?
- Should there be a `rehearse init` command to scaffold a new spec directory with example scenarios?
