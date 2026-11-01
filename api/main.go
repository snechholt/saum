package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"unicode"
	"unicode/utf8"
)

// A project is only a named container. It has no internal workflow yet.
type Project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func parseName(raw string) (string, error) {
	name := strings.TrimSpace(raw)
	if !utf8.ValidString(name) || utf8.RuneCountInString(name) < 1 || utf8.RuneCountInString(name) > 80 {
		return "", fmt.Errorf("Project name must contain 1–80 characters")
	}
	for _, character := range name {
		if unicode.IsControl(character) {
			return "", fmt.Errorf("Project name cannot contain control characters")
		}
	}
	return name, nil
}
func jsonResponse(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func (f *Factory) handler(web http.Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/projects", func(w http.ResponseWriter, r *http.Request) {
		projects, err := f.list()
		if err != nil {
			jsonResponse(w, 500, map[string]string{"error": "Could not load projects"})
			return
		}
		jsonResponse(w, 200, projects)
	})
	mux.HandleFunc("POST /api/projects", func(w http.ResponseWriter, r *http.Request) {
		if strings.Split(r.Header.Get("Content-Type"), ";")[0] != "application/json" {
			jsonResponse(w, 415, map[string]string{"error": "Expected application/json"})
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, 4096)
		var input struct {
			Name string `json:"name"`
		}
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if decoder.Decode(&input) != nil || decoder.Decode(new(any)) != io.EOF {
			jsonResponse(w, 400, map[string]string{"error": "Invalid project request"})
			return
		}
		name, err := parseName(input.Name)
		if err != nil {
			jsonResponse(w, 400, map[string]string{"error": err.Error()})
			return
		}
		project, err := f.create(name)
		if err != nil {
			jsonResponse(w, 500, map[string]string{"error": "Could not save project"})
			return
		}
		jsonResponse(w, 201, project)
	})
	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) { http.NotFound(w, r) })
	mux.Handle("/", web)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Cache-Control", "no-store")
		if r.Host != "127.0.0.1:8080" && r.Host != "localhost:8080" {
			http.Error(w, "Local requests only", 403)
			return
		}
		origin := r.Header.Get("Origin")
		if origin != "" && origin != "http://"+r.Host && origin != "http://127.0.0.1:5173" {
			http.Error(w, "Cross-origin request rejected", 403)
			return
		}
		if r.Header.Get("Sec-Fetch-Site") == "cross-site" {
			http.Error(w, "Cross-site request rejected", 403)
			return
		}
		mux.ServeHTTP(w, r)
	})
}
func main() {
	data := flag.String("data", "../.saum/factory.sqlite", "SQLite database file")
	flag.Parse()
	factory, err := openFactory(*data)
	if err != nil {
		log.Fatal(err)
	}
	defer factory.db.Close()
	server := &http.Server{Addr: "127.0.0.1:8080", Handler: factory.handler(http.FileServer(http.Dir("../sites/dashboard/dist"))), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second, IdleTimeout: 60 * time.Second}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() {
		<-ctx.Done()
		shutdown, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdown)
	}()
	log.Print("Saum API → http://127.0.0.1:8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
