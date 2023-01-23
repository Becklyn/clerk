package postgres

import (
	"context"
	"time"
)

type Config struct {
	Host    Host
	Timeout time.Duration
}

func DefaultConfig(host Host) Config {
	return Config{
		Host:    host,
		Timeout: 5 * time.Second,
	}
}

func (c Config) GetContext(ctx context.Context) (context.Context, context.CancelFunc) {
	if c.Timeout == 0 {
		return context.WithCancel(ctx)
	}
	return context.WithTimeout(ctx, c.Timeout)
}
