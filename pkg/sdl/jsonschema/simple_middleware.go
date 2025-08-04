package jsonschema

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"go.uber.org/zap"

	"go-transport-prac/internal/errors"
	"go-transport-prac/internal/logger"
	"go-transport-prac/internal/types"
)

// SimpleHTTPMiddleware provides simple HTTP middleware for JSON Schema validation
type SimpleHTTPMiddleware struct {
	validator *XeipuuvValidator
	logger    *logger.Logger
}

// NewSimpleHTTPMiddleware creates a new simple HTTP middleware
func NewSimpleHTTPMiddleware(validator *XeipuuvValidator, logger *logger.Logger) *SimpleHTTPMiddleware {
	return &SimpleHTTPMiddleware{
		validator: validator,
		logger:    logger,
	}
}

// ValidateRequest validates HTTP request body against a schema
func (m *SimpleHTTPMiddleware) ValidateRequest(schemaID string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip validation for non-JSON content types
			contentType := r.Header.Get("Content-Type")
			if !strings.Contains(strings.ToLower(contentType), "application/json") {
				next.ServeHTTP(w, r)
				return
			}

			// Read request body
			body, err := io.ReadAll(r.Body)
			if err != nil {
				m.writeErrorResponse(w, http.StatusBadRequest,
					errors.BadRequestError(errors.CodeInvalidInput,
						"failed to read request body"))
				return
			}

			// Restore request body for next handler
			r.Body = io.NopCloser(bytes.NewBuffer(body))

			// Skip validation if body is empty
			if len(body) == 0 {
				next.ServeHTTP(w, r)
				return
			}

			// Validate against schema
			if err := m.validator.ValidateJSON(schemaID, string(body)); err != nil {
				if m.logger != nil {
					m.logger.Warn("Request validation failed",
						zap.String("schema_id", schemaID),
						zap.String("method", r.Method),
						zap.String("path", r.URL.Path),
						zap.Error(err),
					)
				}

				m.writeErrorResponse(w, http.StatusBadRequest, err.(*errors.AppError))
				return
			}

			if m.logger != nil {
				m.logger.Debug("Request validation successful",
					zap.String("schema_id", schemaID),
					zap.String("method", r.Method),
					zap.String("path", r.URL.Path),
				)
			}

			next.ServeHTTP(w, r)
		})
	}
}

// ValidationHandler creates a standalone validation handler
func (m *SimpleHTTPMiddleware) ValidationHandler(schemaID string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			m.writeErrorResponse(w, http.StatusMethodNotAllowed,
				errors.BadRequestError(errors.CodeInvalidInput, "only POST method is allowed"))
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			m.writeErrorResponse(w, http.StatusBadRequest,
				errors.BadRequestError(errors.CodeInvalidInput, "failed to read request body"))
			return
		}

		err = m.validator.ValidateJSON(schemaID, string(body))

		response := SimpleValidationResponse{
			Valid:    err == nil,
			SchemaID: schemaID,
		}

		if err != nil {
			response.Error = err.Error()
		}

		statusCode := http.StatusOK
		if !response.Valid {
			statusCode = http.StatusBadRequest
		}

		m.writeJSONResponse(w, statusCode, response)
	}
}

// SimpleValidationResponse represents the response from validation endpoint
type SimpleValidationResponse struct {
	Valid    bool   `json:"valid"`
	SchemaID string `json:"schema_id"`
	Error    string `json:"error,omitempty"`
	Message  string `json:"message,omitempty"`
}

// Helper methods

func (m *SimpleHTTPMiddleware) writeErrorResponse(w http.ResponseWriter, statusCode int, err *errors.AppError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorResponse := types.APIResponse[interface{}]{
		Success: false,
		Error: &types.APIError{
			Code:    err.Code,
			Message: err.Message,
			Details: err.Details,
			Fields:  err.Fields,
		},
	}

	json.NewEncoder(w).Encode(errorResponse)
}

func (m *SimpleHTTPMiddleware) writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}