# AC: filters-by-tag

**Status:** planned
**Feature:** [cli/run](../README.md)

## Description

Given scenarios with different tags, running `rehearse run <dir> --tag <tag>`
only executes the scenarios that carry the matching tag. Scenarios without the
tag are excluded from the results.

## Inputs

| Name | Required | Description |
|---|---|---|
| binary_path | Yes | Path to the compiled `rehearse` binary |
| dir_path | Yes | Path to a directory containing scenarios with various tags |
| filter_tag | Yes | Tag value to filter on (e.g. `smoke`) |

## Verification

```bash
output=$("$binary_path" run "$dir_path" --tag "$filter_tag" --format json 2>&1)
rc=$?
test $rc -eq 0 || { echo "Expected exit 0, got $rc"; echo "$output"; exit 1; }

# Every reported scenario must carry the filter tag
mismatched=$(echo "$output" | jq --arg tag "$filter_tag" \
  '[.scenarios[] | select(.tags | index($tag) | not)] | length')
test "$mismatched" -eq 0 \
  || { echo "Found $mismatched scenarios missing tag '$filter_tag'"; echo "$output"; exit 1; }

# At least one scenario should have run
count=$(echo "$output" | jq '.scenarios | length')
test "$count" -gt 0 \
  || { echo "No scenarios matched tag '$filter_tag'"; exit 1; }
```

## Scenarios

(None yet.)
