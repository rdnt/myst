package responses

// Error asdsd
type Error struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// AddError adds a new validation error to the current context
func AddError(name string, desc string) Error {
	return Error{
		Name:        name,
		Description: desc,
	}
}

// ResponseMessage is a JSON success/validation failure message (HTTP code 200)
type ResponseMessage struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

// ErrorMessage is an JSON error message
type ErrorMessage struct {
	Status  string      `json:"status"`
	Code    string      `json:"code"`
	Message interface{} `json:"message"`
}

// NewResponse returns a success response message
func NewResponse(data interface{}) *ResponseMessage {
	// Format response message
	return &ResponseMessage{
		Status: "success",
		Data:   data,
	}
}

// NewFail returns a validation failure response message
func NewFail(data interface{}) *ResponseMessage {
	// Format response message
	return &ResponseMessage{
		Status: "fail",
		Data:   data,
	}
}

// NewError returns an error message
func NewError(code string, msg interface{}) *ErrorMessage {
	// Format response message
	return &ErrorMessage{
		Status:  "error",
		Code:    code,
		Message: msg,
	}
}

// NewHTTPError returns an HTTP error message
func NewHTTPError(id int) *ErrorMessage {
	// Format response message
	var code, msg string
	switch id {
	case 400:
		code = "BAD REQUEST"
		msg = "Bad request"
	case 401:
		code = "UNAUTHORIZED"
		msg = "Unauthorized"
	case 403:
		code = "FORBIDDEN"
		msg = "Forbidden"
	case 404:
		code = "NOT FOUND"
		msg = "Not found"
	case 429:
		code = "TOO MANY REQUESTS"
		msg = "You are being rate limited"
	case 500:
		code = "INTERNAL SERVER ERROR"
		msg = "Internal Server Error"
	case 503:
		code = "SERVICE UNAVAILABLE"
		msg = "Service Unavailable"
	default:
		code = "INTERNAL SERVER ERROR"
		msg = "Internal Server Error"
	}
	return &ErrorMessage{
		Status:  "error",
		Code:    code,
		Message: msg,
	}
}
