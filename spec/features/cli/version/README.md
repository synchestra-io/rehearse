# Feature: CLI — version

**Status:** In Progress

## Summary

`rehearse version` prints the Rehearse version, Go version, and build metadata. It is a standalone informational command with no flags and no side effects.

## Behavior

### Usage

```
rehearse version
```

### Output format

```
rehearse {version} ({commit}) {date}
```

Example:

```
rehearse v0.1.0 (a1b2c3d) 2025-03-17
```

The values are injected at build time via `-ldflags`. When not set, the output uses `(unknown)` for missing fields.

### Flags

None.

### Exit code contract

| Exit code | Meaning |
|---|---|
| `0` | Success |
| `10+` | Unexpected errors |

The command always exits `0` under normal operation. Exit codes `10+` are reserved for unexpected errors (e.g., internal failures writing to stdout).

## Interaction with Other Features

None — standalone informational command.

## Acceptance Criteria

| AC | Description | Status |
|---|---|---|
| [prints-version-info](_acs/prints-version-info.ac.md) | Prints version, commit, and date in expected format | planned |

## Outstanding Questions

None at this time.
