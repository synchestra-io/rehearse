# Feature: HTTP Requests

**Status:** Conceptual

## Summary

HTTP request blocks let scenarios make real HTTP(S) calls using a standard, human-readable request format. Declare the method, URL, headers, and body directly in a fenced code block. The runner executes the request, and the response is automatically stored in context for subsequent steps and assertions to consume. Variables from context and environment are substituted before the request is sent — no templating language to learn, same `${{ }}` syntax used everywhere else in Rehearse.

## Problem

Testing HTTP APIs today means choosing between the shell script approach and framework-specific DSLs:

- **curl in bash blocks** works, but produces shell noise that obscures what request is actually being made. Parsing the response requires piping through `jq`, capturing to variables, and managing exit codes manually. The intent — "call this endpoint, check that" — is buried in mechanics.
- **Framework HTTP clients** (Supertest, RestAssured, etc.) are language-specific and require writing code in whatever language the framework uses, not the language best suited to the assertion.
- **Postman / Bruno / Insomnia collections** are JSON/YAML configs that look nothing like HTTP. They are not readable as documentation and do not integrate with scenario AC verification.

HTTP request blocks solve this by using a format that reads like the actual HTTP spec — because it is. The `.http` file format (used by JetBrains HTTP Client, VS Code REST Client, and codified in RFC 7230) is familiar, language-neutral, and renders readably on GitHub. The runner executes it; the response lands on context; assertions verify it.

## Behavior

### HTTP request block format

HTTP requests are defined in fenced code blocks with the `http` annotation. The format follows RFC 7230 HTTP/1.1 message syntax:

```markdown
```http
METHOD https://host/path
Header-Name: value
Another-Header: value

Body content (optional)
```
```

**Components:**

| Component | Required | Description |
|---|---|---|
| Method line | Yes | `METHOD URL` — first line of the block. Method is uppercase (`GET`, `POST`, `PUT`, `PATCH`, `DELETE`, `HEAD`, `OPTIONS`). URL must be absolute (`https://` or `http://`). |
| Headers | No | One per line, immediately after the method line. `Header-Name: value` format. No blank line between method and headers. |
| Blank line | No | Separates headers from body. Required when a body is present. |
| Body | No | Raw body content. Can be JSON, form data, XML, or any string. Follows the blank line separator. |

**Full example:**

```http
POST https://api.example.com/users
Authorization: Bearer ${{ context.auth_token }}
Content-Type: application/json
X-Request-ID: ${{ env.REQUEST_ID }}

{
  "name": "${{ context.user_name }}",
  "email": "${{ context.user_email }}"
}
```

**GET with no body:**

```http
GET https://api.example.com/users/${{ context.user_id }}
Authorization: Bearer ${{ context.auth_token }}
Accept: application/json
```

### Variable substitution

Before the request is sent, the runner substitutes all `${{ }}` references in the method line, headers, and body:

| Syntax | Resolves from | Example |
|---|---|---|
| `${{ context.name }}` | Scenario context (set by previous step outputs) | `${{ context.auth_token }}` |
| `${{ env.NAME }}` | Process environment variables | `${{ env.API_BASE_URL }}` |

Substitution follows the same rules as [scenario variable resolution](../test-scenario/README.md#output-model): unresolved references are a runtime error that fails the step before the request is sent.

### Response context

After a successful HTTP call (any completed request, regardless of status code), the runner stores the response as a structured object in context. By default, the response is stored under the key `response`:

| Context key | Type | Description |
|---|---|---|
| `context.response.status` | Integer | HTTP status code (e.g., `200`, `404`) |
| `context.response.status_text` | String | HTTP reason phrase (e.g., `OK`, `Not Found`) |
| `context.response.headers` | Object | Response headers as key-value pairs (header names lowercased) |
| `context.response.body` | String | Raw response body |
| `context.response.json` | Object | Parsed JSON body — set only when `Content-Type: application/json` is in the response |

**Environment variables set for assertions and subsequent steps:**

| Variable | Description |
|---|---|
| `RESPONSE_STATUS` | HTTP status code as a string (e.g., `"201"`) |
| `RESPONSE_BODY` | Raw response body |
| `RESPONSE_HEADERS_*` | Each response header as `RESPONSE_HEADERS_{UPPERCASED_NAME}` |

These environment variables follow the same convention as `STEP_STDOUT` and `STEP_STDERR` — they are available within the step's [assertions](#) and for use in Extract expressions.

### Custom response name

By default the response is stored as `context.response`. When a step makes multiple HTTP calls or needs a more meaningful name, the store key is configured via the **Outputs** table:

```markdown
## create-user

**Outputs:**

| Name | Store | Extract |
|---|---|---|
| created_user | context | `echo "$RESPONSE_BODY"` |
| user_status | context | `echo "$RESPONSE_STATUS"` |

```http
POST https://api.example.com/users
Content-Type: application/json

{"name": "Alice"}
```
```

When an Outputs table is present, the default `context.response` auto-store is suppressed. Only the declared outputs are stored. The Extract expression runs against the response environment variables (`$RESPONSE_BODY`, `$RESPONSE_STATUS`, `$RESPONSE_HEADERS_*`).

If no Outputs table is declared, the full response object is stored as `context.response` (overwriting any previous value).

### Error handling

| Condition | Behavior |
|---|---|
| Network error (unreachable, DNS failure, timeout) | Step fails with exit code 1. Response is not stored. |
| Unresolved variable reference | Step fails before the request is sent. Error message names the unresolved reference. |
| Non-2xx HTTP status | Step succeeds — an HTTP error status is a valid response, not a runner error. Assertions verify expected status codes. |
| Malformed request block (no method, relative URL) | Validation error at parse time — scenario fails to load. |

A non-2xx status does not fail the step because HTTP semantics are domain-specific: a `404` is expected when testing "user not found" flows, and a `422` is expected when testing validation rejection. Assertions verify the status, not the runner.

### Multiple HTTP requests in one scenario

Each HTTP step makes exactly one request. For workflows requiring multiple calls, use multiple steps and pass response data via context:

```markdown
## authenticate

```http
POST https://api.example.com/auth/token
Content-Type: application/json

{"username": "${{ env.TEST_USER }}", "password": "${{ env.TEST_PASSWORD }}"}
```

## create-resource

**Depends on:** authenticate

```http
POST https://api.example.com/resources
Authorization: Bearer ${{ context.response.json.access_token }}
Content-Type: application/json

{"name": "my-resource"}
```
```

## Interaction with Other Features

| Feature | Interaction |
|---|---|
| [Test Scenario](../test-scenario/README.md) | HTTP blocks are a step code block type alongside `bash`, `python`, `sql`, and `starlark`. Steps with `http` blocks support the same Outputs table, ACs, and Assertions as other step types. |
| [Test Runner](../test-runner/README.md) | The runner detects `http` code block annotation, parses the request format, substitutes variables, executes the HTTP call, and stores the response in context. |
| [Acceptance Criteria](../../acceptance-criteria/README.md) | AC verification scripts can reference `$RESPONSE_BODY`, `$RESPONSE_STATUS`, and `$RESPONSE_HEADERS_*` when the step they verify is an HTTP request step. |

## Acceptance Criteria

| AC | Description | Status |
|---|---|---|
| [substitutes-context-vars](_acs/substitutes-context-vars.ac.md) | context and env vars are resolved in method line, headers, and body before the request is sent | planned |
| [stores-response-to-context](_acs/stores-response-to-context.ac.md) | response object (status, headers, body, json) is stored in context.response by default | planned |
| [custom-response-name](_acs/custom-response-name.ac.md) | when Outputs table is present, response is stored under declared names and default context.response is suppressed | planned |
| [network-error-fails-step](_acs/network-error-fails-step.ac.md) | network errors (unreachable host, timeout) fail the step; non-2xx responses do not | planned |
| [parses-http-block-format](_acs/parses-http-block-format.ac.md) | runner correctly parses method, url, headers, and body from http code block | planned |

## Outstanding Questions

- Should `http` blocks support request timeouts via a `**Timeout:**` metadata field, or inherit from a global runner setting?
- Should response headers be accessible as `context.response.headers.content-type` (dot notation) in addition to the `RESPONSE_HEADERS_CONTENT_TYPE` env var convention?
- Should the runner follow HTTP redirects by default, or require an explicit `**Follow-Redirects:** true` metadata field?
- Should there be a way to define base URLs or common headers at the scenario level to avoid repeating them per step?
