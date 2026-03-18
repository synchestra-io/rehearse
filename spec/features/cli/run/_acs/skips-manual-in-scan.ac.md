# AC: skips-manual-in-scan

**Status:** planned
**Feature:** [cli/run](../README.md)

## Description

Scenarios tagged `manual` are automatically skipped when `rehearse run` is
invoked with a directory path. They do not appear in the results output.

## Inputs

| Name | Required | Description |
|---|---|---|
| binary_path | Yes | Path to the compiled `rehearse` binary |
| dir_path | Yes | Path to a directory containing at least one `manual`-tagged scenario |

## Verification

```bash
output=$("$binary_path" run "$dir_path" --format json 2>&1)
rc=$?
test $rc -eq 0 || { echo "Expected exit 0, got $rc"; echo "$output"; exit 1; }

# No scenario in the results should carry the "manual" tag
manual_count=$(echo "$output" | jq \
  '[.scenarios[] | select(.tags | index("manual"))] | length')
test "$manual_count" -eq 0 \
  || { echo "Expected 0 manual scenarios in results, got $manual_count"; echo "$output"; exit 1; }
```

## Scenarios

(None yet.)
