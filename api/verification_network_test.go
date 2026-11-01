package main

import (
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"sync/atomic"
	"testing"
	"time"
)

func TestVerificationAllowsLocalhostStartupAndRequests(t *testing.T) {
	s := testStore(t)
	plan, _ := submittedPlan(t, s)
	if err := s.approvePlan(plan.ID); err != nil {
		t.Fatal(err)
	}
	goal := startClaimed(t, s)
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		_, _ = io.WriteString(w, "host reachable")
	}))
	defer server.Close()
	binary, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("SAUM_SESSION_TOKEN", "must-not-inherit")
	check := successfulCheck(t, s, goal, []string{"/usr/bin/env", "SAUM_TEST_NETWORK_HELPER=1", "SAUM_TEST_NETWORK_URL=" + server.URL, binary, "-test.run=^TestVerificationNetworkHelper$"})
	if requests.Load() != 1 {
		t.Fatalf("host request missing: %+v", check)
	}
}
func TestVerificationNetworkHelper(t *testing.T) {
	if os.Getenv("SAUM_TEST_NETWORK_HELPER") != "1" {
		return
	}
	if os.Getenv("SAUM_SESSION_TOKEN") != "" {
		t.Fatal("inherited session token")
	}
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	server := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { _, _ = io.WriteString(w, "started") }), ReadHeaderTimeout: time.Second}
	defer server.Close()
	go func() { _ = server.Serve(listener) }()
	client := http.Client{Timeout: 3 * time.Second}
	for _, url := range []string{"http://" + listener.Addr().String(), os.Getenv("SAUM_TEST_NETWORK_URL")} {
		response, err := client.Get(url)
		if err != nil {
			t.Fatal(err)
		}
		_ = response.Body.Close()
		if response.StatusCode != http.StatusOK {
			t.Fatal(response.Status)
		}
	}
}
