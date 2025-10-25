package tsviewer

import "context"

// DummyProvider is a static implementation of Provider that returns hard-coded data
type DummyProvider struct{}

// NewDummyProvider creates a new DummyProvider instance
func NewDummyProvider() *DummyProvider {
	return &DummyProvider{}
}

// FetchOverview returns static dummy data shaped like a real TeamSpeak server
func (d *DummyProvider) FetchOverview(ctx context.Context) (*ServerOverview, error) {
	return &ServerOverview{
		ServerName: "My TeamSpeak Server",
		Channels: []Channel{
			{
				ID:       1,
				Name:     "Lobby",
				ParentID: 0,
				Clients: []Client{
					{
						ID:       101,
						Nickname: "Alice",
						IsMuted:  false,
						IsDeaf:   false,
					},
					{
						ID:       102,
						Nickname: "Bob",
						IsMuted:  true,
						IsDeaf:   false,
					},
				},
				Children: []Channel{},
			},
			{
				ID:       2,
				Name:     "Gaming",
				ParentID: 0,
				Clients: []Client{
					{
						ID:       103,
						Nickname: "Charlie",
						IsMuted:  false,
						IsDeaf:   false,
					},
				},
				Children: []Channel{
					{
						ID:       3,
						Name:     "Squad A",
						ParentID: 2,
						Clients: []Client{
							{
								ID:       104,
								Nickname: "Diana",
								IsMuted:  false,
								IsDeaf:   true,
							},
							{
								ID:       105,
								Nickname: "Eve",
								IsMuted:  true,
								IsDeaf:   true,
							},
						},
						Children: []Channel{},
					},
				},
			},
			{
				ID:       4,
				Name:     "AFK",
				ParentID: 0,
				Clients:  []Client{},
				Children: []Channel{},
			},
		},
	}, nil
}
