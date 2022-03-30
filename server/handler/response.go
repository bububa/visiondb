package handler

const (
	// SuccessCode code for success response
	SuccessCode = 0
	// SuccessMessage message for success response
	SuccessMessage = "success"
	// ErrorCode code for general error response
	ErrorCode = -1
)

type ResultType string

const (
	SUCCESS ResultType = "success"
	ERROR   ResultType = "error"
	WARNING ResultType = "warning"
)

// Response represents general response
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Type    ResultType  `json:"type,omitempty"`
	Result  interface{} `json:"result,omitempty"`
}

// SuccessResponse success response wrapper
func SuccessResponse(data interface{}) Response {
	return Response{
		Code:   SuccessCode,
		Result: data,
		Type:   SUCCESS,
	}
}

// ErrorResponse error response wrapper
func ErrorResponse(code int, msg string) Response {
	return Response{
		Code:    code,
		Message: msg,
		Type:    ERROR,
	}
}
