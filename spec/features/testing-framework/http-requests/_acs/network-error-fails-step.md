# AC: network-error-fails-step

**Status:** planned
**Feature:** [testing-framework/http-requests](../README.md)

## Description

A network-level error (unreachable host, DNS failure, connection timeout) causes the HTTP step to fail with exit code 1 and an error message. The step's context is not updated. A non-2xx HTTP status code (4xx, 5xx) does NOT fail the step — it is a valid completed request with a response that assertions may verify.

## Inputs

| Name | Required | Description |
|---|---|---|
| unreachable_scenario_path | Yes | Path to a scenario making a request to an unreachable host |
| non2xx_scenario_path | Yes | Path to a scenario making a request that returns a 4xx or 5xx response |
| expected_error_fragment | No | Substring that must appear in the error output for the network failure case |

## Verification

```bash
# Case 1: network error → step fails, exit code non-zero
rehearse run "$unreachable_scenario_path" --format json > /tmp/network_result.json 2>&1
exit_code=$?

test $exit_code -ne 0

step_status=$(python3 -c "
import json
d = json.load(open('/tmp/network_result.json'))
print(d['steps'][0]['status'])
")
test "$step_status" = "failed"

if [ -n "$expected_error_fragment" ]; then
  grep -q "$expected_error_fragment" /tmp/network_result.json
fi

# Case 2: non-2xx response → step succeeds (status "passed"), response stored in context
rehearse run "$non2xx_scenario_path" --format json > /tmp/non2xx_result.json 2>&1
non2xx_exit=$?

step_status_non2xx=$(python3 -c "
import json
d = json.load(open('/tmp/non2xx_result.json'))
print(d['steps'][0]['status'])
")
test "$step_status_non2xx" = "passed"

# Response should still be stored even for error status
response_status=$(python3 -c "
import json
d = json.load(open('/tmp/non2xx_result.json'))
print(d['steps'][0]['context']['response']['status'])
")
# Must be a 4xx or 5xx (>= 400)
python3 -c "assert int('$response_status') >= 400"

rm -f /tmp/network_result.json /tmp/non2xx_result.json
```
