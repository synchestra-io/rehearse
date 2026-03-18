# AC: json-output

**Status:** planned
**Feature:** [cli/run](../README.md)

## Description

When `--format json` is passed, `rehearse run` produces valid JSON output
containing the scenario name, an ordered list of steps, and a pass/fail status
for each step and for the scenario overall.

## Inputs

| Name | Required | Description |
|---|---|---|
| binary_path | Yes | Path to the compiled `rehearse` binary |
| scenario_path | Yes | Path to a valid scenario `.test.md` file |

## Verification

```bash
output=$("$binary_path" run "$scenario_path" --format json 2>&1)
rc=$?
test $rc -eq 0 || { echo "Expected exit 0, got $rc"; echo "$output"; exit 1; }

# Must be valid JSON
echo "$output" | jq empty 2>/dev/null \
  || { echo "Output is not valid JSON"; echo "$output"; exit 1; }

# Required top-level fields
for field in scenario steps status; do
  echo "$output" | jq -e ".$field" > /dev/null \
    || { echo "Missing required field: $field"; exit 1; }
done

# Scenario must have a name
echo "$output" | jq -e '.scenario.name | length > 0' > /dev/null \
  || { echo "Scenario name is empty"; exit 1; }

# Each step must have a name and status
bad_steps=$(echo "$output" | jq '[.steps[] | select(.name == null or .status == null)] | length')
test "$bad_steps" -eq 0 \
  || { echo "Found $bad_steps steps missing name or status"; exit 1; }
```

## Scenarios

(None yet.)
