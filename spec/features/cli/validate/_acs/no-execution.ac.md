# AC: no-execution

**Status:** planned
**Feature:** [cli/validate](../README.md)

## Description

`rehearse validate` must not execute any verification scripts or test step code
blocks. Given a scenario with a step that writes a marker file to disk, after running
`rehearse validate` on that scenario the marker file must NOT exist. This proves
validation is purely structural — it reads markdown, not bash.

## Inputs

| Name | Required | Description |
|---|---|---|
| binary_path | Yes | Path to the compiled `rehearse` binary |
| scenario_path | Yes | Path to a scenario whose steps would create a marker file if executed |
| marker_file | Yes | Path where the step code block would write a file |

## Verification

```bash
# Ensure marker does not already exist
rm -f "$marker_file"

# Run validate (not run) on the scenario
output=$("$binary_path" validate "$scenario_path" 2>&1)

# Marker file must NOT exist — validate should not execute code blocks
if [ -f "$marker_file" ]; then
  echo "FAIL: marker file '$marker_file' exists — validate executed step code"
  rm -f "$marker_file"
  exit 1
fi
```

## Scenarios

(None yet.)
