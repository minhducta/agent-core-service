package domain

import "errors"

// Error code constants
const (
	ErrCodeValidation   = "VALIDATION_ERROR"
	ErrCodeNotFound     = "NOT_FOUND"
	ErrCodeConflict     = "CONFLICT"
	ErrCodeUnauthorized = "UNAUTHORIZED"
	ErrCodeForbidden    = "FORBIDDEN"
	ErrCodeInternal     = "INTERNAL_ERROR"
)

// Sentinel errors
var (
	ErrBotNotFound          = errors.New("bot not found")
	ErrBotInactive          = errors.New("bot is inactive or banned")
	ErrInvalidAPIKey        = errors.New("invalid or missing API key")
	ErrMemoryNotFound       = errors.New("memory not found")
	ErrSkillNotFound        = errors.New("skill not found")
	ErrPolicyNotFound       = errors.New("policy not found")
	ErrTodoNotFound         = errors.New("todo not found")
	ErrChecklistItemNotFound = errors.New("checklist item not found")
	ErrHeartbeatNotFound    = errors.New("no heartbeat recorded for this bot")
)

// ErrorResponse represents an HTTP error response body
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail contains error code and message
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
