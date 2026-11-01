package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestDashboardDeepLinks(t *testing.T) {
	directory := t.TempDir()
	if err := os.WriteFile(filepath.Join(directory, "index.html"), []byte("Saum dashboard"), 0600); err != nil {
		t.Fatal(err)
	}
	handler := testStore(t).handler(dashboardHandler(directory))
	for _, path := range []string{"/", "/factories/saum", "/factories/saum/sessions/session-1", "/stories/story-1", "/stories/story-1/sessions/session-1"} {
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, httptest.NewRequest("GET", "http://127.0.0.1:8080"+path, nil))
		if response.Code != http.StatusOK || response.Body.String() != "Saum dashboard" {
			t.Fatalf("deep link %s failed: %d %s", path, response.Code, response.Body.String())
		}
	}
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest("GET", "http://127.0.0.1:8080/api/missing", nil))
	if response.Code != 404 {
		t.Fatal("API miss rendered dashboard")
	}
}
