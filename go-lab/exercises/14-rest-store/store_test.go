package store

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestStoreCRUD(t *testing.T) {
	s := NewStore()

	a := s.Create("write tests")
	b := s.Create("ship it")
	if a.ID != 1 || b.ID != 2 {
		t.Fatalf("IDs = %d, %d; want 1, 2", a.ID, b.ID)
	}

	got, ok := s.Get(1)
	if !ok || got.Title != "write tests" {
		t.Errorf("Get(1) = %+v, %v", got, ok)
	}
	if _, ok := s.Get(999); ok {
		t.Error("Get(999) should be false")
	}

	list := s.List()
	if len(list) != 2 || list[0].ID != 1 || list[1].ID != 2 {
		t.Errorf("List = %+v, want sorted [1,2]", list)
	}

	if !s.Delete(1) {
		t.Error("Delete(1) should be true")
	}
	if s.Delete(1) {
		t.Error("second Delete(1) should be false")
	}
	if len(s.List()) != 1 {
		t.Errorf("after delete, len = %d, want 1", len(s.List()))
	}
}

func TestStoreConcurrentCreate(t *testing.T) {
	s := NewStore()
	const n = 500
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			s.Create("t")
		}()
	}
	wg.Wait()

	list := s.List()
	if len(list) != n {
		t.Fatalf("created %d tasks, want %d", len(list), n)
	}
	seen := make(map[int]bool)
	for _, task := range list {
		if seen[task.ID] {
			t.Fatalf("duplicate ID %d — the store is not concurrency-safe", task.ID)
		}
		seen[task.ID] = true
	}
}

func TestRESTAPI(t *testing.T) {
	srv := httptest.NewServer(NewRouter(NewStore()))
	defer srv.Close()

	// POST /tasks
	resp, err := http.Post(srv.URL+"/tasks", "application/json",
		bytes.NewBufferString(`{"title":"learn go"}`))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("POST status = %d, want 201", resp.StatusCode)
	}
	var created Task
	json.NewDecoder(resp.Body).Decode(&created)
	resp.Body.Close()
	if created.ID != 1 || created.Title != "learn go" {
		t.Errorf("created = %+v", created)
	}

	// GET /tasks/1
	resp, _ = http.Get(fmt.Sprintf("%s/tasks/%d", srv.URL, created.ID))
	if resp.StatusCode != http.StatusOK {
		t.Errorf("GET one status = %d, want 200", resp.StatusCode)
	}
	resp.Body.Close()

	// GET /tasks (list)
	resp, _ = http.Get(srv.URL + "/tasks")
	var list []Task
	json.NewDecoder(resp.Body).Decode(&list)
	resp.Body.Close()
	if len(list) != 1 {
		t.Errorf("list len = %d, want 1", len(list))
	}

	// GET missing -> 404
	resp, _ = http.Get(srv.URL + "/tasks/999")
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("GET missing status = %d, want 404", resp.StatusCode)
	}
	resp.Body.Close()

	// DELETE /tasks/1 -> 204
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/tasks/%d", srv.URL, created.ID), nil)
	resp, _ = http.DefaultClient.Do(req)
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("DELETE status = %d, want 204", resp.StatusCode)
	}
	resp.Body.Close()

	// Bad ID -> 400
	resp, _ = http.Get(srv.URL + "/tasks/abc")
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("bad id status = %d, want 400", resp.StatusCode)
	}
	resp.Body.Close()
}
