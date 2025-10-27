package tsviewer

import (
	"context"
	"fmt"
	"strconv"
)

// Service handles business logic for TeamSpeak operations
type Service struct {
	defaultProvider Provider
}

// NewService creates a new TeamSpeak service instance
func NewService(defaultProvider Provider) *Service {
	return &Service{
		defaultProvider: defaultProvider,
	}
}

// GetServerOverview retrieves the server overview from the default provider
func (s *Service) GetServerOverview(ctx context.Context) (*ServerOverview, error) {
	return s.defaultProvider.FetchOverview(ctx)
}

// GetServerOverviewByAddress retrieves the server overview from a specific TeamSpeak server
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - host: The hostname or IP address of the TeamSpeak server
//   - portStr: The port number as a string (optional, defaults to 10011)
//
// Returns the server overview or an error if connection fails
func (s *Service) GetServerOverviewByAddress(ctx context.Context, host string, portStr string) (*ServerOverview, error) {
	// Validate host
	if host == "" {
		return nil, fmt.Errorf("host cannot be empty")
	}

	// Parse port if provided
	port := 0 // 0 will default to 10011 in the provider
	if portStr != "" {
		parsedPort, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, fmt.Errorf("invalid port number: %w", err)
		}
		if parsedPort < 1 || parsedPort > 65535 {
			return nil, fmt.Errorf("port must be between 1 and 65535")
		}
		port = parsedPort
	}

	// Create provider for this specific server
	provider := NewTeamSpeakProvider(host, port)

	// Fetch overview
	return provider.FetchOverview(ctx)
}
