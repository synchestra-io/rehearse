# AC: stores-response-to-context

**Status:** planned
**Feature:** [testing-framework/http-requests](../README.md)

## Description

After a successful HTTP request, the runner stores the full response as `context.response` containing `status` (integer), `status_text` (string), `headers` (object with lowercased keys), `body` (raw string), and `json` (parsed object when response Content-Type is application/json). When no Outputs table is declared, the default `context.response` store is used automatically.

## Inputs

| Name | Required | Description |
|---|---|---|
| scenario_path | Yes | Path to a scenario making an HTTP request with no Outputs table |
| expected_status | Yes | Expected HTTP status code as a string (e.g., `"200"`) |
| expected_body_fragment | No | A substring that must appear in `context.response.body` |
| expected_json_key | No | A key that must exist in `context.response.json` |

## Verification

```bash
# Run the scenario and capture the context state after the HTTP step
result=$(rehearse run "$scenario_path" --format json 2>&1)

status=$(echo "$result" | python3 -c "
import json, sys
d = json.load(sys.stdin)
print(d['steps'][0]['context']['response']['status'])
")

test "$status" = "$expected_status"

if [ -n "$expected_body_fragment" ]; then
  body=$(echo "$result" | python3 -c "
import json, sys
d = json.load(sys.stdin)
print(d['steps'][0]['context']['response']['body'])
")
  echo "$body" | grep -q "$expected_body_fragment"
fi

if [ -n "$expected_json_key" ]; then
  python3 -c "
import json, sys
result = json.loads(open('/dev/stdin').read())
response_json = result['steps'][0]['context']['response']['json']
assert '$expected_json_key' in response_json, f'Key $expected_json_key not found in response.json'
" <<< "$result"
fi
```
