package server

import (
	"io/fs"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/template/html/v2"

	"github.com/MetaEMK/ts-viewer/internal/assets"
	"github.com/MetaEMK/ts-viewer/internal/tsviewer"
)

// Server wraps the Fiber app and dependencies
type Server struct {
	provider tsviewer.Provider
	app      *fiber.App
}

// New creates a new Server instance with Fiber
func New(provider tsviewer.Provider) (*Server, error) {
	// Create template engine from embedded assets
	engine := html.NewFileSystem(http.FS(assets.FS), ".html")

	// Create Fiber app with template engine
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	s := &Server{
		provider: provider,
		app:      app,
	}

	// Setup routes
	s.setupRoutes()

	return s, nil
}

// setupRoutes configures all application routes
func (s *Server) setupRoutes() {
	// Serve static files from embedded assets
	staticFS, err := fs.Sub(assets.FS, "static")
	if err != nil {
		log.Fatalf("Failed to create static filesystem: %v", err)
	}

	s.app.Use("/static", filesystem.New(filesystem.Config{
		Root: http.FS(staticFS),
	}))

	// Routes
	s.app.Get("/", s.handleIndex)
	s.app.Get("/healthz", s.handleHealthz)
}

// App returns the Fiber app instance
func (s *Server) App() *fiber.App {
	return s.app
}

// handleIndex renders the main TeamSpeak viewer page
func (s *Server) handleIndex(c *fiber.Ctx) error {
	// Fetch overview from provider
	overview, err := s.provider.FetchOverview(c.Context())
	if err != nil {
		log.Printf("Error fetching overview: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	// Render template
	return c.Render("templates/index.tmpl", overview, "")
}

// handleHealthz returns a simple health check response
func (s *Server) handleHealthz(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "ok",
	})
}
