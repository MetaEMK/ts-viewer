package tsviewer

import (
	"context"
	"fmt"

	"github.com/MetaEMK/ts-viewer/internal/config"
)

// Service handles business logic for TeamSpeak operations
type Service struct {
	defaultProvider Provider
	config          *config.Config
}

// NewService creates a new TeamSpeak service instance
func NewService(defaultProvider Provider, cfg *config.Config) *Service {
	return &Service{
		defaultProvider: defaultProvider,
		config:          cfg,
	}
}

// GetServerOverview retrieves the server overview from the default provider
func (s *Service) GetServerOverview(ctx context.Context) (*ServerOverview, error) {
	return s.defaultProvider.FetchOverview(ctx)
}

// GetServerOverviewByName retrieves the server overview from a configured TeamSpeak server
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - serverName: The name of the server as defined in the configuration
//
// Returns the server overview or an error if connection fails or server is not found
func (s *Service) GetServerOverviewByName(ctx context.Context, serverName string) (*ServerOverview, error) {
	// Get server configuration
	serverCfg, ok := s.config.GetServer(serverName)
	if !ok {
		return nil, fmt.Errorf("server '%s' not found in configuration", serverName)
	}

	// Create provider for this specific server
	provider := NewTeamSpeakProvider(serverCfg.Host, serverCfg.Port, serverCfg.Sid)

	// Fetch overview
	return provider.FetchOverview(ctx)
}

// ListServers returns a list of configured server names
func (s *Service) ListServers() []string {
	names := make([]string, 0, len(s.config.Servers))
	for name := range s.config.Servers {
		names = append(names, name)
	}
	return names
}

// GetServersOverview returns information about all configured servers
func (s *Service) GetServersOverview() *ServersOverview {
	servers := make([]ServerInfo, 0, len(s.config.Servers))
	for name, cfg := range s.config.Servers {
		servers = append(servers, ServerInfo{
			Name: name,
			Host: cfg.Host,
			Port: cfg.Port,
			Sid:  cfg.Sid,
		})
	}
	return &ServersOverview{
		Servers: servers,
	}
}
