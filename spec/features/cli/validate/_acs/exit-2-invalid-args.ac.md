# AC: exit-2-invalid-args

**Status:** planned
**Feature:** [cli/validate](../README.md)

## Description

When `rehearse validate` is invoked with invalid arguments — such as an unknown flag
— it exits 2 without performing any validation. This distinguishes argument errors
from validation errors (exit 1) and successful runs (exit 0).

## Inputs

| Name | Required | Description |
|---|---|---|
| binary_path | Yes | Path to the compiled `rehearse` binary |

## Verification

```bash
# Pass an unknown flag to trigger an argument error
output=$("$binary_path" validate --no-such-flag 2>&1)
rc=$?
test $rc -eq 2 || {
  echo "Expected exit 2 for invalid arguments, got $rc"
  echo "$output"
  exit 1
}
```

## Scenarios

(None yet.)
