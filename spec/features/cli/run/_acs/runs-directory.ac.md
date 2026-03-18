# AC: runs-directory

**Status:** planned
**Feature:** [cli/run](../README.md)

## Description

Given a directory containing multiple `.test.md` files, `rehearse run <dir>`
discovers and runs all of them. The command exits successfully and the output
reports results for every discovered scenario.

## Inputs

| Name | Required | Description |
|---|---|---|
| binary_path | Yes | Path to the compiled `rehearse` binary |
| dir_path | Yes | Path to a directory containing multiple `.test.md` scenario files |

## Verification

```bash
output=$("$binary_path" run "$dir_path" --format json 2>&1)
rc=$?
test $rc -eq 0 || { echo "Expected exit 0, got $rc"; echo "$output"; exit 1; }

# Count reported scenarios — expect more than one
scenario_count=$(echo "$output" | jq '.scenarios | length')
test "$scenario_count" -gt 1 \
  || { echo "Expected multiple scenarios, got $scenario_count"; echo "$output"; exit 1; }

# Each scenario must have a name
unnamed=$(echo "$output" | jq '[.scenarios[] | select(.name == null or .name == "")] | length')
test "$unnamed" -eq 0 \
  || { echo "Found $unnamed scenarios without a name"; exit 1; }
```

## Scenarios

(None yet.)
