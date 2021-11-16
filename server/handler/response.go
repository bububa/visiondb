package handler

const (
	// SuccessCode code for success response
	SuccessCode = 0
	// SuccessMessage message for success response
	SuccessMessage = "success"
	// ErrorCode code for general error response
	ErrorCode = -1
)

// Response represents general response
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// SuccessResponse success response wrapper
func SuccessResponse(data interface{}) Response {
	return Response{
		Code:    SuccessCode,
		Message: "success",
		Data:    data,
	}
}

// ErrorResponse error response wrapper
func ErrorResponse(code int, msg string) Response {
	return Response{
		Code:    code,
		Message: msg,
	}
}
