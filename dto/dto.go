package dto

// This package contains Data Transfer Objects (DTOs) used for transferring data
type ResponsePattern struct {
	RequestID    string `json:"request_id"`
	Status       string `json:"status"`
	Data         any    `json:"data,omitempty"`
	Message      string `json:"message,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
	Code         int    `json:"code"`
	Meta         any    `json:"meta,omitempty"`
}
