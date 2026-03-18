# AC: parses-http-block-format

**Status:** planned
**Feature:** [testing-framework/http-requests](../README.md)

## Description

The runner correctly parses an `http` fenced code block into its constituent parts: HTTP method, absolute URL, headers, and body. Malformed blocks (no method, relative URL, invalid header syntax) are detected at parse time and prevent the scenario from loading.

## Inputs

| Name | Required | Description |
|---|---|---|
| scenario_path | Yes | Path to the scenario file containing an `http` code block |
| expected_method | Yes | Expected HTTP method (e.g., `POST`) |
| expected_url | Yes | Expected request URL |
| expected_header_name | No | A header name that must appear in the parsed request |
| expected_header_value | No | Expected value for the named header |

## Verification

```bash
# Run the scenario through the parser (dry-run / parse-only mode)
# and verify parsed request fields are reported correctly
parsed_output=$(rehearse parse "$scenario_path" --format json 2>&1)

method=$(echo "$parsed_output" | python3 -c "import json,sys; d=json.load(sys.stdin); print(d['steps'][0]['http_request']['method'])")
url=$(echo "$parsed_output" | python3 -c "import json,sys; d=json.load(sys.stdin); print(d['steps'][0]['http_request']['url'])")

test "$method" = "$expected_method"
test "$url" = "$expected_url"

if [ -n "$expected_header_name" ]; then
  header_val=$(echo "$parsed_output" | python3 -c "
import json,sys
d = json.load(sys.stdin)
headers = d['steps'][0]['http_request']['headers']
print(headers.get('$expected_header_name', ''))
")
  test "$header_val" = "$expected_header_value"
fi
```
