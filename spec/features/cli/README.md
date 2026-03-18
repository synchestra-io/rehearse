# Feature: CLI

**Status:** In Progress

## Summary

The Rehearse CLI (`rehearse`) is the primary interface for running and managing markdown-native test scenarios. It provides commands for executing scenarios, listing available tests, validating spec artifacts, and producing structured output for both humans and CI pipelines.

## Contents

| Command | Description |
|---|---|
| [run](run/README.md) | Execute test scenario files or directories |
| [list](list/README.md) | List available test scenarios without executing them |
| [validate](validate/README.md) | Check structural validity of scenarios and ACs without execution |
| [version](version/README.md) | Print version information |

### run

Executes one or more test scenarios. Accepts a file path (single scenario) or directory path (all scenarios, recursive). Supports tag filtering, JSON output, manual scenario inclusion, and configurable spec root. See [run/README.md](run/README.md).

### list

Lists available scenarios without executing them. Discovers scenarios from `spec/tests/` and `spec/features/*/_tests/` directories. Supports tag filtering. See [list/README.md](list/README.md).

### validate

Checks the structural validity of test scenarios and acceptance criteria files without executing any code. Catches format errors, missing references, and orphaned files before runtime. See [validate/README.md](validate/README.md).

### version

Prints the Rehearse version, Go version, and build metadata. See [version/README.md](version/README.md).

## Design Principles

### Command hierarchy

Commands follow a `rehearse <action>` pattern:

```
rehearse run [path]       — run scenario file or directory
rehearse list             — list available scenarios
rehearse validate [path]  — validate scenarios and ACs
rehearse version          — print version information
```

### Exit code contract

All commands share a consistent exit code contract:

| Exit code | Meaning |
|---|---|
| `0` | Success |
| `1` | Failure (test failure for `run`, validation errors for `validate`) |
| `2` | Invalid arguments |
| `3` | Resource not found (path does not exist) |
| `10+` | Unexpected errors |

Individual commands may use a subset of these codes. See each command's spec for details.

## Interaction with Other Features

| Feature | Interaction |
|---|---|
| [Testing Framework](../testing-framework/README.md) | CLI is the user-facing entry point for the testing framework. `run` and `list` delegate to the test runner. |
| [Test Runner](../testing-framework/test-runner/README.md) | `rehearse run` instantiates the runner; `rehearse validate` shares parsing logic. |
| [Acceptance Criteria](../acceptance-criteria/README.md) | `rehearse validate` checks AC file structure and cross-references. |

## Acceptance Criteria

Acceptance criteria are defined at the per-command level. See each command's spec for its ACs:

- [cli/run ACs](run/_acs/README.md) — 12 ACs
- [cli/list ACs](list/_acs/README.md) — 2 ACs
- [cli/validate ACs](validate/_acs/README.md) — 8 ACs
- [cli/version ACs](version/_acs/README.md) — 1 AC

## Outstanding Questions

None at this time.
