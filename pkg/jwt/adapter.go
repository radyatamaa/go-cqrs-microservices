package jwt

import (
	"context"
	"time"
)

// Adapter interface cache redis
type (
	Adapter interface {
		Get(ctx context.Context, key string) (interface{}, error)
		Put(ctx context.Context, key string, val interface{}, timeout time.Duration) error
		Delete(ctx context.Context, key string) error
	}
)


