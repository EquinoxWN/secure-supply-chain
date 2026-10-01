package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func get(t *testing.T, h http.Handler, method, path string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(method, path, nil))
	return rec
}

func TestHealthz(t *testing.T) {
	rec := get(t, New(BuildInfo{}), http.MethodGet, "/healthz")
	if rec.Code != http.StatusOK || rec.Body.String() != "{\"status\":\"ok\"}\n" {
		t.Fatalf("got %d %q", rec.Code, rec.Body.String())
	}
}

func TestVersionReportsBuild(t *testing.T) {
	rec := get(t, New(BuildInfo{Version: "1.2.3", Commit: "abc123"}), http.MethodGet, "/version")
	var got BuildInfo
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.Version != "1.2.3" || got.Commit != "abc123" || got.Go == "" {
		t.Fatalf("got %+v", got)
	}
	if rec.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("content type %q", rec.Header().Get("Content-Type"))
	}
}

func TestUnknownRoutesAndMethods(t *testing.T) {
	h := New(BuildInfo{})
	if rec := get(t, h, http.MethodGet, "/admin"); rec.Code != http.StatusNotFound {
		t.Errorf("unknown path: %d", rec.Code)
	}
	if rec := get(t, h, http.MethodPost, "/healthz"); rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("POST /healthz: %d", rec.Code)
	}
}
