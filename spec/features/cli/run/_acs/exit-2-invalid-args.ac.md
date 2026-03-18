# AC: exit-2-invalid-args

**Status:** planned
**Feature:** [cli/run](../README.md)

## Description

When `rehearse run` is invoked with invalid arguments (e.g. unknown flags or
missing required arguments), it exits with code 2 and prints a usage error.

## Inputs

| Name | Required | Description |
|---|---|---|
| binary_path | Yes | Path to the compiled `rehearse` binary |

## Verification

```bash
output=$("$binary_path" run --nonexistent-flag 2>&1)
rc=$?
test $rc -eq 2 \
  || { echo "Expected exit code 2, got $rc"; echo "$output"; exit 1; }

# Error message should mention the bad flag or show usage
echo "$output" | grep -qi 'unknown\|invalid\|usage\|unrecognized' \
  || { echo "Expected usage/error message in output"; echo "$output"; exit 1; }
```

## Scenarios

(None yet.)
