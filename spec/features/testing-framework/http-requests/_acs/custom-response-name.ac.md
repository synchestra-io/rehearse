# AC: custom-response-name

**Status:** planned
**Feature:** [testing-framework/http-requests](../README.md)

## Description

When an HTTP step declares an Outputs table, the runner stores response data under the declared output names and suppresses the default `context.response` auto-store. Extract expressions in the Outputs table can reference `$RESPONSE_BODY`, `$RESPONSE_STATUS`, and `$RESPONSE_HEADERS_*`.

## Inputs

| Name | Required | Description |
|---|---|---|
| scenario_path | Yes | Path to a scenario with an HTTP step that has an Outputs table |
| output_name | Yes | Name of the declared output (e.g., `created_user`) |
| expected_value_fragment | Yes | A substring that must appear in the stored output value |

## Verification

```bash
result=$(rehearse run "$scenario_path" --format json 2>&1)

# The declared output must be present in context
output_value=$(echo "$result" | python3 -c "
import json, sys
d = json.load(sys.stdin)
print(d['steps'][0]['context']['$output_name'])
")

echo "$output_value" | grep -q "$expected_value_fragment"

# The default context.response must NOT be present when Outputs table is declared
default_response=$(echo "$result" | python3 -c "
import json, sys
d = json.load(sys.stdin)
ctx = d['steps'][0]['context']
print('present' if 'response' in ctx else 'absent')
")

test "$default_response" = "absent"
```
