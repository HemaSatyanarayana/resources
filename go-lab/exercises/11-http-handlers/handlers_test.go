package httpsvc

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()
	HealthHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rr.Code)
	}
	if body := strings.TrimSpace(rr.Body.String()); body != "ok" {
		t.Errorf("body = %q, want %q", body, "ok")
	}
}

func TestGreetHandler(t *testing.T) {
	cases := map[string]string{
		"/greet?name=Ada": "Hello, Ada!",
		"/greet":          "Hello, World!",
		"/greet?name=":    "Hello, World!",
	}
	for url, want := range cases {
		req := httptest.NewRequest(http.MethodGet, url, nil)
		rr := httptest.NewRecorder()
		GreetHandler(rr, req)
		if got := strings.TrimSpace(rr.Body.String()); got != want {
			t.Errorf("GET %s body = %q, want %q", url, got, want)
		}
	}
}

func TestItemHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/item", nil)
	rr := httptest.NewRecorder()
	ItemHandler(rr, req)

	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}
	var got Item
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("body was not valid JSON: %v (%s)", err, rr.Body.String())
	}
	want := Item{ID: 7, Name: "widget", Price: 500}
	if got != want {
		t.Errorf("item = %+v, want %+v", got, want)
	}
}

func TestRouter(t *testing.T) {
	srv := httptest.NewServer(NewRouter())
	defer srv.Close()

	get := func(path string) (int, string) {
		resp, err := http.Get(srv.URL + path)
		if err != nil {
			t.Fatalf("GET %s: %v", path, err)
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return resp.StatusCode, strings.TrimSpace(string(body))
	}

	if code, body := get("/health"); code != 200 || body != "ok" {
		t.Errorf("/health = (%d, %q)", code, body)
	}
	if code, body := get("/greet?name=Go"); code != 200 || body != "Hello, Go!" {
		t.Errorf("/greet = (%d, %q)", code, body)
	}
	if code, body := get("/items/42"); code != 200 || body != "item 42" {
		t.Errorf("/items/42 = (%d, %q)", code, body)
	}
	if code, _ := get("/nope"); code != 404 {
		t.Errorf("/nope status = %d, want 404", code)
	}
}
