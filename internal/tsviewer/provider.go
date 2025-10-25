package tsviewer

import "context"

// Provider defines the interface for fetching TeamSpeak server data
type Provider interface {
	FetchOverview(ctx context.Context) (*ServerOverview, error)
}
