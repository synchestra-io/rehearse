# AC: exit-1-validation-errors

**Status:** planned
**Feature:** [cli/validate](../README.md)

## Description

When any checked file has structural errors, `rehearse validate` exits 1 and prints
each error with its file path, line number (where applicable), and a description of
the problem. The output ends with a summary line listing the count of files and errors
found.

## Inputs

| Name | Required | Description |
|---|---|---|
| binary_path | Yes | Path to the compiled `rehearse` binary |
| invalid_dir | Yes | Path to a directory containing at least one structurally invalid file |

## Verification

```bash
output=$("$binary_path" validate "$invalid_dir" 2>&1)
rc=$?
test $rc -eq 1 || {
  echo "Expected exit 1 for directory with errors, got $rc"
  echo "$output"
  exit 1
}

# Output should list at least one error
echo "$output" | grep -qiE '(error|invalid|missing|malformed)' || {
  echo "Expected error descriptions in output"
  echo "$output"
  exit 1
}

# Output should include a summary line with counts
echo "$output" | grep -qE '[0-9]+ (file|error)' || {
  echo "Expected summary line with file/error counts"
  echo "$output"
  exit 1
}
```

## Scenarios

(None yet.)
