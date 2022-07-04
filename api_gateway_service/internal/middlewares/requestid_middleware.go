package middlewares

import (
	beego "github.com/beego/beego/v2/server/web"
	beegoContext "github.com/beego/beego/v2/server/web/context"
	"github.com/google/uuid"
)

type (
	// RequestIDConfig defines the config for RequestID middleware.
	RequestIDConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper Skipper

		// Generator defines a function to generate an ID.
		// Optional. Default value random.String(32).
		Generator func() string

		// RequestIDHandler defines a function which is executed for a request id.
		RequestIDHandler func(*beegoContext.Context, string)
	}
)

var (
	// DefaultRequestIDConfig is the default RequestID middleware config.
	DefaultRequestIDConfig = RequestIDConfig{
		Skipper:   DefaultSkipper,
		Generator: generator,
	}
)

// RequestID returns a X-Request-ID middleware.
func RequestID() beego.FilterChain {
	return RequestIDWithConfig(DefaultRequestIDConfig)
}

// RequestIDWithConfig returns a X-Request-ID middleware with config.
func RequestIDWithConfig(config RequestIDConfig) beego.FilterChain  {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultRequestIDConfig.Skipper
	}
	if config.Generator == nil {
		config.Generator = generator
	}

	return func(next beego.FilterFunc) beego.FilterFunc {
		return func(ctx *beegoContext.Context) {
			if config.Skipper(ctx) {
				next(ctx)
				return
			}

			req := ctx.Request
			res := ctx.ResponseWriter.ResponseWriter
			rid := req.Header.Get("X-REQUEST-ID")
			if rid == "" {
				rid = config.Generator()
			}
			res.Header().Set("X-REQUEST-ID", rid)
			if config.RequestIDHandler != nil {
				config.RequestIDHandler(ctx, rid)
			}
			next(ctx)
		}
	}
}

func generator() string {
	return uuid.New().String()
}

