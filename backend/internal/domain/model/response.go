package model

type SuccessResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

type ErrorResponse struct {
	Success bool       `json:"success"`
	Error   ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    string        `json:"code"`
	Message string        `json:"message"`
	Details []FieldError  `json:"details,omitempty"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type Pagination struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

func NewSuccessResponse(data interface{}) SuccessResponse {
	return SuccessResponse{Success: true, Data: data}
}

func NewListResponse(data interface{}, pagination Pagination) SuccessResponse {
	return SuccessResponse{Success: true, Data: data, Pagination: &pagination}
}

func NewErrorResponse(code, message string) ErrorResponse {
	return ErrorResponse{
		Success: false,
		Error:   ErrorDetail{Code: code, Message: message},
	}
}

func NewValidationErrorResponse(code, message string, details []FieldError) ErrorResponse {
	return ErrorResponse{
		Success: false,
		Error:   ErrorDetail{Code: code, Message: message, Details: details},
	}
}
