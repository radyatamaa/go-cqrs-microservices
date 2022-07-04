package middlewares

import (
	"bufio"
	"bytes"
	beego "github.com/beego/beego/v2/server/web"
	beegoContext "github.com/beego/beego/v2/server/web/context"
	"io"
	"io/ioutil"
	"net"
	"net/http"
)


type (
	// BodyDumpConfig defines the config for BodyDump middleware.
	BodyDumpConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper Skipper
		// Handler receives request and response payload.
		// Required.
		Handler BodyDumpHandler
	}

	// BodyDumpHandler receives the request and response payload.
	BodyDumpHandler func(*beegoContext.Context, []byte, []byte)

	bodyDumpResponseWriter struct {
		io.Writer
		http.ResponseWriter
	}
)

var (
	// DefaultBodyDumpConfig is the default BodyDump middleware config.
	DefaultBodyDumpConfig = BodyDumpConfig{
		Skipper: DefaultSkipper,
	}
)

// BodyDump returns a BodyDump middleware.
//
// BodyDump middleware captures the request and response payload and calls the
// registered handler.
func BodyDump(handler BodyDumpHandler) beego.FilterChain {
	c := DefaultBodyDumpConfig
	c.Handler = handler
	return BodyDumpWithConfig(c)
}

// BodyDumpWithConfig returns a BodyDump middleware with config.
func BodyDumpWithConfig(config BodyDumpConfig) beego.FilterChain {
	// Defaults
	if config.Handler == nil {
		panic("body-dump middleware requires a handler function")
	}
	if config.Skipper == nil {
		config.Skipper = DefaultBodyDumpConfig.Skipper
	}

	return func(next beego.FilterFunc) beego.FilterFunc {
		return func(ctx *beegoContext.Context) {
			if config.Skipper(ctx) {
				next(ctx)
				return
			}

			var reqBody []byte
			if ctx.Request.Body != nil { // Read
				reqBody, _ = ioutil.ReadAll(ctx.Request.Body)
			}
			ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody)) // Reset


			// Response
			resBody := new(bytes.Buffer)
			mw := io.MultiWriter(ctx.ResponseWriter.ResponseWriter, resBody)
			writer := &bodyDumpResponseWriter{Writer: mw, ResponseWriter: ctx.ResponseWriter.ResponseWriter}
			ctx.ResponseWriter.ResponseWriter = writer

			next(ctx)


			// Callback
			config.Handler(ctx, reqBody, resBody.Bytes())

			return
		}
	}
}

func (w *bodyDumpResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
}

func (w *bodyDumpResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (w *bodyDumpResponseWriter) Flush() {
	w.ResponseWriter.(http.Flusher).Flush()
}

func (w *bodyDumpResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}
