package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MetaEMK/ts-viewer/internal/config"
	"github.com/MetaEMK/ts-viewer/internal/server"
	"github.com/MetaEMK/ts-viewer/internal/tsviewer"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	log.Printf("Starting TeamSpeak Viewer with config: %s", cfg)

	// Create provider (currently dummy)
	provider := tsviewer.NewDummyProvider()

	// Create HTTP server
	srv, err := server.New(provider)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Setup HTTP server
	httpServer := &http.Server{
		Addr:         cfg.HTTPAddr,
		Handler:      srv.Handler(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server listening on %s", cfg.HTTPAddr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped gracefully")
}
