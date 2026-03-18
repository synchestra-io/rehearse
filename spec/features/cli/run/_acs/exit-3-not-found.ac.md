# AC: exit-3-not-found

**Status:** planned
**Feature:** [cli/run](../README.md)

## Description

When `rehearse run` is given a path that does not exist, it exits with code 3
and prints an error indicating the path was not found.

## Inputs

| Name | Required | Description |
|---|---|---|
| binary_path | Yes | Path to the compiled `rehearse` binary |

## Verification

```bash
output=$("$binary_path" run "/tmp/does-not-exist-$$" 2>&1)
rc=$?
test $rc -eq 3 \
  || { echo "Expected exit code 3, got $rc"; echo "$output"; exit 1; }

# Error message should reference the missing path
echo "$output" | grep -qi 'not found\|no such\|does not exist' \
  || { echo "Expected 'not found' message in output"; echo "$output"; exit 1; }
```

## Scenarios

(None yet.)
