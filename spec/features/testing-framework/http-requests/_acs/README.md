# Acceptance Criteria: HTTP Requests

| AC | Description | Status |
|---|---|---|
| [parses-http-block-format](parses-http-block-format.ac.md) | runner correctly parses method, url, headers, and body from http code block | planned |
| [substitutes-context-vars](substitutes-context-vars.ac.md) | context and env vars are resolved in method line, headers, and body before the request is sent | planned |
| [stores-response-to-context](stores-response-to-context.ac.md) | response object (status, headers, body, json) is stored in context.response by default | planned |
| [custom-response-name](custom-response-name.ac.md) | when Outputs table is present, response is stored under declared names and default context.response is suppressed | planned |
| [network-error-fails-step](network-error-fails-step.ac.md) | network errors fail the step; non-2xx HTTP responses do not | planned |
