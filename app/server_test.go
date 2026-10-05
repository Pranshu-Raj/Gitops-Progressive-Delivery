package main

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"
	"time"
)

func testServer(t *testing.T, cfg config) (*server, *httptest.Server) {
	t.Helper()
	cfg.env = "test"
	s := newServer(cfg, slog.New(slog.NewTextHandler(io.Discard, nil)))
	ts := httptest.NewServer(s)
	t.Cleanup(ts.Close)
	return s, ts
}

func get(t *testing.T, url string) (*http.Response, string) {
	t.Helper()
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return resp, string(body)
}

func TestRootReportsVersion(t *testing.T) {
	version = "abc123"
	_, ts := testServer(t, config{})

	resp, body := get(t, ts.URL+"/")
	if resp.StatusCode != 200 {
		t.Fatalf("status %d", resp.StatusCode)
	}
	var got map[string]string
	if err := json.Unmarshal([]byte(body), &got); err != nil {
		t.Fatal(err)
	}
	if got["version"] != "abc123" || got["env"] != "test" {
		t.Errorf("got %v", got)
	}
}

func TestFailEvery(t *testing.T) {
	_, ts := testServer(t, config{failEvery: 3})

	var codes []int
	for range 6 {
		resp, _ := get(t, ts.URL+"/")
		codes = append(codes, resp.StatusCode)
	}
	want := []int{200, 200, 500, 200, 200, 500}
	if !slices.Equal(codes, want) {
		t.Errorf("got %v want %v", codes, want)
	}
}

func TestSlowMs(t *testing.T) {
	_, ts := testServer(t, config{slowMs: 50})

	start := time.Now()
	resp, _ := get(t, ts.URL+"/")
	if resp.StatusCode != 200 {
		t.Fatalf("status %d", resp.StatusCode)
	}
	if d := time.Since(start); d < 50*time.Millisecond {
		t.Errorf("took %v, expected at least 50ms", d)
	}
}

func TestReadyzWhileDraining(t *testing.T) {
	s, ts := testServer(t, config{})

	if resp, _ := get(t, ts.URL+"/readyz"); resp.StatusCode != 200 {
		t.Fatalf("before drain: %d", resp.StatusCode)
	}
	s.drain()
	if resp, _ := get(t, ts.URL+"/readyz"); resp.StatusCode != 503 {
		t.Fatalf("after drain: %d", resp.StatusCode)
	}
	if resp, _ := get(t, ts.URL+"/healthz"); resp.StatusCode != 200 {
		t.Fatalf("healthz during drain: %d", resp.StatusCode)
	}
}

func TestRequestID(t *testing.T) {
	_, ts := testServer(t, config{})

	req, _ := http.NewRequest("GET", ts.URL+"/", nil)
	req.Header.Set("X-Request-Id", "req-42")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if got := resp.Header.Get("X-Request-Id"); got != "req-42" {
		t.Errorf("echoed %q", got)
	}

	resp, _ = get(t, ts.URL+"/")
	if got := resp.Header.Get("X-Request-Id"); len(got) != 16 {
		t.Errorf("generated %q", got)
	}
}

func TestUnknownPath(t *testing.T) {
	_, ts := testServer(t, config{})
	if resp, _ := get(t, ts.URL+"/nope"); resp.StatusCode != 404 {
		t.Errorf("status %d", resp.StatusCode)
	}
}

func TestMetrics(t *testing.T) {
	version = "abc123"
	_, ts := testServer(t, config{failEvery: 2})

	get(t, ts.URL+"/")
	get(t, ts.URL+"/")
	_, body := get(t, ts.URL+"/metrics")

	for _, want := range []string{
		`http_requests_total{path="/",code="200",version="abc123"} 1`,
		`http_requests_total{path="/",code="500",version="abc123"} 1`,
		`http_request_duration_seconds_count{path="/",version="abc123"} 2`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("missing %q in:\n%s", want, body)
		}
	}
}
