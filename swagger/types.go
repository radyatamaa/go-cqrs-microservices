package swagger

type BaseResponse struct {
	Code      string      `json:"code" example:"OK"`
	Message   string      `json:"message" example:"operasi berhasil dieksekusi."`
	Data      interface{} `json:"data"`
	Errors    interface{} `json:"errors"`
	RequestId string      `json:"request_id" example:"24fa3770-628c-49de-aa17-3a338f73d99b"`
	Timestamp string      `json:"timestamp" example:"2022-04-27 23:19:56"`
}

type BadRequestResponse struct {
	Code      string      `json:"code" example:"KDMU-02-011"`
	Message   string      `json:"message" example:"data yang anda minta tidak ditemukan."`
	Data      interface{} `json:"data"`
	Errors    interface{} `json:"errors"`
	RequestId string      `json:"request_id" example:"24fa3770-628c-49de-aa17-3a338f73d99b"`
	Timestamp string      `json:"timestamp" example:"2022-04-27 23:19:56"`
}

type UnauthorizedResponse struct {
	Code      string      `json:"code" example:"KDMU-02-012"`
	Message   string      `json:"message" example:"token tidak valid."`
	Data      interface{} `json:"data"`
	Errors    interface{} `json:"errors"`
	RequestId string      `json:"request_id" example:"24fa3770-628c-49de-aa17-3a338f73d99b"`
	Timestamp string      `json:"timestamp" example:"2022-04-27 23:19:56"`
}

type BadRequestErrorValidationResponse struct {
	Code      string      `json:"code" example:"KDMU-02-006"`
	Message   string      `json:"message" example:"permintaan tidak valid, kesalahan muncul ketika permintaan Anda memiliki parameter yang tidak valid."`
	Data      interface{} `json:"data"`
	Errors    interface{} `json:"errors"`
	RequestId string      `json:"request_id" example:"24fa3770-628c-49de-aa17-3a338f73d99b"`
	Timestamp string      `json:"timestamp" example:"2022-04-27 23:19:56"`
}

type UnprocessableEntityResponse struct {
	Code      string      `json:"code" example:"KDMU-02-006"`
	Message   string      `json:"message" example:"permintaan tidak valid, kesalahan muncul ketika permintaan Anda memiliki parameter yang tidak valid."`
	Data      interface{} `json:"data"`
	Errors    interface{} `json:"errors"`
	RequestId string      `json:"request_id" example:"24fa3770-628c-49de-aa17-3a338f73d99b"`
	Timestamp string      `json:"timestamp" example:"2022-04-27 23:19:56"`
}

type RequestTimeoutResponse struct {
	Code      string      `json:"code" example:"KDMU-02-009"`
	Message   string      `json:"message" example:"permintaan telah melampaui batas waktu, harap request kembali."`
	Data      interface{} `json:"data"`
	Errors    interface{} `json:"errors"`
	RequestId string      `json:"request_id" example:"24fa3770-628c-49de-aa17-3a338f73d99b"`
	Timestamp string      `json:"timestamp" example:"2022-04-27 23:19:56"`
}

type InternalServerErrorResponse struct {
	Code      string      `json:"code" example:"KDMU-02-008"`
	Message   string      `json:"message" example:"terjadi kesalahan, silakan hubungi administrator."`
	Data      interface{} `json:"data"`
	Errors    interface{} `json:"errors"`
	RequestId string      `json:"request_id" example:"24fa3770-628c-49de-aa17-3a338f73d99b"`
	Timestamp string      `json:"timestamp" example:"2022-04-27 23:19:56"`
}

type ValidationErrors struct {
	Field       string `json:"field" example:"MobilePhone wajib diisi."`
	Description string `json:"message" example:"ActiveDate harus format yang benar yyyy-mm-dd."`
}