package clerk

import (
	"context"
)

type Config interface {
	GetContext(ctx context.Context) (context.Context, context.CancelFunc)
}
