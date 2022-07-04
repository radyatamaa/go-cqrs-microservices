package internal

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/i18n"
)

type BaseController struct {
	beego.Controller
	i18n.Locale
}

// SetLangVersion sets site language version.
func (r *BaseController) SetLangVersion() {
	// 1. Check URL arguments.
	lang := r.Ctx.Input.Query("lang")

	// Check again in case someone modifies on purpose.
	if !i18n.IsExist(lang) {
		lang = ""
	}

	// 2. Get language information from 'Accept-Language'.
	if len(lang) == 0 {
		al := r.Ctx.Request.Header.Get("Accept-Language")
		if i18n.IsExist(al) {
			lang = al
		}
	}

	// 3. Default language is Indonesia.
	if len(lang) == 0 {
		lang = "id"
	}

	// Set language properties.
	r.Lang = lang
}
