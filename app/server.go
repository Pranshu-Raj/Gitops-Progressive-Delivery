package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"sync/atomic"
	"time"
)

type server struct {
	cfg      config
	log      *slog.Logger
	mux      *http.ServeMux
	metrics  *metrics
	hostname string
	hits     atomic.Int64
	draining atomic.Bool
}

func newServer(cfg config, log *slog.Logger) *server {
	hostname, _ := os.Hostname()
	s := &server{
		cfg:      cfg,
		log:      log,
		mux:      http.NewServeMux(),
		metrics:  newMetrics(),
		hostname: hostname,
	}
	s.mux.HandleFunc("GET /{$}", s.handleRoot)
	s.mux.HandleFunc("GET /healthz", s.handleHealthz)
	s.mux.HandleFunc("GET /readyz", s.handleReadyz)
	s.mux.HandleFunc("GET /metrics", s.handleMetrics)
	return s
}

func (s *server) drain() {
	s.draining.Store(true)
}

var quietPaths = map[string]bool{"/healthz": true, "/readyz": true, "/metrics": true}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	rid := r.Header.Get("X-Request-Id")
	if rid == "" {
		rid = newRequestID()
	}
	w.Header().Set("X-Request-Id", rid)

	rw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
	s.mux.ServeHTTP(rw, r)

	elapsed := time.Since(start)
	s.metrics.observe(routeLabel(r), rw.status, elapsed.Seconds())
	if !quietPaths[r.URL.Path] {
		s.log.Info("request", "rid", rid, "method", r.Method, "path", r.URL.Path, "status", rw.status, "ms", elapsed.Milliseconds())
	}
}

func (s *server) handleRoot(w http.ResponseWriter, r *http.Request) {
	n := s.hits.Add(1)

	if s.cfg.slowMs > 0 {
		select {
		case <-time.After(time.Duration(s.cfg.slowMs) * time.Millisecond):
		case <-r.Context().Done():
			w.WriteHeader(499)
			return
		}
	}
	if s.cfg.failEvery > 0 && n%int64(s.cfg.failEvery) == 0 {
		http.Error(w, "injected failure", http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{
		"version":  version,
		"env":      s.cfg.env,
		"hostname": s.hostname,
	})
}

func (s *server) handleHealthz(w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w, "ok\n")
}

func (s *server) handleReadyz(w http.ResponseWriter, _ *http.Request) {
	if s.draining.Load() {
		http.Error(w, "draining", http.StatusServiceUnavailable)
		return
	}
	io.WriteString(w, "ok\n")
}

func (s *server) handleMetrics(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	s.metrics.write(w)
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func routeLabel(r *http.Request) string {
	switch r.URL.Path {
	case "/", "/healthz", "/readyz", "/metrics":
		return r.URL.Path
	}
	return "other"
}

func newRequestID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}
