package statements

import (
	"fmt"
	"strings"
)

// ErrorType represents the type of error
type ErrorType string

const (
	// ErrorTypeValidation indicates a validation error
	ErrorTypeValidation ErrorType = "validation"
	// ErrorTypeCalculation indicates a calculation error
	ErrorTypeCalculation ErrorType = "calculation"
	// ErrorTypeRendering indicates a rendering error
	ErrorTypeRendering ErrorType = "rendering"
	// ErrorTypeFormat indicates a formatting error
	ErrorTypeFormat ErrorType = "format"
)

// Error represents a detailed error with context
type Error struct {
	Type    ErrorType
	Field   string
	Message string
	Details map[string]interface{}
}

func (e Error) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("[%s] %s: %s", e.Type, e.Field, e.Message)
	}
	return fmt.Sprintf("[%s] %s", e.Type, e.Message)
}

// ValidationErrors represents multiple validation errors
type ValidationErrors struct {
	Errors []ValidationError
}

func (ve ValidationErrors) Error() string {
	if len(ve.Errors) == 0 {
		return "validation failed"
	}

	messages := make([]string, len(ve.Errors))
	for i, err := range ve.Errors {
		messages[i] = err.Error()
	}
	return "validation failed: " + strings.Join(messages, "; ")
}

// Add adds a new validation error
func (ve *ValidationErrors) Add(field, message string) {
	ve.Errors = append(ve.Errors, ValidationError{
		Field:   field,
		Message: message,
	})
}

// HasErrors returns true if there are validation errors
func (ve *ValidationErrors) HasErrors() bool {
	return len(ve.Errors) > 0
}

// NewValidationError creates a single validation error
func NewValidationError(field, message string) error {
	return ValidationError{
		Field:   field,
		Message: message,
	}
}

// NewCalculationError creates a calculation error
func NewCalculationError(message string) error {
	return Error{
		Type:    ErrorTypeCalculation,
		Message: message,
	}
}

// NewRenderingError creates a rendering error
func NewRenderingError(format, message string) error {
	return Error{
		Type:    ErrorTypeRendering,
		Field:   format,
		Message: message,
	}
}

// NewFormatError creates a formatting error
func NewFormatError(field, message string) error {
	return Error{
		Type:    ErrorTypeFormat,
		Field:   field,
		Message: message,
	}
}