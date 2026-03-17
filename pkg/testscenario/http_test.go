package testscenario

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestParseHTTPBlock_GetNoBody(t *testing.T) {
	code := "GET https://example.com/\nAccept: text/html\n"
	method, rawURL, headers, body, err := parseHTTPBlock(code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if method != "GET" {
		t.Errorf("expected method GET, got %q", method)
	}
	if rawURL != "https://example.com/" {
		t.Errorf("expected url https://example.com/, got %q", rawURL)
	}
	if headers["accept"] != "text/html" {
		t.Errorf("expected accept header 'text/html', got %q", headers["accept"])
	}
	if body != "" {
		t.Errorf("expected empty body, got %q", body)
	}
}

func TestParseHTTPBlock_PostWithBody(t *testing.T) {
	code := "POST https://api.example.com/users\nContent-Type: application/json\n\n{\"name\":\"alice\"}"
	method, rawURL, headers, body, err := parseHTTPBlock(code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if method != "POST" {
		t.Errorf("expected POST, got %q", method)
	}
	if headers["content-type"] != "application/json" {
		t.Errorf("expected content-type application/json, got %q", headers["content-type"])
	}
	if !strings.Contains(body, "alice") {
		t.Errorf("expected body to contain 'alice', got %q", body)
	}
	_ = rawURL
}

func TestParseHTTPBlock_MissingMethod(t *testing.T) {
	_, _, _, _, err := parseHTTPBlock("")
	if err == nil {
		t.Error("expected error for empty block")
	}
}

func TestParseHTTPBlock_RelativeURL(t *testing.T) {
	_, _, _, _, err := parseHTTPBlock("GET /relative/path\n")
	if err == nil {
		t.Error("expected error for relative URL")
	}
}

func TestExecHTTPRequest_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(200)
		_, _ = w.Write([]byte("hello"))
	}))
	defer srv.Close()

	code := "GET " + srv.URL + "\n"
	result, err := execHTTPRequest(code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Status != 200 {
		t.Errorf("expected status 200, got %d", result.Status)
	}
	if result.Body != "hello" {
		t.Errorf("expected body 'hello', got %q", result.Body)
	}
	if result.Headers["content-type"] != "text/plain" {
		t.Errorf("expected content-type text/plain, got %q", result.Headers["content-type"])
	}
}

func TestExecHTTPRequest_NetworkError(t *testing.T) {
	code := "GET http://127.0.0.1:1\n" // nothing listening there
	_, err := execHTTPRequest(code)
	if err == nil {
		t.Error("expected error for unreachable host")
	}
}

func TestHTTPResult_Env(t *testing.T) {
	r := HTTPResult{
		Status:     200,
		StatusText: "OK",
		Headers:    map[string]string{"content-type": "text/html; charset=utf-8"},
		Body:       "hello",
	}
	env := r.Env()
	envMap := make(map[string]string)
	for _, e := range env {
		k, v, _ := strings.Cut(e, "=")
		envMap[k] = v
	}
	if envMap["RESPONSE_STATUS"] != "200" {
		t.Errorf("expected RESPONSE_STATUS=200, got %q", envMap["RESPONSE_STATUS"])
	}
	if envMap["RESPONSE_BODY"] != "hello" {
		t.Errorf("expected RESPONSE_BODY=hello, got %q", envMap["RESPONSE_BODY"])
	}
	if envMap["RESPONSE_HEADERS_CONTENT_TYPE"] != "text/html; charset=utf-8" {
		t.Errorf("wrong RESPONSE_HEADERS_CONTENT_TYPE: %q", envMap["RESPONSE_HEADERS_CONTENT_TYPE"])
	}
}

func TestHTTPResult_ContextVars(t *testing.T) {
	r := HTTPResult{
		Status:     404,
		StatusText: "Not Found",
		Headers:    map[string]string{"content-type": "text/plain"},
		Body:       "not found",
	}
	vars := r.ContextVars()
	if vars["response.status"] != "404" {
		t.Errorf("expected response.status=404, got %q", vars["response.status"])
	}
	if vars["response.status_text"] != "Not Found" {
		t.Errorf("expected response.status_text='Not Found', got %q", vars["response.status_text"])
	}
	if vars["response.headers.content-type"] != "text/plain" {
		t.Errorf("expected response.headers.content-type=text/plain, got %q", vars["response.headers.content-type"])
	}
	if vars["response.body"] != "not found" {
		t.Errorf("expected response.body='not found', got %q", vars["response.body"])
	}
}
