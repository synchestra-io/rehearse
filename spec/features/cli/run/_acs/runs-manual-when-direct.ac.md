# AC: runs-manual-when-direct

**Status:** planned
**Feature:** [cli/run](../README.md)

## Description

A scenario tagged `manual` is executed normally when its file path is passed
directly to `rehearse run`. The manual skip behavior only applies to directory
scans, not to explicit file paths.

## Inputs

| Name | Required | Description |
|---|---|---|
| binary_path | Yes | Path to the compiled `rehearse` binary |
| scenario_path | Yes | Path to a `manual`-tagged scenario `.test.md` file |

## Verification

```bash
output=$("$binary_path" run "$scenario_path" --format json 2>&1)
rc=$?
test $rc -eq 0 || { echo "Expected exit 0, got $rc"; echo "$output"; exit 1; }

# Scenario must actually have been executed (not skipped)
status=$(echo "$output" | jq -r '.status')
test "$status" != "skipped" \
  || { echo "Manual scenario was skipped when given directly"; echo "$output"; exit 1; }

# Verify the scenario name is present
echo "$output" | jq -e '.scenario.name' > /dev/null \
  || { echo "Scenario name missing from output"; exit 1; }
```

## Scenarios

(None yet.)
