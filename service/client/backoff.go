package client

import (
	"context"
	"time"
)

type BackoffFunc func(ctx context.Context, req Request, attempts int) (time.Duration, error)
