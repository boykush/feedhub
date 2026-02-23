package provider

import (
	"fmt"

	"github.com/samber/do/v2"

	"github.com/boykush/feedhub/server/feed/internal/infra/ent"
	_ "github.com/lib/pq"
)

// EntClient wraps *ent.Client to implement do.Shutdownable.
type EntClient struct {
	*ent.Client
}

func (c *EntClient) Shutdown() error {
	return c.Close()
}

// ProvideEntClient creates a new ent database client.
func ProvideEntClient(i do.Injector) (*EntClient, error) {
	cfg := do.MustInvoke[Config](i)
	client, err := ent.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}
	return &EntClient{Client: client}, nil
}
