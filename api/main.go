package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"unicode"
	"unicode/utf8"
)

func parseName(raw string) (string, error) {
	name := strings.TrimSpace(raw)
	if !utf8.ValidString(name) || utf8.RuneCountInString(name) < 1 || utf8.RuneCountInString(name) > 80 {
		return "", fmt.Errorf("Name must contain 1–80 characters")
	}
	for _, character := range name {
		if unicode.IsControl(character) {
			return "", fmt.Errorf("Name cannot contain control characters")
		}
	}
	return name, nil
}
func jsonResponse(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func problem(w http.ResponseWriter, status int, message string) {
	jsonResponse(w, status, map[string]string{"error": message})
}

func (s *Store) handler(web http.Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/factories", s.getFactories)
	mux.HandleFunc("GET /api/settings", s.settingsHTTP)
	mux.HandleFunc("PUT /api/settings", s.settingsHTTP)
	mux.HandleFunc("GET /api/audit", s.auditHTTP)
	mux.HandleFunc("POST /api/factories", s.saveFactory)
	mux.HandleFunc("PUT /api/factories/{factory}", s.saveFactory)
	mux.HandleFunc("DELETE /api/factories/{factory}", s.removeFactory)
	mux.HandleFunc("GET /api/factories/{factory}/stories", s.getStories)
	mux.HandleFunc("POST /api/factories/{factory}/stories", s.postStory)
	mux.Handle("/mcp", s.mcpHandler())
	mux.HandleFunc("GET /api/sessions/{session}", s.getSession)
	mux.HandleFunc("POST /api/sessions/{session}/comments", s.postPlanComment)
	mux.HandleFunc("DELETE /api/sessions/{session}/comments/{comment}", s.deletePlanComment)
	mux.HandleFunc("POST /api/sessions/{session}/approve", s.postPlanApproval)
	mux.HandleFunc("POST /api/sessions/{session}/feedback", s.postPlanFeedback)
	mux.HandleFunc("GET /api/sessions/{session}/events", s.getSessionEvents)
	mux.HandleFunc("DELETE /api/sessions/{session}", s.cancelSession)
	mux.HandleFunc("POST /api/stories/{story}/planning", s.retryPlanning)
	mux.HandleFunc("POST /api/stories/{story}/implementation", s.retryImplementationHTTP)
	mux.HandleFunc("POST /api/work-items/{item}/approve", s.workItemDecisionHTTP(true))
	mux.HandleFunc("POST /api/work-items/{item}/feedback", s.workItemDecisionHTTP(false))
	mux.HandleFunc("POST /api/work-items/{item}/reply", s.replyWorkItemHTTP)
	mux.HandleFunc("GET /api/stories/{story}", s.getStory)
	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) { problem(w, 404, "Not found") })
	mux.Handle("/", web)
	return localOnly(mux)
}
func (s *Store) getFactories(w http.ResponseWriter, r *http.Request) {
	factories, err := s.listFactories()
	if err != nil {
		problem(w, 500, "Could not load factories")
		return
	}
	jsonResponse(w, 200, factories)
}
func (s *Store) saveFactory(w http.ResponseWriter, r *http.Request) {
	if strings.Split(r.Header.Get("Content-Type"), ";")[0] != "application/json" {
		problem(w, 415, "Expected application/json")
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	var input struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if decoder.Decode(&input) != nil || decoder.Decode(new(any)) != io.EOF {
		problem(w, 400, "Invalid factory request")
		return
	}
	name, err := parseName(input.Name)
	if err != nil {
		problem(w, 400, err.Error())
		return
	}
	var factory Factory
	status := http.StatusCreated
	if r.Method == http.MethodPut {
		factory, err = s.updateFactory(r.PathValue("factory"), name)
		status = http.StatusOK
	} else {
		factory, err = s.createFactory(name)
	}
	if errors.Is(err, sql.ErrNoRows) {
		problem(w, 404, "Factory not found")
		return
	}
	if err != nil {
		problem(w, 500, "Could not save factory")
		return
	}
	jsonResponse(w, status, factory)
}
func (s *Store) removeFactory(w http.ResponseWriter, r *http.Request) {
	err := s.deleteFactory(r.PathValue("factory"))
	switch {
	case errors.Is(err, sql.ErrNoRows):
		problem(w, 404, "Factory not found")
	case errors.Is(err, ErrFactoryNotEmpty):
		problem(w, 409, "Factory still contains user stories")
	case err != nil:
		problem(w, 500, "Could not delete factory")
	default:
		w.WriteHeader(http.StatusNoContent)
	}
}

func (s *Store) getStories(w http.ResponseWriter, r *http.Request) {
	stories, err := s.listStories(r.PathValue("factory"))
	if errors.Is(err, sql.ErrNoRows) {
		problem(w, 404, "Factory not found")
		return
	}
	if err != nil {
		problem(w, 500, "Could not load user stories")
		return
	}
	jsonResponse(w, 200, stories)
}
func localOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Cache-Control", "no-store")
		if r.Host != "127.0.0.1:8080" && r.Host != "localhost:8080" {
			problem(w, 403, "Local requests only")
			return
		}
		origin := r.Header.Get("Origin")
		if origin != "" && origin != "http://"+r.Host && origin != "http://127.0.0.1:5173" && origin != "http://localhost:5173" {
			problem(w, 403, "Cross-origin request rejected")
			return
		}
		if r.Header.Get("Sec-Fetch-Site") == "cross-site" {
			problem(w, 403, "Cross-site request rejected")
			return
		}
		next.ServeHTTP(w, r)
	})
}
func main() {
	data := flag.String("data", "../.saum/factory.sqlite", "SQLite database file")
	codex := flag.String("codex", "codex", "Codex executable")
	workerEnabled := flag.Bool("worker", true, "Automatically decompose user stories")
	flag.Parse()
	store, err := openStore(*data)
	if err != nil {
		log.Fatal(err)
	}
	defer store.db.Close()
	planner := &Planner{Store: store, Codex: *codex, MCPURL: "http://127.0.0.1:8080/mcp"}
	store.planner = planner
	server := &http.Server{
		Addr:              "127.0.0.1:8080",
		Handler:           store.handler(dashboardHandler("../sites/dashboard/dist")),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		log.Fatal(err)
	}
	workerDone := make(chan struct{})
	go func() {
		defer close(workerDone)
		if *workerEnabled {
			planner.Run(ctx)
		}
	}()
	go func() {
		<-ctx.Done()
		shutdown, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdown)
	}()
	log.Print("Saum API → http://127.0.0.1:8080")
	if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
		log.Print(err)
	}
	stop()
	<-workerDone
}
