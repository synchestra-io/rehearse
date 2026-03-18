# AC: validates-ac-structure

**Status:** planned
**Feature:** [cli/validate](../README.md)

## Description

Given a well-formed `.ac.md` file with a title, Status, Feature back-reference,
Description, Inputs table, and Verification section, `rehearse validate` exits 0.
Given a malformed AC file — such as a missing Status field, a slug that does not
match the filename, or a missing Verification code block when the status is
`implemented` — `rehearse validate` exits 1 with a specific error describing the
problem.

## Inputs

| Name | Required | Description |
|---|---|---|
| binary_path | Yes | Path to the compiled `rehearse` binary |
| valid_ac_path | Yes | Path to a well-formed `.ac.md` file |
| invalid_ac_path | Yes | Path to a malformed `.ac.md` file |

## Verification

```bash
# --- Case 1: valid AC passes ---
output=$("$binary_path" validate "$valid_ac_path" 2>&1)
rc=$?
test $rc -eq 0 || {
  echo "Expected exit 0 for valid AC, got $rc"
  echo "$output"
  exit 1
}

# --- Case 2: malformed AC fails with descriptive error ---
output=$("$binary_path" validate "$invalid_ac_path" 2>&1)
rc=$?
test $rc -eq 1 || {
  echo "Expected exit 1 for invalid AC, got $rc"
  echo "$output"
  exit 1
}

# Error should reference the file path
echo "$output" | grep -q "$invalid_ac_path" || {
  echo "Error output does not mention file path '$invalid_ac_path'"
  echo "$output"
  exit 1
}

# Error should describe the structural problem
echo "$output" | grep -qiE '(status|slug|verification|missing)' || {
  echo "Error output does not describe the structural issue"
  echo "$output"
  exit 1
}
```

## Scenarios

(None yet.)
