# AC: overrides-spec-root

**Status:** planned
**Feature:** [cli/run](../README.md)

## Description

The `--spec-root` flag changes the root directory used for resolving AC
references in scenarios. When provided, all AC paths are resolved relative to
the custom root instead of the default spec directory.

## Inputs

| Name | Required | Description |
|---|---|---|
| binary_path | Yes | Path to the compiled `rehearse` binary |
| scenario_path | Yes | Path to a scenario `.test.md` file with AC references |
| custom_spec_root | Yes | Alternative directory to use as the spec root |

## Verification

```bash
# Run with the custom spec root
output=$("$binary_path" run "$scenario_path" --spec-root "$custom_spec_root" --format json 2>&1)
rc=$?
test $rc -eq 0 || { echo "Expected exit 0, got $rc"; echo "$output"; exit 1; }

# Verify ACs were resolved — none should be unresolved
unresolved=$(echo "$output" | jq '[.steps[]
  | select(.ac_refs != null)
  | .ac_refs[]
  | select(.resolved == false)] | length')
test "$unresolved" -eq 0 \
  || { echo "Found $unresolved unresolved AC refs with custom spec root"; echo "$output"; exit 1; }

# Verify the spec root recorded in output matches the override
reported_root=$(echo "$output" | jq -r '.spec_root // empty')
test "$reported_root" = "$custom_spec_root" \
  || { echo "Spec root mismatch: expected '$custom_spec_root', got '$reported_root'"; exit 1; }
```

## Scenarios

(None yet.)
