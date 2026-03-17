# AC: substitutes-context-vars

**Status:** planned
**Feature:** [testing-framework/http-requests](../README.md)

## Description

Before sending an HTTP request, the runner substitutes all `${{ context.* }}` and `${{ env.* }}` references in the method line, headers, and body. Unresolved references cause the step to fail before the network call is made, with an error message identifying the missing variable.

## Inputs

| Name | Required | Description |
|---|---|---|
| scenario_path | Yes | Path to a scenario with an `http` block containing `${{ }}` references |
| context_var_name | Yes | Name of a context variable used in the request (without the `context.` prefix) |
| context_var_value | Yes | Value to set for that context variable |
| env_var_name | No | Name of an env var used in the request (without the `env.` prefix) |
| env_var_value | No | Value to set for that env var |
| expected_in_request | Yes | String that must appear in the captured outbound request (proves substitution occurred) |

## Verification

```bash
# Run the scenario with the specified context and env vars set
# Capture the actual outbound request (via a mock server or request log)
export REHEARSE_CONTEXT_VARS="${context_var_name}=${context_var_value}"

if [ -n "$env_var_name" ]; then
  export "$env_var_name"="$env_var_value"
fi

# Run against a local mock server that records received requests
mock_log=$(mktemp)
# (Mock server setup is provided by the test scenario wrapping this AC)

rehearse run "$scenario_path" 2>&1

# Verify the substituted value appeared in the request received by the mock
grep -q "$expected_in_request" "$mock_log"
rm -f "$mock_log"
```
