package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func okHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
}

func do(h http.Handler, setup func(*http.Request)) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if setup != nil {
		setup(req)
	}
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	return rr
}

func TestRequireAuth(t *testing.T) {
	h := RequireAuth(okHandler())

	// Missing/incorrect token -> 401.
	rr := do(h, nil)
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("no auth: status = %d, want 401", rr.Code)
	}
	if got := strings.TrimSpace(rr.Body.String()); got != "unauthorized" {
		t.Errorf("no auth: body = %q", got)
	}

	// Correct token -> passes through.
	rr = do(h, func(r *http.Request) { r.Header.Set("Authorization", "Bearer secret") })
	if rr.Code != http.StatusOK || strings.TrimSpace(rr.Body.String()) != "ok" {
		t.Errorf("with auth: (%d, %q)", rr.Code, rr.Body.String())
	}
}

func TestRecover(t *testing.T) {
	boom := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("kaboom")
	})
	h := Recover(boom)

	rr := do(h, nil) // must NOT panic the test
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", rr.Code)
	}
	if got := strings.TrimSpace(rr.Body.String()); got != "internal error" {
		t.Errorf("body = %q, want %q", got, "internal error")
	}
}

func TestChainOrderAndCount(t *testing.T) {
	hits := 0
	// Recover (outer) -> Count -> RequireAuth (inner) -> handler.
	h := Chain(okHandler(), Recover, Count(&hits), RequireAuth)

	// Unauthorized: chain still ran Count before RequireAuth rejected.
	rr := do(h, nil)
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want 401", rr.Code)
	}
	if hits != 1 {
		t.Errorf("hits = %d, want 1 (Count runs before RequireAuth)", hits)
	}

	// Authorized: reaches the handler.
	rr = do(h, func(r *http.Request) { r.Header.Set("Authorization", "Bearer secret") })
	if rr.Code != http.StatusOK || strings.TrimSpace(rr.Body.String()) != "ok" {
		t.Errorf("authorized: (%d, %q)", rr.Code, rr.Body.String())
	}
	if hits != 2 {
		t.Errorf("hits = %d, want 2", hits)
	}
}
