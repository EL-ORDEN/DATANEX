package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientDoGET(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET request, got %s", r.Method)
		}
		if r.Header.Get("Authorization") != "Bearer demo-token" {
			t.Fatalf("missing bearer token")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	client := NewClient()
	resp, err := client.Do(Request{
		Method:      http.MethodGet,
		URL:         server.URL,
		Headers:     map[string]string{"X-Test": "true"},
		BearerToken: "demo-token",
	})
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if len(resp.Body) == 0 {
		t.Fatal("expected response body")
	}
}

func TestClientDoPOST(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST request, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Fatalf("missing JSON content type")
		}
		w.WriteHeader(http.StatusCreated)
		_, _ = fmt.Fprint(w, `{"created":true}`)
	}))
	defer server.Close()

	client := NewClient()
	resp, err := client.Do(Request{
		Method: http.MethodPost,
		URL:    server.URL,
		Body: map[string]any{
			"name": "Ada",
		},
	})
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}
}
