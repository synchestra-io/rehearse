# AC: validates-scenario-structure

**Status:** planned
**Feature:** [cli/validate](../README.md)

## Description

Given a well-formed `.test.md` file that contains a title, description, and properly
annotated steps, `rehearse validate` exits 0. Given a malformed scenario — such as a
missing title, a bare code fence without a language annotation, or duplicate step
names — `rehearse validate` exits 1 and reports the specific error with the file path
and line number where the problem was detected.

## Inputs

| Name | Required | Description |
|---|---|---|
| binary_path | Yes | Path to the compiled `rehearse` binary |
| valid_scenario_path | Yes | Path to a well-formed `.test.md` file |
| invalid_scenario_path | Yes | Path to a malformed `.test.md` file |

## Verification

```bash
# --- Case 1: valid scenario passes ---
output=$("$binary_path" validate "$valid_scenario_path" 2>&1)
rc=$?
test $rc -eq 0 || {
  echo "Expected exit 0 for valid scenario, got $rc"
  echo "$output"
  exit 1
}

# --- Case 2: malformed scenario fails with descriptive error ---
output=$("$binary_path" validate "$invalid_scenario_path" 2>&1)
rc=$?
test $rc -eq 1 || {
  echo "Expected exit 1 for invalid scenario, got $rc"
  echo "$output"
  exit 1
}

# Error should reference the file path
echo "$output" | grep -q "$invalid_scenario_path" || {
  echo "Error output does not mention file path '$invalid_scenario_path'"
  echo "$output"
  exit 1
}

# Error should include a line number
echo "$output" | grep -qE 'line [0-9]+' || {
  echo "Error output does not include a line number"
  echo "$output"
  exit 1
}
```

## Scenarios

(None yet.)
