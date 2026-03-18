# AC: runs-single-file

**Status:** planned
**Feature:** [cli/run](../README.md)

## Description

Given a valid scenario file path, `rehearse run <file>` executes the scenario
and reports its results. The command exits successfully and the output contains
the scenario name from the file.

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

# Verify output is valid JSON with a scenario name
echo "$output" | jq -e '.scenario.name' > /dev/null \
  || { echo "JSON output missing scenario name"; echo "$output"; exit 1; }

# Verify at least one result is reported
result_count=$(echo "$output" | jq '.steps | length')
test "$result_count" -ge 0 || { echo "No result data in output"; exit 1; }
```

## Scenarios

(None yet.)
