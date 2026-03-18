# AC: filters-by-tag

**Status:** planned
**Feature:** [cli/list](../README.md)

## Description

Given scenarios tagged with different tags, running `rehearse list --tag <tag>`
returns only the scenarios whose Tags field includes the specified tag. Scenarios
that do not carry the tag are absent from the output.

## Inputs

| Name | Required | Description |
|---|---|---|
| binary_path | Yes | Path to the compiled `rehearse` binary |
| spec_root | Yes | Path to the spec root directory containing test scenarios |
| filter_tag | Yes | Tag value to filter by (e.g. `e2e`) |

## Verification

```bash
# Run unfiltered list to capture all scenario names
all_output=$("$binary_path" list 2>&1)
all_rc=$?
test $all_rc -eq 0 || { echo "Unfiltered list failed with exit $all_rc"; echo "$all_output"; exit 1; }

# Run filtered list
filtered_output=$("$binary_path" list --tag "$filter_tag" 2>&1)
filtered_rc=$?
test $filtered_rc -eq 0 || { echo "Filtered list failed with exit $filtered_rc"; echo "$filtered_output"; exit 1; }

# Every row in the filtered output (excluding the header) must contain the filter tag
echo "$filtered_output" | tail -n +2 | while read -r line; do
  [ -z "$line" ] && continue
  echo "$line" | grep -qi "$filter_tag" || { echo "Row missing tag '$filter_tag': $line"; exit 1; }
done

# At least one matching scenario should appear (test assumes the tag exists)
data_rows=$(echo "$filtered_output" | tail -n +2 | grep -c '[^[:space:]]')
test "$data_rows" -ge 1 || { echo "No scenarios matched tag '$filter_tag'"; exit 1; }

# Scenarios without the tag must be absent — compare against unfiltered list
echo "$all_output" | tail -n +2 | while read -r line; do
  [ -z "$line" ] && continue
  echo "$line" | grep -qi "$filter_tag" && continue
  # Extract the scenario name (first column)
  scenario_name=$(echo "$line" | awk '{print $1}')
  echo "$filtered_output" | grep -q "$scenario_name" && { echo "Non-matching scenario '$scenario_name' present in filtered output"; exit 1; }
done

echo "PASS: --tag filter returns only matching scenarios"
```

## Scenarios

(None yet.)
