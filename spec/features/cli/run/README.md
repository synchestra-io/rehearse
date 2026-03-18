# Feature: CLI — `rehearse run`

**Status:** In Progress

## Summary

`rehearse run` executes one or more test scenarios. Accepts a file path (single scenario) or directory path (all scenarios, recursive). It is the primary command for running [test scenarios](../../testing-framework/test-scenario/README.md) — parsing them, resolving [acceptance criteria](../../acceptance-criteria/README.md), delegating execution to the [test runner](../../testing-framework/test-runner/README.md), and reporting structured results.

## Behavior

### Usage

```
rehearse run [path]                  — run scenario file or directory
rehearse run --tag e2e               — filter by tag
rehearse run --format json           — machine-readable output
rehearse run --run-manual-tests      — include scenarios tagged 'manual'
rehearse run --spec-root ./my-spec   — override spec root directory
```

When no path is given, defaults to scanning the current directory.

### Flags

| Flag | Default | Description |
|---|---|---|
| `--format` | `text` | Output format: `text` (styled with live progress) or `json` |
| `--spec-root` | `spec` | Override the spec root directory for AC resolution |
| `--tag` | | Filter scenarios by tag (repeatable) |
| `--run-manual-tests` | `false` | Include scenarios tagged `manual` in directory scans |

### Path resolution

- **File path:** Runs the single scenario at the given path. The scenario is always executed, regardless of tags — including `manual`.
- **Directory path:** Recursively discovers all `*.test.md` files under the directory and runs them. Tag filters and manual-test exclusion apply.
- **No path:** Equivalent to passing `.` — scans the current directory recursively.

### Manual scenario filtering

Scenarios tagged `manual` are skipped during directory scans unless `--run-manual-tests` is set. This prevents demo scenarios, stress tests, and interactive verification from running in CI or routine sweeps.

Manual scenarios run when:
- The scenario file is **passed directly by path** (e.g., `rehearse run path/to/demo.test.md`)
- The `--run-manual-tests` flag is set (e.g., `rehearse run tests/ --run-manual-tests`)

### Tag filtering

The `--tag` flag filters scenarios by their declared tags. Only scenarios that include at least one of the specified tags are executed. The flag is repeatable:

```
rehearse run --tag e2e --tag smoke
```

This runs scenarios tagged `e2e` **or** `smoke`.

### Output formats

- **`text`** (default): Human-readable, colored terminal output with real-time progress. Uses lipgloss for styled output with checkmarks/crosses and inline duration. Each step shows a dimmed `▸ step-name` indicator while running, replaced in-place with a colored `✔`/`✘` result when the step completes.
- **`json`**: Machine-readable output for CI integration, dashboards, and programmatic analysis. Emitted as a single JSON object after all scenarios complete.

### Spec root resolution

The `--spec-root` flag (default: `spec`) determines where [acceptance criteria](../../acceptance-criteria/README.md) are resolved from. All AC references resolve relative to this root — `cli/project/remove/*` becomes `{spec_root}/features/cli/project/remove/_acs/`. The spec root is read once at initialization and threaded through the AC resolver.

## Exit Code Contract

All exit codes follow the shared [CLI exit code contract](../README.md#exit-code-contract):

| Exit code | Meaning |
|---|---|
| `0` | Success (all scenarios/steps passed) |
| `1` | Test failure (one or more scenarios/steps failed) |
| `2` | Invalid arguments |
| `3` | Resource not found (scenario file or directory does not exist) |
| `10+` | Unexpected errors |

## Interaction with Other Features

| Feature | Interaction |
|---|---|
| [Testing Framework](../../testing-framework/README.md) | CLI is the user-facing entry point for the testing framework. |
| [Test Runner](../../testing-framework/test-runner/README.md) | `rehearse run` instantiates the runner, passing parsed flags and resolved paths. |
| [Acceptance Criteria](../../acceptance-criteria/README.md) | AC resolution during runs — the runner resolves `_acs/` references relative to `--spec-root`. |

## Acceptance Criteria

| AC | Description | Status |
|---|---|---|
| [runs-single-file](_acs/runs-single-file.ac.md) | Single scenario file executed and results reported | planned |
| [runs-directory](_acs/runs-directory.ac.md) | All scenarios in directory discovered and run recursively | planned |
| [filters-by-tag](_acs/filters-by-tag.ac.md) | --tag flag filters to matching scenarios only | planned |
| [json-output](_acs/json-output.ac.md) | --format json produces valid structured JSON output | planned |
| [skips-manual-in-scan](_acs/skips-manual-in-scan.ac.md) | Manual-tagged scenarios skipped in directory scans | planned |
| [runs-manual-when-direct](_acs/runs-manual-when-direct.ac.md) | Manual scenarios run when file path passed directly | planned |
| [runs-manual-with-flag](_acs/runs-manual-with-flag.ac.md) | --run-manual-tests includes manual scenarios in scans | planned |
| [overrides-spec-root](_acs/overrides-spec-root.ac.md) | --spec-root changes AC resolution root | planned |
| [exit-0-on-pass](_acs/exit-0-on-pass.ac.md) | Exits 0 when all scenarios pass | planned |
| [exit-1-on-failure](_acs/exit-1-on-failure.ac.md) | Exits 1 when any scenario fails | planned |
| [exit-2-invalid-args](_acs/exit-2-invalid-args.ac.md) | Exits 2 on invalid arguments | planned |
| [exit-3-not-found](_acs/exit-3-not-found.ac.md) | Exits 3 when path does not exist | planned |

## Outstanding Questions

None at this time.
