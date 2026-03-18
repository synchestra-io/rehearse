# AC: exit-0-on-pass

**Status:** planned
**Feature:** [cli/run](../README.md)

## Description

When all scenarios pass, `rehearse run` exits with code 0, indicating success.

## Inputs

| Name | Required | Description |
|---|---|---|
| binary_path | Yes | Path to the compiled `rehearse` binary |
| scenario_path | Yes | Path to a scenario `.test.md` file that is expected to pass |

## Verification

```bash
"$binary_path" run "$scenario_path" > /dev/null 2>&1
rc=$?
test $rc -eq 0 \
  || { echo "Expected exit code 0, got $rc"; exit 1; }
```

## Scenarios

(None yet.)
