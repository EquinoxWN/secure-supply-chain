// Package server is the small HTTP service that the pipeline builds, scans and ships.
package server

import (
	"encoding/json"
	"net/http"
	"runtime"
)

// BuildInfo identifies exactly which build is running, so a deployed image can be traced to a commit.
type BuildInfo struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Go      string `json:"go"`
}

// New returns the service's routes.
func New(info BuildInfo) http.Handler {
	if info.Go == "" {
		info.Go = runtime.Version()
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("GET /version", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, info)
	})
	return mux
}

// writeJSON sends v as a JSON response.
func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	_ = json.NewEncoder(w).Encode(v)
}
