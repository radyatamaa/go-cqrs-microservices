package helper

import (
	"reflect"

	beegoContext "github.com/beego/beego/v2/server/web/context"
	"github.com/beego/i18n"
)

func ItemExists(arrayType interface{}, item interface{}) bool {
	arr := reflect.ValueOf(arrayType)
	for i := 0; i < arr.Len(); i++ {
		if arr.Index(i).Interface() == item {
			return true
		}
	}
	return false
}

// GetLangVersion sets site language version.
func GetLangVersion(ctx *beegoContext.Context) string {
	// 1. Check URL arguments.
	lang := ctx.Input.Query("lang")

	// Check again in case someone modifies on purpose.
	if !i18n.IsExist(lang) {
		lang = ""
	}

	// 2. Get language information from 'Accept-Language'.
	if len(lang) == 0 {
		al := ctx.Request.Header.Get("Accept-Language")
		if i18n.IsExist(al) {
			lang = al
		}
	}

	// 3. Default language is Indonesia.
	if len(lang) == 0 {
		lang = "id"
	}

	// Set language properties.
	return lang
}
