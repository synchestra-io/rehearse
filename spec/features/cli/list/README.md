# Feature: CLI — list

**Status:** In Progress

## Summary

`rehearse list` lists available test scenarios without executing them. It discovers scenarios from `spec/tests/` and `spec/features/*/_tests/` directories, and prints a summary table showing each scenario's name, description, and tags. Supports tag-based filtering to narrow results.

## Behavior

### Usage

```
rehearse list              — list all scenarios in current directory
rehearse list --tag e2e    — list filtered by tag
```

### Flags

| Flag | Default | Description |
|---|---|---|
| `--tag` | | Filter scenarios by tag (repeatable) |

### Output format

Columnar text with the following columns:

```
SCENARIO          DESCRIPTION                          TAGS
runner-core       Core test runner lifecycle            e2e, dogfood
progress-demo     Live progress indicator demo          manual
```

Each row includes:
- **Scenario name** — extracted from the `# Scenario:` heading
- **Description** — the scenario's Description metadata field
- **Tags** — comma-separated list of tags

### Discovery

The command scans two directory trees for scenario files:

- **`spec/tests/`** — cross-feature scenarios that verify end-to-end workflows spanning multiple features
- **`spec/features/*/_tests/`** — feature-scoped scenarios that verify a single feature's behavior

Scenario files are identified by the `.test.md` extension. Files without this extension are ignored.

### Exit code contract

| Exit code | Meaning |
|---|---|
| `0` | Success |
| `2` | Invalid arguments |
| `10+` | Unexpected errors |

## Interaction with Other Features

| Feature | Interaction |
|---|---|
| [Testing Framework](../../testing-framework/README.md) | `rehearse list` uses the testing framework's scenario parser to discover and read scenario metadata. |
| [Test Scenario](../../testing-framework/test-scenario/README.md) | Scenarios are defined by the test scenario format — `list` parses the `# Scenario:` heading, Description, and Tags fields. |

## Acceptance Criteria

| AC | Description | Status |
|---|---|---|
| [lists-all-scenarios](_acs/lists-all-scenarios.ac.md) | All discovered scenarios shown with name, description, and tags | planned |
| [filters-by-tag](_acs/filters-by-tag.ac.md) | --tag flag filters listed scenarios | planned |

## Outstanding Questions

None at this time.
