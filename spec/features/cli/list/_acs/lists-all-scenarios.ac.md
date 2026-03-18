# AC: lists-all-scenarios

**Status:** planned
**Feature:** [cli/list](../README.md)

## Description

Given a spec root containing scenario files in `spec/tests/` and
`spec/features/*/_tests/`, `rehearse list` discovers all `.test.md` files and
outputs every scenario's name, description, and tags in columnar format. No
scenarios are omitted and no extra rows appear.

## Inputs

| Name | Required | Description |
|---|---|---|
| binary_path | Yes | Path to the compiled `rehearse` binary |
| spec_root | Yes | Path to the spec root directory containing test scenarios |

## Verification

```bash
output=$("$binary_path" list 2>&1)
rc=$?
test $rc -eq 0 || { echo "Expected exit 0, got $rc"; echo "$output"; exit 1; }

# Verify header row is present with expected columns
echo "$output" | head -1 | grep -q "SCENARIO" || { echo "Missing SCENARIO column header"; exit 1; }
echo "$output" | head -1 | grep -q "DESCRIPTION" || { echo "Missing DESCRIPTION column header"; exit 1; }
echo "$output" | head -1 | grep -q "TAGS" || { echo "Missing TAGS column header"; exit 1; }

# Verify at least one scenario row appears (header + 1 data row minimum)
row_count=$(echo "$output" | wc -l | tr -d ' ')
test "$row_count" -ge 2 || { echo "Expected at least 2 lines (header + data), got $row_count"; exit 1; }

# Verify known scenario names from spec/tests/ and spec/features/*/_tests/ appear
# (scenario names are derived from .test.md files present under spec_root)
find "$spec_root/spec/tests" "$spec_root/spec/features" -name '*.test.md' 2>/dev/null | while read -r f; do
  name=$(head -20 "$f" | grep -m1 '^# Scenario:' | sed 's/^# Scenario:[[:space:]]*//')
  if [ -n "$name" ]; then
    echo "$output" | grep -q "$name" || { echo "Scenario '$name' from $f not found in output"; exit 1; }
  fi
done

echo "PASS: all discovered scenarios listed with columnar format"
```

## Scenarios

(None yet.)
