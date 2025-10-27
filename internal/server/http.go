package server

import (
	"fmt"
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
	service *tsviewer.Service
	app     *fiber.App
}

// New creates a new Server instance with Fiber
func New(provider tsviewer.Provider) (*Server, error) {
	// Create template engine from embedded assets
	engine := html.NewFileSystem(http.FS(assets.FS), ".html")

	// Create Fiber app with template engine
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Create service layer
	service := tsviewer.NewService(provider)

	s := &Server{
		service: service,
		app:     app,
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
	s.app.Get("/ts-view", s.handleTSView)
	s.app.Get("/healthz", s.handleHealthz)
}

// App returns the Fiber app instance
func (s *Server) App() *fiber.App {
	return s.app
}

// handleIndex renders the main TeamSpeak viewer page
func (s *Server) handleIndex(c *fiber.Ctx) error {
	// Fetch overview from service
	overview, err := s.service.GetServerOverview(c.Context())
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

// handleTSView renders the TeamSpeak viewer page for a specific server
// Accepts query parameters: ip or host (the TeamSpeak server address), and port (optional)
func (s *Server) handleTSView(c *fiber.Ctx) error {
	// Get the server address from query parameters
	host := c.Query("ip")
	if host == "" {
		host = c.Query("host")
	}

	if host == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Missing required parameter: ip or host")
	}

	// Get optional port parameter
	port := c.Query("port")

	// Use service layer to fetch overview
	overview, err := s.service.GetServerOverviewByAddress(c.Context(), host, port)
	if err != nil {
		log.Printf("Error fetching TeamSpeak data from %s: %v", host, err)
		return c.Status(fiber.StatusInternalServerError).SendString(
			fmt.Sprintf("Failed to connect to TeamSpeak server at %s: %v", host, err),
		)
	}

	// Render template
	return c.Render("templates/index.tmpl", overview, "")
}
