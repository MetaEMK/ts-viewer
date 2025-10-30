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
	"github.com/MetaEMK/ts-viewer/internal/config"
	"github.com/MetaEMK/ts-viewer/internal/tsviewer"
)

// Server wraps the Fiber app and dependencies
type Server struct {
	service *tsviewer.Service
	app     *fiber.App
}

// New creates a new Server instance with Fiber
func New(provider tsviewer.Provider, cfg *config.Config) (*Server, error) {
	// Create template engine from embedded assets
	engine := html.NewFileSystem(http.FS(assets.FS), ".html")

	// Create Fiber app with template engine
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Create service layer
	service := tsviewer.NewService(provider, cfg)

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
	s.app.Get("/", s.handleServersOverview)
	s.app.Get("/ts-view/:server", s.handleTSView)
	s.app.Get("/healthz", s.handleHealthz)
}

// App returns the Fiber app instance
func (s *Server) App() *fiber.App {
	return s.app
}

// handleServersOverview renders the servers overview page
func (s *Server) handleServersOverview(c *fiber.Ctx) error {
	// Get servers overview from service (with live data)
	overview := s.service.GetServersOverview(c.Context())

	// Render template
	return c.Render("templates/overview.tmpl", overview, "")
}

// handleHealthz returns a simple health check response
func (s *Server) handleHealthz(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "ok",
	})
}

// handleTSView renders the TeamSpeak viewer page for a specific server
// Accepts path parameter: server (the server name as configured in config.yaml)
func (s *Server) handleTSView(c *fiber.Ctx) error {
	// Get the server name from path parameter
	serverName := c.Params("server")

	if serverName == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Missing required parameter: server name")
	}

	// Use service layer to fetch overview by server name
	overview, err := s.service.GetServerOverviewByName(c.Context(), serverName)
	if err != nil {
		log.Printf("Error fetching TeamSpeak data for server '%s': %v", serverName, err)
		return c.Status(fiber.StatusInternalServerError).SendString(
			fmt.Sprintf("Failed to connect to TeamSpeak server '%s': %v", serverName, err),
		)
	}

	// Render template
	return c.Render("templates/index.tmpl", overview, "")
}
