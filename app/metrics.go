package main

import (
	"bytes"
	"fmt"
	"io"
	"maps"
	"slices"
	"strconv"
	"sync"
)

var buckets = []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5}

type metrics struct {
	mu       sync.Mutex
	requests map[string]map[int]uint64
	latency  map[string]*histogram
}

type histogram struct {
	counts []uint64
	sum    float64
	total  uint64
}

func newMetrics() *metrics {
	return &metrics{
		requests: make(map[string]map[int]uint64),
		latency:  make(map[string]*histogram),
	}
}

func (m *metrics) observe(path string, code int, seconds float64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.requests[path] == nil {
		m.requests[path] = make(map[int]uint64)
	}
	m.requests[path][code]++

	h := m.latency[path]
	if h == nil {
		h = &histogram{counts: make([]uint64, len(buckets))}
		m.latency[path] = h
	}
	for i, le := range buckets {
		if seconds <= le {
			h.counts[i]++
		}
	}
	h.sum += seconds
	h.total++
}

func (m *metrics) write(w io.Writer) {
	var buf bytes.Buffer
	m.mu.Lock()

	fmt.Fprintln(&buf, "# TYPE http_requests_total counter")
	for _, path := range slices.Sorted(maps.Keys(m.requests)) {
		for _, code := range slices.Sorted(maps.Keys(m.requests[path])) {
			fmt.Fprintf(&buf, "http_requests_total{path=%q,code=\"%d\",version=%q} %d\n",
				path, code, version, m.requests[path][code])
		}
	}

	fmt.Fprintln(&buf, "# TYPE http_request_duration_seconds histogram")
	for _, path := range slices.Sorted(maps.Keys(m.latency)) {
		h := m.latency[path]
		for i, le := range buckets {
			fmt.Fprintf(&buf, "http_request_duration_seconds_bucket{path=%q,version=%q,le=%q} %d\n",
				path, version, strconv.FormatFloat(le, 'g', -1, 64), h.counts[i])
		}
		fmt.Fprintf(&buf, "http_request_duration_seconds_bucket{path=%q,version=%q,le=\"+Inf\"} %d\n", path, version, h.total)
		fmt.Fprintf(&buf, "http_request_duration_seconds_sum{path=%q,version=%q} %g\n", path, version, h.sum)
		fmt.Fprintf(&buf, "http_request_duration_seconds_count{path=%q,version=%q} %d\n", path, version, h.total)
	}

	m.mu.Unlock()
	buf.WriteTo(w)
}
