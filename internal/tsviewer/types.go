package tsviewer

// ServerOverview represents the complete view of a TeamSpeak server
type ServerOverview struct {
	ServerName string
	Channels   []Channel
}

// ServerInfo represents basic information about a configured server
type ServerInfo struct {
	Name          string // Config name
	ServerName    string // Actual server name from TeamSpeak
	OnlineClients int
	TotalChannels int
	IsOnline      bool
	ErrorMessage  string
}

// ServersOverview represents the list of configured servers
type ServersOverview struct {
	Servers []ServerInfo
}

// Channel represents a TeamSpeak channel
type Channel struct {
	ID       int
	Name     string
	ParentID int
	Clients  []Client
	Children []Channel
}

// Client represents a connected user in a TeamSpeak channel
type Client struct {
	ID       int
	Nickname string
	IsMuted  bool
	IsDeaf   bool
}
