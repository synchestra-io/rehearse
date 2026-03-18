# AC: fail-fast-stops-early

**Status:** planned
**Feature:** [cli/validate](../README.md)

## Description

When `--fail-fast` or `--fail-fast=N` is passed, the validate command stops collecting errors once the limit is reached. The output includes only the errors found before the limit and a truncation note. The exit code is still `1` (validation errors found).

Without `--fail-fast`, the validate command collects all errors across all files. With `--fail-fast` (no value), the limit defaults to 1. With `--fail-fast=N`, the limit is N.

## Inputs

| Name | Required | Description |
|---|---|---|
| STEP_STDOUT | yes | Standard output from the validate command |
| STEP_EXIT_CODE | yes | Exit code from the validate command |

## Verification

```bash
# Verify that the output was truncated
echo "$STEP_STDOUT" | grep -q "fail-fast" || {
  echo "Expected truncation note in output"
  exit 1
}

# Verify exit code is 1 (validation errors)
test "$STEP_EXIT_CODE" -eq 1 || {
  echo "Expected exit code 1, got $STEP_EXIT_CODE"
  exit 1
}
```

## Scenarios

- [validate-core](../_tests/validate-core.test.md): `validate-fail-fast`
