# AC: validates-ac-refs-resolve

**Status:** planned
**Feature:** [cli/validate](../README.md)

## Description

Given a scenario that references an AC that exists on disk, validation passes. Given
a scenario that references a non-existent AC, validation fails with an error naming
the unresolved reference and the file path where the broken reference was found.

## Inputs

| Name | Required | Description |
|---|---|---|
| binary_path | Yes | Path to the compiled `rehearse` binary |
| scenario_path | Yes | Path to a scenario `.test.md` file that references a non-existent AC |
| spec_root | Yes | Path to the spec root directory used for resolving AC references |

## Verification

```bash
# Validate with --spec-root so references resolve against the fixture tree
output=$("$binary_path" validate "$scenario_path" --spec-root "$spec_root" 2>&1)
rc=$?
test $rc -eq 1 || {
  echo "Expected exit 1 for unresolvable AC reference, got $rc"
  echo "$output"
  exit 1
}

# Error should mention the missing AC path or slug
echo "$output" | grep -qiE '(unresolved|not found|missing|does not exist)' || {
  echo "Error output does not mention the unresolved reference"
  echo "$output"
  exit 1
}

# Error should name the scenario file where the bad reference lives
echo "$output" | grep -q "$scenario_path" || {
  echo "Error output does not mention the scenario file path"
  echo "$output"
  exit 1
}
```

## Scenarios

(None yet.)
