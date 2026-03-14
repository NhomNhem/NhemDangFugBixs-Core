package models

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

// APIError represents an API error
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	TraceID string `json:"trace_id,omitempty"`
}

// Common error codes
const (
	ErrCodeInvalidRequest    = "INVALID_REQUEST"
	ErrCodeUnauthorized      = "UNAUTHORIZED"
	ErrCodeInvalidToken      = "INVALID_TOKEN"
	ErrCodeUserBanned        = "USER_BANNED"
	ErrCodeInsufficientFunds = "INSUFFICIENT_FUNDS"
	ErrCodeInvalidLevel      = "INVALID_LEVEL"
	ErrCodeCheatingDetected  = "CHEATING_DETECTED"
	ErrCodeRateLimitExceeded = "RATE_LIMIT_EXCEEDED"
	ErrCodeValidationError   = "VALIDATION_ERROR"
	ErrCodeInternalError     = "INTERNAL_ERROR"
)
