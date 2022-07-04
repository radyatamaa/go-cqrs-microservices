package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/beego/i18n"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	validatorGo "github.com/go-playground/validator/v10"
	"github.com/radyatamaa/go-cqrs-microservices/pkg/validator"
)

type (
	ApiResponseInterface interface {
		Ok(ctx *context.Context, data interface{}) error
		OkWithCustomMessage(ctx *context.Context, message string, data interface{}) error
		ResponseError(ctx *context.Context, httpStatus int, errorCode string, err error) error
		ResponseErrorWithCustomMessage(ctx *context.Context, httpStatus int, errorCode string, message string, err error) error
	}
)

type ApiResponse struct {
	Code      string      `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	Errors    []Errors    `json:"errors"`
	RequestId string      `json:"request_id"`
	Timestamp string      `json:"timestamp"`
}

type Errors struct {
	Field       string `json:"field"`
	Description string `json:"message"`
}

func (r ApiResponse) Ok(ctx *context.Context, message string, data interface{}) error {
	ctx.Output.SetStatus(http.StatusOK)

	return ctx.Output.JSON(ApiResponse{
		Code:      http.StatusText(http.StatusOK),
		RequestId: ctx.ResponseWriter.ResponseWriter.Header().Get("X-REQUEST-ID"),
		Message:   message,
		Data:      data,
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
	}, beego.BConfig.RunMode != "prod", false)
}

func (r ApiResponse) ResponseError(ctx *context.Context, httpStatus int, errorCode string, message string, err error) error {
	var apiResponse ApiResponse
	var errorValidations []Errors = nil

	ctx.Output.SetStatus(httpStatus)

	if err != nil {
		if ctx.Input.RequestBody != nil {
			validateJsonError := checkJsonRequest(err)
			if len(validateJsonError) > 0 {
				errorValidations = validateJsonError
			} else if ute, ok := err.(*validator.ValidateDynamicStructError); ok {
				lang := "id"
				acceptLang := ctx.Request.Header.Get("Accept-Language")
				if i18n.IsExist(acceptLang) {
					lang = acceptLang
				}
				if trans, found := validator.Validate.GetTranslator(lang); found {
					errorValidations = append(errorValidations, Errors{
						Field:       ute.Field,
						Description: ute.Field + ute.Msg.Translate(trans),
					})
				}
			} else {
				if fields, ok := err.(validatorGo.ValidationErrors); ok {
					lang := "id"
					acceptLang := ctx.Request.Header.Get("Accept-Language")
					if i18n.IsExist(acceptLang) {
						lang = acceptLang
					}
					if trans, found := validator.Validate.GetTranslator(lang); found {
						for _, v := range fields {
							errorValidations = append(errorValidations, Errors{
								Field:       v.Field(),
								Description: v.Translate(trans),
							})
						}
					}
				}
			}
		}
	}

	apiResponse.RequestId = ctx.ResponseWriter.ResponseWriter.Header().Get("X-REQUEST-ID")
	apiResponse.Code = errorCode
	apiResponse.Message = message
	apiResponse.Timestamp = time.Now().Format("2006-01-02 15:04:05")
	apiResponse.Errors = errorValidations

	return ctx.Output.JSON(apiResponse, beego.BConfig.RunMode != "prod", false)
}

// checkJsonRequest Response API
func checkJsonRequest(err error) (response []Errors) {
	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError
	var invalidUnmarshalError *json.InvalidUnmarshalError

	switch {
	case errors.As(err, &syntaxError):
		msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
		errorValidation := Errors{
			Field:       "json",
			Description: msg,
		}
		response = append(response, errorValidation)
		return
	case errors.Is(err, io.ErrUnexpectedEOF):
		msg := fmt.Sprintf("Request body contains badly-formed JSON")
		errorValidation := Errors{
			Field:       "json",
			Description: msg,
		}
		response = append(response, errorValidation)
		return
	case errors.As(err, &unmarshalTypeError):
		if ute, ok := err.(*json.UnmarshalTypeError); ok {
			errorValidation := Errors{
				Field:       ute.Field,
				Description: fmt.Sprintf("Parameter %s is invalid (type: %s)", ute.Field, ute.Type),
			}
			response = append(response, errorValidation)
			return
		}
	case errors.As(err, &invalidUnmarshalError):
		if ute, ok := err.(*json.InvalidUnmarshalError); ok {
			errorValidation := Errors{
				Field:       ute.Type.Name(),
				Description: ute.Error(),
			}
			response = append(response, errorValidation)
			return
		}
	}
	return response
}
