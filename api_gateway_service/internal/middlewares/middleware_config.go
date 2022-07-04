package middlewares

import beegoContext "github.com/beego/beego/v2/server/web/context"

type (
	Skipper func(*beegoContext.Context) bool
)

func DefaultSkipper(*beegoContext.Context) bool {
	return false
}
