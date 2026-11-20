package requestid

import (
	"context"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

var idFormat = regexp.MustCompile(`^[0-9a-f]{32}$`)

// serve runs one request through the middleware and returns the id the
// backend saw and the id in the response.
func serve(t *testing.T, cfg *Config, incoming string) (backend, response string) {
	t.Helper()
	next := http.HandlerFunc(func(_ http.ResponseWriter, req *http.Request) {
		backend = req.Header.Get(cfg.HeaderName)
	})
	h, err := New(context.Background(), next, cfg, "requestid")
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
	if incoming != "" {
		req.Header.Set(cfg.HeaderName, incoming)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return backend, rec.Header().Get(cfg.HeaderName)
}

func TestSetsSameIDOnRequestAndResponse(t *testing.T) {
	backend, response := serve(t, CreateConfig(), "")
	if !idFormat.MatchString(backend) {
		t.Fatalf("id %q is not 32 hex characters", backend)
	}
	if backend != response {
		t.Fatalf("backend got %q, response has %q", backend, response)
	}
}

func TestNewIDPerRequest(t *testing.T) {
	first, _ := serve(t, CreateConfig(), "")
	second, _ := serve(t, CreateConfig(), "")
	if first == second {
		t.Fatalf("two requests got the same id %q", first)
	}
}

func TestReplacesIncomingID(t *testing.T) {
	backend, _ := serve(t, CreateConfig(), "spoofed")
	if backend == "spoofed" {
		t.Fatal("incoming id was kept without trustIncoming")
	}
}

func TestTrustIncomingKeepsID(t *testing.T) {
	cfg := CreateConfig()
	cfg.TrustIncoming = true
	backend, response := serve(t, cfg, "from-upstream")
	if backend != "from-upstream" || response != "from-upstream" {
		t.Fatalf("got backend %q, response %q, want from-upstream", backend, response)
	}
}

func TestCustomHeaderName(t *testing.T) {
	cfg := &Config{HeaderName: "X-Trace"}
	backend, response := serve(t, cfg, "")
	if !idFormat.MatchString(backend) || backend != response {
		t.Fatalf("got backend %q, response %q", backend, response)
	}
}
