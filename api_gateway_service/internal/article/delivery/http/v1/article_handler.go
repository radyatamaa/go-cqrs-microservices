package v1

import (
	"context"
	"errors"
	"net/http"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/radyatamaa/go-cqrs-microservices/api_gateway_service/internal"
	"github.com/radyatamaa/go-cqrs-microservices/api_gateway_service/internal/domain"
	"github.com/radyatamaa/go-cqrs-microservices/pkg/response"
	"github.com/radyatamaa/go-cqrs-microservices/pkg/validator"
	"github.com/radyatamaa/go-cqrs-microservices/pkg/zaplogger"
	"gorm.io/gorm"
)

type ArticleHandler struct {
	ZapLogger zaplogger.Logger
	internal.BaseController
	response.ApiResponse
	ArticleUsecase domain.ArticleUseCase
}

func NewArticleHandler(articleUsecase domain.ArticleUseCase, zapLogger zaplogger.Logger) {
	pHandler := &ArticleHandler{
		ZapLogger:      zapLogger,
		ArticleUsecase: articleUsecase,
	}
	beego.Router("/api/v1/articles", pHandler, "post:CreateArticle")
	beego.Router("/api/v1/articles", pHandler, "get:GetArticles")
}

func (h *ArticleHandler) Prepare() {
	// check user access when needed
	h.SetLangVersion()
}

// CreateArticle
// @Title Create Article
// @Tags Article
// @Summary Create Data Article
// @Produce json
// @Param Accept-Language header string false "lang"
// @Success 200 {object} swagger.BaseResponse{errors=[]object,data=domain.CreateArticleRequest}
// @Failure 400 {object} swagger.BadRequestErrorValidationResponse{errors=[]swagger.ValidationErrors,data=object}
// @Failure 408 {object} swagger.RequestTimeoutResponse{errors=[]object,data=object}
// @Failure 500 {object} swagger.InternalServerErrorResponse{errors=[]object,data=object}
// @Param body body domain.CreateArticleRequest true "request payload"
// @Router /v1/articles [post]
func (h *ArticleHandler) CreateArticle() {
	var request domain.CreateArticleRequest

	if err := h.BindJSON(&request); err != nil {
		h.Ctx.Input.SetData("stackTrace", h.ZapLogger.SetMessageLog(err))
		h.ResponseError(h.Ctx, http.StatusBadRequest, response.ApiValidationCodeError, response.ErrorCodeText(response.ApiValidationCodeError, h.Locale.Lang), err)
		return
	}
	if err := validator.Validate.ValidateStruct(&request); err != nil {
		h.Ctx.Input.SetData("stackTrace", h.ZapLogger.SetMessageLog(err))
		h.ResponseError(h.Ctx, http.StatusBadRequest, response.ApiValidationCodeError, response.ErrorCodeText(response.ApiValidationCodeError, h.Locale.Lang), err)
		return
	}

	err := h.ArticleUsecase.CreateArticle(h.Ctx, request)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			h.ResponseError(h.Ctx, http.StatusRequestTimeout, response.RequestTimeoutCodeError, response.ErrorCodeText(response.RequestTimeoutCodeError, h.Locale.Lang), err)
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.ResponseError(h.Ctx, http.StatusBadRequest, response.DataNotFoundCodeError, response.ErrorCodeText(response.DataNotFoundCodeError, h.Locale.Lang), err)
			return
		}
		h.ResponseError(h.Ctx, http.StatusInternalServerError, response.ServerErrorCode, response.ErrorCodeText(response.ServerErrorCode, h.Locale.Lang), err)
		return
	}
	h.Ok(h.Ctx, h.Tr("message.success"), request)
	return
}

// GetArticles
// @Title Get All Articles
// @Tags Article
// @Summary Get All Articles
// @Produce json
// @Param Accept-Language header string false "lang"
// @Param size query int false "size"
// @Param page query int false "page"
// @Param search query string false "search by body or title"
// @Param author query string false "filter by author"
// @Success 200 {object} swagger.BaseResponse{data=[]domain.ArticlePaginationResponse,errors=[]object}
// @Failure 408 {object} swagger.RequestTimeoutResponse{errors=[]object,data=object}
// @Failure 500 {object} swagger.InternalServerErrorResponse{errors=[]object,data=object}
// @router /v1/articles [get]
func (h *ArticleHandler) GetArticles() {
	pageSize, page, err := domain.PaginationQueryParamValidation(h.Ctx.Input.Query("size"), h.Ctx.Input.Query("page"))
	if err != nil {
		h.Ctx.Input.SetData("stackTrace", h.ZapLogger.SetMessageLog(err))
		h.ResponseError(h.Ctx, http.StatusBadRequest, response.QueryParamInvalidCode, response.ErrorCodeText(response.QueryParamInvalidCode, h.Locale.Lang), err)
		return
	}

	result, err := h.ArticleUsecase.GetArticles(h.Ctx, page, pageSize, h.Ctx.Input.Query("search"), h.Ctx.Input.Query("author"))
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			h.ResponseError(h.Ctx, http.StatusRequestTimeout, response.RequestTimeoutCodeError, response.ErrorCodeText(response.RequestTimeoutCodeError, h.Locale.Lang), err)
			return
		}
		h.ResponseError(h.Ctx, http.StatusInternalServerError, response.ServerErrorCode, response.ErrorCodeText(response.ServerErrorCode, h.Locale.Lang), err)
		return
	}
	h.Ok(h.Ctx, h.Tr("message.success"), result)
	return
}
