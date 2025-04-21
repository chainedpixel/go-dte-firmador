package response

// Response represents the standard API response structure
type Response struct {
	Status string      `json:"status"`
	Body   interface{} `json:"body"`
}

// ErrorBody represents the structure of an error response body
type ErrorBody struct {
	Code    string      `json:"error_code"`
	Message interface{} `json:"message"`
}

// NewSuccessResponse creates a success response with the provided data
func NewSuccessResponse(data interface{}) *Response {
	return &Response{
		Status: "OK",
		Body:   data,
	}
}

// NewErrorResponse creates an error response with the provided code and message
func NewErrorResponse(code string, message interface{}) *Response {
	return &Response{
		Status: "error",
		Body: ErrorBody{
			Code:    code,
			Message: message,
		},
	}
}
