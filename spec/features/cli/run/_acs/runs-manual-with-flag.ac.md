# AC: runs-manual-with-flag

**Status:** planned
**Feature:** [cli/run](../README.md)

## Description

When `--run-manual-tests` is passed, `rehearse run <dir>` includes scenarios
tagged `manual` in the directory scan. The manual-tagged scenarios appear in
the results alongside non-manual ones.

## Inputs

| Name | Required | Description |
|---|---|---|
| binary_path | Yes | Path to the compiled `rehearse` binary |
| dir_path | Yes | Path to a directory containing at least one `manual`-tagged scenario |

## Verification

```bash
output=$("$binary_path" run "$dir_path" --run-manual-tests --format json 2>&1)
rc=$?
test $rc -eq 0 || { echo "Expected exit 0, got $rc"; echo "$output"; exit 1; }

# At least one manual-tagged scenario must be present
manual_count=$(echo "$output" | jq \
  '[.scenarios[] | select(.tags | index("manual"))] | length')
test "$manual_count" -gt 0 \
  || { echo "No manual scenarios found in results despite --run-manual-tests"; echo "$output"; exit 1; }
```

## Scenarios

(None yet.)
