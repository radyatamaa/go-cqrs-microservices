package response

import (
	"errors"

	"github.com/beego/i18n"
)

/*
Rules penulisan error code

Format : Format XXXX-XXX-XXX
- 4 digit pertama adalah nama service / aplikasi.
- 3 digit selanjutnya adalah sub service / module tersebut.
- 3 digit terakhir adalah kode unik dari error tersebut.

- Contoh : CORE-AUTH-001
		   CORE-API-001
           CORE-KDM-001
*/

const (
	ApiKeyNotRegisteredCodeError    = "KDMU-AUTH-001"
	MissingApiKeyCodeError          = "KDMU-AUTH-002"
	InvalidApiKeyCodeError          = "KDMU-AUTH-003"
	UnauthorizedCodeError           = "KDMU-AUTH-004"
	RequestForbiddenCodeError       = "KDMU-API-001"
	ResourceNotFoundCodeError       = "KDMU-API-002"
	RequestTimeoutCodeError         = "KDMU-API-003"
	ApiValidationCodeError          = "KDMU-API-004"
	DataNotFoundCodeError           = "KDMU-API-005"
	ServiceCommunicationErrorCode   = "KDMU-API-006"
	InvalidCredentialCodeError      = "KDMU-API-007"
	InvalidTokenCodeError           = "KDMU-API-008"
	ExpiredTokenCodeError           = "KDMU-API-009"
	MissingTokenCodeError           = "KDMU-API-010"
	AuthElseWhereCodeError          = "KDMU-API-011"
	NotAllowedTransaction           = "KDMU-API-012"
	TransactionAlreadyExist         = "KDMU-API-013"
	TransactionRejected             = "KDMU-API-014"
	TransactionNotFound             = "KDMU-API-015"
	InsufficientLimit               = "KDMU-API-016"
	InvalidReturnAmount             = "KDMU-API-017"
	DataAlreadyExistCodeError       = "KDMU-API-018"
	InvalidMinMax                   = "KDMU-API-019"
	InvalidActiveDate               = "KDMU-API-020"
	CustomerStatusNotFoundErrorCode = "KDMU-API-021"
	LimitStatusNotFoundErrorCode    = "KDMU-API-022"
	CustomerIDNotFoundErrorCode     = "KDMU-API-023"
	TenorIDNotFoundErrorCode        = "KDMU-API-024"
	InvalidActiveEndDate            = "KDMU-API-025"
	QueryParamInvalidCode           = "KDMU-API-026"
	PathParamInvalidCode            = "KDMU-API-027"
	ServerErrorCode                 = "KDMU-API-999"
)

var (
	ErrInsufficientLimit         = errors.New("insufficient limit")
	ErrRejectTransaction         = errors.New("transaction rejected")
	ErrTransactionNotAllowed     = errors.New("not allowed transaction")
	ErrTransactionAlreadyExist   = errors.New("transaction already exist")
	ErrTransactionNotFound       = errors.New("transaction not found")
	ErrInvalidReturnAmount       = errors.New("invalid return amount")
	ErrDataAlreadyExist          = errors.New("data already exist")
	ErrMinMoreThanMax            = errors.New("minimal can't be more than maximal amount")
	ErrActiveMoreThanExpired     = errors.New("active date can't be more than expired date")
	ErrActiveMoreThanEnd         = errors.New("active date can't be more than end date")
	ErrCustomerStatusNotFound    = errors.New("customer_status_id not found")
	ErrLimitStatusNotFound       = errors.New("limit_status_id not found")
	ErrQueryParamInvalid         = errors.New("query param is invalid")
	ErrPathParamInvalid          = errors.New("path param is invalid")
	ErrCustomerIDNotFound        = errors.New("customer_id not found")
	ErrTenorIDNotFound           = errors.New("tenor id not found")
	ErrServiceCommunicationError = errors.New("service communication error")
)

func ErrorCodeText(code, locale string, args ...interface{}) string {
	switch code {
	case ApiKeyNotRegisteredCodeError:
		return i18n.Tr(locale, "message.errorApiKeyNotRegistered", args)
	case MissingApiKeyCodeError:
		return i18n.Tr(locale, "message.errorMissingApiKey", args)
	case ApiValidationCodeError:
		return i18n.Tr(locale, "message.errorValidation", args)
	case InvalidApiKeyCodeError:
		return i18n.Tr(locale, "message.errorInvalidApiKey", args)
	case UnauthorizedCodeError:
		return i18n.Tr(locale, "message.errorUnauthorized", args)
	case RequestForbiddenCodeError:
		return i18n.Tr(locale, "message.errorRequestForbidden", args)
	case ResourceNotFoundCodeError:
		return i18n.Tr(locale, "message.errorResourceNotFound", args)
	case ServerErrorCode:
		return i18n.Tr(locale, "message.errorServerError", args)
	case RequestTimeoutCodeError:
		return i18n.Tr(locale, "message.errorRequestTimeout", args)
	case InvalidCredentialCodeError:
		return i18n.Tr(locale, "message.errorInvalidCredential", args)
	case DataNotFoundCodeError:
		return i18n.Tr(locale, "message.errorDataNotFound", args)
	case InvalidTokenCodeError:
		return i18n.Tr(locale, "message.errorInvalidToken", args)
	case ExpiredTokenCodeError:
		return i18n.Tr(locale, "message.errorExpiredToken", args)
	case MissingTokenCodeError:
		return i18n.Tr(locale, "message.errorMissingToken", args)
	case AuthElseWhereCodeError:
		return i18n.Tr(locale, "message.errorAuthElseWhere", args)
	case NotAllowedTransaction:
		return i18n.Tr(locale, "message.errorNotAllowedTransaction", args)
	case TransactionAlreadyExist:
		return i18n.Tr(locale, "message.errorTransactionAlreadyExist", args)
	case TransactionRejected:
		return i18n.Tr(locale, "message.errorTransactionRejected", args)
	case TransactionNotFound:
		return i18n.Tr(locale, "message.errorTransactionNotFound", args)
	case InsufficientLimit:
		return i18n.Tr(locale, "message.errorInsufficientLimit", args)
	case InvalidReturnAmount:
		return i18n.Tr(locale, "message.errorInvalidReturnAmount", args)
	case DataAlreadyExistCodeError:
		return i18n.Tr(locale, "message.errorDataAlreadyExist", args)
	case InvalidMinMax:
		return i18n.Tr(locale, "message.errorInvalidMinMax", args)
	case InvalidActiveDate:
		return i18n.Tr(locale, "message.errorActiveMoreThanExpired", args)
	case InvalidActiveEndDate:
		return i18n.Tr(locale, "message.errorActiveMoreThanEnd", args)
	case CustomerStatusNotFoundErrorCode:
		return i18n.Tr(locale, "message.errorCustomerStatusNotFound", args)
	case LimitStatusNotFoundErrorCode:
		return i18n.Tr(locale, "message.errorLimitStatusNotFound", args)
	case CustomerIDNotFoundErrorCode:
		return i18n.Tr(locale, "message.errorCustomerIDNotFound", args)
	case TenorIDNotFoundErrorCode:
		return i18n.Tr(locale, "message.errorTenorIDNotFound", args)
	case QueryParamInvalidCode:
		return i18n.Tr(locale, "message.errorQueryParamInvalid", args)
	case PathParamInvalidCode:
		return i18n.Tr(locale, "message.errorPathParamInvalid", args)
	default:
		return ""
	}
}
