# AC: prints-version-info

**Status:** planned
**Feature:** [cli/version](../README.md)

## Description

Running `rehearse version` prints a single line matching the format
`rehearse {version} ({commit}) {date}` and exits with code 0.
The version, commit, and date values are injected at build time via `-ldflags`.
When not set, the output uses `(unknown)` for missing fields.

## Inputs

| Name | Required | Description |
|---|---|---|
| binary_path | Yes | Path to the compiled `rehearse` binary |

## Verification

```bash
output=$("$binary_path" version 2>&1)
rc=$?
test $rc -eq 0 || { echo "Expected exit 0, got $rc"; echo "$output"; exit 1; }

# Verify output matches expected format: rehearse {version} ({commit}) {date}
# version: v-prefixed semver or (unknown)
# commit: short hash or (unknown)
# date: YYYY-MM-DD or (unknown)
pattern='^rehearse (v[0-9]+\.[0-9]+\.[0-9]+[^ ]*|\(unknown\)) \(([a-f0-9]+|\(unknown\))\) ([0-9]{4}-[0-9]{2}-[0-9]{2}|\(unknown\))$'
echo "$output" | grep -Eq "$pattern" || { echo "Output does not match expected format"; echo "Expected pattern: rehearse {version} ({commit}) {date}"; echo "Got: $output"; exit 1; }

echo "PASS: output matches expected format"
```

## Scenarios

(None yet.)
