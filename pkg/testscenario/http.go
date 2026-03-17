package testscenario

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// httpClient is the HTTP client used for executing http block requests.
// A 30-second timeout prevents hung servers from blocking the test runner indefinitely.
var httpClient = &http.Client{Timeout: 30 * time.Second}

// HTTPResult holds the outcome of an HTTP request.
type HTTPResult struct {
	Status     int
	StatusText string
	Headers    map[string]string // lowercased header names
	Body       string
}

// Env returns HTTP response as RESPONSE_* environment variable pairs.
// RESPONSE_STATUS, RESPONSE_BODY, and RESPONSE_HEADERS_<UPPERCASED_NAME> are set.
// NOTE: values containing newlines produce multi-line env entries, which are silently
// dropped by some OS/shell implementations. This mirrors the existing ContextVarsAsEnv
// behaviour in this package.
func (r HTTPResult) Env() []string {
	env := []string{
		fmt.Sprintf("RESPONSE_STATUS=%d", r.Status),
		"RESPONSE_BODY=" + r.Body,
	}
	for name, val := range r.Headers {
		key := "RESPONSE_HEADERS_" + strings.ToUpper(strings.ReplaceAll(name, "-", "_"))
		env = append(env, key+"="+val)
	}
	return env
}

// ContextVars returns flat key=value pairs for storing in the execution context.
// Keys: "response.status", "response.status_text", "response.body",
//       "response.headers.<lowercased-name>".
func (r HTTPResult) ContextVars() map[string]string {
	vars := map[string]string{
		"response.status":      fmt.Sprintf("%d", r.Status),
		"response.status_text": r.StatusText,
		"response.body":        r.Body,
	}
	for name, val := range r.Headers {
		vars["response.headers."+name] = val
	}
	return vars
}

// parseHTTPBlock parses the content of an http fenced code block.
// Format: METHOD URL\nHeader: value\n\nbody
// Returns method, url, headers (lowercased keys), body, and any error.
func parseHTTPBlock(code string) (method, rawURL string, headers map[string]string, body string, err error) {
	lines := strings.Split(strings.TrimRight(code, "\n"), "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) == "" {
		return "", "", nil, "", fmt.Errorf("http block is empty")
	}

	// First line: METHOD URL
	parts := strings.SplitN(strings.TrimSpace(lines[0]), " ", 2)
	if len(parts) != 2 {
		return "", "", nil, "", fmt.Errorf("http block first line must be 'METHOD URL', got %q", lines[0])
	}
	method = strings.ToUpper(strings.TrimSpace(parts[0]))
	rawURL = strings.TrimSpace(parts[1])

	// Validate absolute URL.
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		return "", "", nil, "", fmt.Errorf("http block URL must be absolute (http:// or https://), got %q", rawURL)
	}

	// Parse headers until blank line.
	headers = make(map[string]string)
	bodyStart := len(lines) // default: no body
	for i := 1; i < len(lines); i++ {
		line := lines[i]
		if strings.TrimSpace(line) == "" {
			bodyStart = i + 1
			break
		}
		name, val, ok := strings.Cut(line, ":")
		if !ok {
			return "", "", nil, "", fmt.Errorf("invalid header line %q", line)
		}
		headers[strings.ToLower(strings.TrimSpace(name))] = strings.TrimSpace(val)
	}

	// Collect body lines.
	if bodyStart < len(lines) {
		body = strings.Join(lines[bodyStart:], "\n")
	}

	return method, rawURL, headers, body, nil
}

// execHTTPRequest executes an HTTP request described by an http block.
// The code must already have ${{ }} variables resolved by the caller.
// Returns HTTPResult on success, or an error for network-level failures.
// Non-2xx HTTP responses are NOT errors — they are valid results.
func execHTTPRequest(code string) (HTTPResult, error) {
	method, rawURL, reqHeaders, body, err := parseHTTPBlock(code)
	if err != nil {
		return HTTPResult{}, fmt.Errorf("parsing http block: %w", err)
	}

	var bodyReader io.Reader
	if body != "" {
		bodyReader = strings.NewReader(body)
	}

	req, err := http.NewRequest(method, rawURL, bodyReader)
	if err != nil {
		return HTTPResult{}, fmt.Errorf("building request: %w", err)
	}
	for name, val := range reqHeaders {
		req.Header.Set(name, val)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return HTTPResult{}, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return HTTPResult{}, fmt.Errorf("reading response body: %w", err)
	}

	// Collect response headers (lowercase keys, first value only).
	respHeaders := make(map[string]string, len(resp.Header))
	for name, vals := range resp.Header {
		if len(vals) > 0 {
			respHeaders[strings.ToLower(name)] = vals[0]
		}
	}

	return HTTPResult{
		Status:     resp.StatusCode,
		StatusText: http.StatusText(resp.StatusCode),
		Headers:    respHeaders,
		Body:       string(respBody),
	}, nil
}
