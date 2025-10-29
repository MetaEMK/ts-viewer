package tsviewer

import (
	"context"
	"fmt"
	"time"

	"github.com/multiplay/go-ts3"
)

// TeamSpeakProvider connects to real TeamSpeak servers via ServerQuery
type TeamSpeakProvider struct {
	host string
	port int
	sid  int
}

// NewTeamSpeakProvider creates a new TeamSpeak provider for the given host
// If port is 0, it defaults to 10011 (standard ServerQuery port)
// If sid is 0, it defaults to 1 (first virtual server)
func NewTeamSpeakProvider(host string, port int, sid int) *TeamSpeakProvider {
	if port == 0 {
		port = 10011
	}
	if sid == 0 {
		sid = 1
	}
	return &TeamSpeakProvider{
		host: host,
		port: port,
		sid:  sid,
	}
}

// FetchOverview connects to the TeamSpeak server and fetches the current state
func (t *TeamSpeakProvider) FetchOverview(ctx context.Context) (*ServerOverview, error) {
	// Check if context is already cancelled
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	// Create client connection
	addr := fmt.Sprintf("%s:%d", t.host, t.port)
	client, err := ts3.NewClient(addr, ts3.Timeout(10*time.Second))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to TeamSpeak server at %s: %w", addr, err)
	}
	defer client.Close()

	// Use the configured virtual server
	if err := client.Use(t.sid); err != nil {
		return nil, fmt.Errorf("failed to select virtual server %d: %w", t.sid, err)
	}

	// Fetch server info
	serverInfo, err := client.Server.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch server info: %w", err)
	}

	// Fetch channel list
	channels, err := client.Server.ChannelList()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch channel list: %w", err)
	}

	// Fetch online clients with voice information (for mute/deaf status)
	clients, err := client.Server.ClientList("-voice")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch client list: %w", err)
	}

	// Build the overview
	overview := &ServerOverview{
		ServerName: serverInfo.Name,
		Channels:   t.buildChannelTree(channels, clients, 0),
	}

	return overview, nil
}

// buildChannelTree recursively builds the channel tree structure
func (t *TeamSpeakProvider) buildChannelTree(channels []*ts3.Channel, clients []*ts3.OnlineClient, parentID int) []Channel {
	var result []Channel

	for _, ch := range channels {
		if ch.ParentID != parentID {
			continue
		}

		channel := Channel{
			ID:       ch.ID,
			Name:     ch.ChannelName,
			ParentID: ch.ParentID,
			Clients:  t.getClientsInChannel(clients, ch.ID),
			Children: t.buildChannelTree(channels, clients, ch.ID),
		}

		result = append(result, channel)
	}

	return result
}

// getClientsInChannel returns all clients in the specified channel
func (t *TeamSpeakProvider) getClientsInChannel(clients []*ts3.OnlineClient, channelID int) []Client {
	var result []Client

	for _, c := range clients {
		// Skip ServerQuery clients (type 1)
		if c.Type == 1 {
			continue
		}

		if c.ChannelID == channelID {
			client := Client{
				ID:       c.ID,
				Nickname: c.Nickname,
				IsMuted:  false,
				IsDeaf:   false,
			}

			// Check mute/deaf status if available
			if c.OnlineClientExt != nil && c.OnlineClientExt.OnlineClientVoice != nil {
				if c.OnlineClientExt.OnlineClientVoice.InputMuted != nil && *c.OnlineClientExt.OnlineClientVoice.InputMuted {
					client.IsMuted = true
				}
				if c.OnlineClientExt.OnlineClientVoice.OutputMuted != nil && *c.OnlineClientExt.OnlineClientVoice.OutputMuted {
					client.IsDeaf = true
				}
			}

			result = append(result, client)
		}
	}

	return result
}
