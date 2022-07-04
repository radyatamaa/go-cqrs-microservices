package response

import (
	"net/http"

	"github.com/beego/beego/v2/server/web"
	"github.com/radyatamaa/go-cqrs-microservices/pkg/helper"
)

type ErrorController struct {
	web.Controller
	ApiResponse
}

func (c *ErrorController) Error404() {
	c.ResponseError(c.Ctx, http.StatusNotFound, ResourceNotFoundCodeError, ErrorCodeText(ResourceNotFoundCodeError, helper.GetLangVersion(c.Ctx)), nil)
	return
}

func (c *ErrorController) Error500() {
	c.ResponseError(c.Ctx, http.StatusInternalServerError, ServerErrorCode, ErrorCodeText(ServerErrorCode, helper.GetLangVersion(c.Ctx)), nil)
	return
}
