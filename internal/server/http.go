package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"

	"github.com/MetaEMK/ts-viewer/internal/assets"
	"github.com/MetaEMK/ts-viewer/internal/tsviewer"
)

// Server wraps the HTTP server and dependencies
type Server struct {
	provider tsviewer.Provider
	tmpl     *template.Template
}

// New creates a new Server instance
func New(provider tsviewer.Provider) (*Server, error) {
	// Parse embedded templates
	tmpl, err := template.ParseFS(assets.FS, "templates/*.html")
	if err != nil {
		return nil, err
	}

	return &Server{
		provider: provider,
		tmpl:     tmpl,
	}, nil
}

// Handler returns the HTTP handler with all routes configured
func (s *Server) Handler() (http.Handler, error) {
	mux := http.NewServeMux()

	// Serve static files
	staticFS, err := fs.Sub(assets.FS, "static")
	if err != nil {
		return nil, fmt.Errorf("failed to extract static assets subdirectory: %w", err)
	}
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	// Routes
	mux.HandleFunc("/", s.handleIndex)
	mux.HandleFunc("/healthz", s.handleHealthz)

	return mux, nil
}

// handleIndex renders the main TeamSpeak viewer page
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Fetch overview from provider
	overview, err := s.provider.FetchOverview(r.Context())
	if err != nil {
		log.Printf("Error fetching overview: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Render template
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := s.tmpl.ExecuteTemplate(w, "index.tmpl.html", overview); err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// handleHealthz returns a simple health check response
func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}
