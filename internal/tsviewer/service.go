package tsviewer

import (
	"context"
	"fmt"

	"github.com/MetaEMK/ts-viewer/internal/config"
)

// Service handles business logic for TeamSpeak operations
type Service struct {
	config *config.Config
}

// NewService creates a new TeamSpeak service instance
func NewService(cfg *config.Config) *Service {
	return &Service{
		config: cfg,
	}
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
// This fetches live data from each server (name, online clients, channels)
func (s *Service) GetServersOverview(ctx context.Context) *ServersOverview {
	servers := make([]ServerInfo, 0, len(s.config.Servers))
	for name, cfg := range s.config.Servers {
		serverInfo := ServerInfo{
			Name:     name,
			IsOnline: false,
		}

		// Try to fetch live data from the server
		provider := NewTeamSpeakProvider(cfg.Host, cfg.Port, cfg.Sid)
		overview, err := provider.FetchOverview(ctx)
		if err != nil {
			// Server is offline or unreachable
			serverInfo.ErrorMessage = err.Error()
		} else {
			// Server is online, collect stats
			serverInfo.IsOnline = true
			serverInfo.ServerName = overview.ServerName
			serverInfo.OnlineClients = countClients(overview.Channels)
			serverInfo.TotalChannels = countChannels(overview.Channels)
		}

		servers = append(servers, serverInfo)
	}
	return &ServersOverview{
		Servers: servers,
	}
}

// countClients recursively counts all clients in channels
func countClients(channels []Channel) int {
	count := 0
	for _, ch := range channels {
		count += len(ch.Clients)
		count += countClients(ch.Children)
	}
	return count
}

// countChannels recursively counts all channels
func countChannels(channels []Channel) int {
	count := len(channels)
	for _, ch := range channels {
		count += countChannels(ch.Children)
	}
	return count
}
