# AC: exit-0-all-valid

**Status:** planned
**Feature:** [cli/validate](../README.md)

## Description

When all checked files are structurally valid — scenarios have proper titles, code
blocks have language annotations, AC files have required sections, references resolve,
and indexes are in sync — `rehearse validate` exits 0 and prints a summary line.

## Inputs

| Name | Required | Description |
|---|---|---|
| binary_path | Yes | Path to the compiled `rehearse` binary |
| valid_dir | Yes | Path to a directory containing only well-formed scenario and AC files |

## Verification

```bash
output=$("$binary_path" validate "$valid_dir" 2>&1)
rc=$?
test $rc -eq 0 || {
  echo "Expected exit 0 for all-valid directory, got $rc"
  echo "$output"
  exit 1
}

# Output should contain a summary (e.g. "Validated X scenarios, Y ACs — no errors.")
echo "$output" | grep -qiE '(no errors|valid|pass)' || {
  echo "Expected summary line indicating success"
  echo "$output"
  exit 1
}
```

## Scenarios

(None yet.)
