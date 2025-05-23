package responses

// StandardResponse represents the standard API response structure
type StandardResponse struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Code    string      `json:"code,omitempty"`
}

// NewSuccessResponse creates a new success response
func NewSuccessResponse(data interface{}) *StandardResponse {
	return &StandardResponse{
		Status: "success",
		Data:   data,
	}
}

// NewErrorResponse creates a new error response
func NewErrorResponse(message string, code string) *StandardResponse {
	return &StandardResponse{
		Status:  "error",
		Message: message,
		Code:    code,
	}
}
