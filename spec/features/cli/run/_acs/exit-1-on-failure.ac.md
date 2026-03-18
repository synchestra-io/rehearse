# AC: exit-1-on-failure

**Status:** planned
**Feature:** [cli/run](../README.md)

## Description

When any scenario fails, `rehearse run` exits with code 1, indicating a test
failure occurred.

## Inputs

| Name | Required | Description |
|---|---|---|
| binary_path | Yes | Path to the compiled `rehearse` binary |
| scenario_path | Yes | Path to a scenario `.test.md` file that is expected to fail |

## Verification

```bash
"$binary_path" run "$scenario_path" > /dev/null 2>&1
rc=$?
test $rc -eq 1 \
  || { echo "Expected exit code 1, got $rc"; exit 1; }
```

## Scenarios

(None yet.)
