package mongodb

import (
	"context"
	"time"
)

type Config struct {
	Uri     string
	Timeout time.Duration
}

func DefaultConfig(uri string) Config {
	return Config{
		Uri:     uri,
		Timeout: 5 * time.Second,
	}
}

// ---

func (c Config) GetContext(ctx context.Context) (context.Context, context.CancelFunc) {
	if c.Timeout == 0 {
		return context.WithCancel(ctx)
	}
	return context.WithTimeout(ctx, c.Timeout)
}
