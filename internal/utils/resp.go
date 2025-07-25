package utils

import (
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	StatusCode int         `json:"status_code,omitempty"`
	Status     string      `json:"status,omitempty"`
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Error      interface{} `json:"error,omitempty"`
}

func ErrorResponse(statusCode int, message string, error interface{}) Response {
	if statusCode >= 500 {
		message = ""
		error = "server error"
	}
	statText := http.StatusText(statusCode)
	return Response{
		StatusCode: statusCode,
		Status:     statText,
		Message:    message,
		Error:      error,
	}
}

func SuccessfulResponse(statusCode int, message string, data interface{}) Response {
	return Response{
		StatusCode: statusCode,
		Status:     http.StatusText(statusCode),
		Message:    message,
		Data:       data,
	}
}

func ValidationErrorResponse(verror validator.ValidationErrors) Response {
	fieldError := make([]map[string]interface{}, len(verror))
	for ind, err := range verror {
		fieldError[ind] = map[string]interface{}{
			"field": err.Field(),
			"error": validationErrorMessage(err),
		}
	}
	return ErrorResponse(http.StatusUnprocessableEntity, "validation error", fieldError)
}

func validationErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "this field is required"
	case "email":
		return "invalid email address"
	case "min":
		return "value is too short"
	case "max":
		return "value is too large"
	}
	return "invalid value"
}
