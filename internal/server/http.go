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

// ErrorData holds data for error template rendering
type ErrorData struct {
	Title   string
	Message string
	Details string
}

// New creates a new Server instance with Fiber
func New(cfg *config.Config) (*Server, error) {
	// Create template engine from embedded assets
	engine := html.NewFileSystem(http.FS(assets.FS), ".html")

	// Create Fiber app with template engine
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Create service layer
	service := tsviewer.NewService(cfg)

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

// renderError renders an error page with the given status code and error data
func (s *Server) renderError(c *fiber.Ctx, status int, title, message, details string) error {
	c.Status(status)
	return c.Render("templates/error.tmpl", ErrorData{
		Title:   title,
		Message: message,
		Details: details,
	}, "")
}

// handleIndex renders the main TeamSpeak viewer page
func (s *Server) handleIndex(c *fiber.Ctx) error {
	// Fetch overview from service
	overview, err := s.service.GetServerOverview(c.Context())
	if err != nil {
		log.Printf("Error fetching overview: %v", err)
		return s.renderError(c, fiber.StatusInternalServerError,
			"Internal Server Error",
			"We encountered an error while loading the server overview.",
			err.Error())
	}

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
		return s.renderError(c, fiber.StatusBadRequest,
			"Bad Request",
			"The server name is required but was not provided.",
			"Please specify a server name in the URL path.")
	}

	// Use service layer to fetch overview by server name
	overview, err := s.service.GetServerOverviewByName(c.Context(), serverName)
	if err != nil {
		log.Printf("Error fetching TeamSpeak data for server '%s': %v", serverName, err)
		return s.renderError(c, fiber.StatusInternalServerError,
			"Connection Error",
			fmt.Sprintf("Failed to connect to TeamSpeak server '%s'.", serverName),
			err.Error())
	}

	// Render template
	return c.Render("templates/index.tmpl", overview, "")
}
