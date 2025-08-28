package utils

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"
)

// Common error codes
const (
	CodeBadRequest          = "BAD_REQUEST"
	CodeUnauthorized        = "UNAUTHORIZED"
	CodeForbidden           = "FORBIDDEN"
	CodeNotFound            = "NOT_FOUND"
	CodeInternalServerError = "INTERNAL_SERVER_ERROR"
)

const (
	timeFormat = "2006-01-02T15:04:05.000Z"
)

var (
	ErrNilWriter      = errors.New("response writer cannot be nil")
	ErrEmptyMessage   = errors.New("message cannot be empty")
	ErrInvalidCode    = errors.New("response code cannot be empty")
	ErrEncodingFailed = errors.New("failed to encode response")
)

// Response represents the standard API response structure
type Response struct {
	Code     string      `json:"responseCode"` // Machine-readable status code
	Message  string      `json:"message"`      // Human-readable message
	Datetime string      `json:"datetime"`     // Timestamp in ISO 8601 format
	Data     interface{} `json:"data,omitempty"` // Optional response payload
}

// NewResponse creates a new standard API response
// Returns error if code or message is empty
func NewResponse(code, message string, data interface{}) (*Response, error) {
	if code == "" {
		return nil, ErrInvalidCode
	}
	if message == "" {
		return nil, ErrEmptyMessage
	}

	return &Response{
		Code:     code,
		Message:  message,
		Datetime: time.Now().UTC().Format(timeFormat),
		Data:     data,
	}, nil
}

// WriteJSON writes the response as JSON with the given status code
// Returns an error if writing the response fails
func (r *Response) WriteJSON(w http.ResponseWriter, statusCode int) error {
	if w == nil {
		return ErrNilWriter
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(r); err != nil {
		log.Printf("Failed to encode response: %v", err)
		return ErrEncodingFailed
	}

	return nil
}
// sendResponse is a helper function to send a response with the given parameters
func sendResponse(w http.ResponseWriter, code, message string, data interface{}, statusCode int) error {
	resp, err := NewResponse(code, message, data)
	if err != nil {
		return err
	}

	return resp.WriteJSON(w, statusCode)
}

// SendSuccess sends a successful JSON response
// Returns an error if response writing fails
func SendSuccess(w http.ResponseWriter, code, message string, statusCode int, data interface{}) error {
	return sendResponse(w, code, message, data, statusCode)
}

// SendError sends an error JSON response with the specified status code
// Returns an error if response writing fails
func SendError(w http.ResponseWriter, statusCode int, code, message string) error {
	return sendResponse(w, code, message, nil, statusCode)
}

// --- HTTP Status Code Helpers ---

// BadRequest writes a 400 Bad Request response
// Returns an error if response writing fails
func BadRequest(w http.ResponseWriter, message string) error {
	return SendError(w, http.StatusBadRequest, CodeBadRequest, message)
}

// Unauthorized writes a 401 Unauthorized response
// Returns an error if response writing fails
func Unauthorized(w http.ResponseWriter, message string) error {
	return SendError(w, http.StatusUnauthorized, CodeUnauthorized, message)
}

// Forbidden writes a 403 Forbidden response
// Returns an error if response writing fails
func Forbidden(w http.ResponseWriter, message string) error {
	return SendError(w, http.StatusForbidden, CodeForbidden, message)
}

// NotFound writes a 404 Not Found response
// Returns an error if response writing fails
func NotFound(w http.ResponseWriter, message string) error {
	return SendError(w, http.StatusNotFound, CodeNotFound, message)
}

// InternalServerError writes a 500 Internal Server Error response
// Returns an error if response writing fails
func InternalServerError(w http.ResponseWriter, message string) error {
	return SendError(w, http.StatusInternalServerError, CodeInternalServerError, message)
}
