package middlewares

import (
	"encoding/json"
	"fmt"
	"strings"

	contextBeego "github.com/beego/beego/v2/server/web/context"
	"github.com/radyatamaa/go-cqrs-microservices/pkg/zaplogger"
)

type AccessLogMiddleware struct {
	ZapLogger  zaplogger.Logger
	AppVersion string
}

func NewAccessLogMiddleware(ZapLogger zaplogger.Logger, appVersion string) *AccessLogMiddleware {
	return &AccessLogMiddleware{
		ZapLogger:  ZapLogger,
		AppVersion: appVersion,
	}
}
func (m *AccessLogMiddleware) Logger() BodyDumpConfig {
	return BodyDumpConfig{
		Skipper: func(ctx *contextBeego.Context) bool {
			if strings.EqualFold(ctx.Request.URL.Path, "/swagger/index.html") {
				return true
			}
			return false
		},
		Handler: func(context *contextBeego.Context, request []byte, response []byte) {
			if context.ResponseWriter.Status > 399 {
				if errorData, ok := context.Input.GetData("stackTrace").(*zaplogger.ListErrors); ok {
					m.ZapLogger.Errorf(zaplogger.StdFormatErrorLog,
						m.AppVersion,
						context.Request.Host,
						context.Request.URL.String(),
						context.ResponseWriter.ResponseWriter.Header().Get("X-REQUEST-ID"),
						json.RawMessage(request),
						json.RawMessage(response),
						fmt.Sprintf("%s %s %s %d %v", errorData.Error, errorData.File, errorData.Function, errorData.Line, errorData.Extra))
				} else {
					m.ZapLogger.Errorf(zaplogger.StdFormatErrorLog,
						m.AppVersion,
						context.Request.Host,
						context.Request.URL.String(),
						context.ResponseWriter.ResponseWriter.Header().Get("X-REQUEST-ID"),
						json.RawMessage(request),
						json.RawMessage(response), "")
				}
			} else {
				m.ZapLogger.Infof(zaplogger.StdFormatLog,
					m.AppVersion,
					context.Request.Host,
					context.Request.URL.String(),
					context.ResponseWriter.ResponseWriter.Header().Get("X-REQUEST-ID"),
					json.RawMessage(request),
					json.RawMessage(response))
			}
		},
	}
}
