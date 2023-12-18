package response

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

// Response response types
type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

// constants of statuses
const (
	StatusOK    = "OK"
	StatusError = "Error"
)

// OK ok response status
func OK() Response {
	return Response{
		Status: StatusOK,
	}
}

// Error error response status
func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}

func ValidationError(errs validator.ValidationErrors) Response {
	var errMessages []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMessages = append(errMessages, fmt.Sprintf("field %s is a required field", err.Field()))
		case "url":
			errMessages = append(errMessages, fmt.Sprintf("field %s is not a valid URL", err.Field()))
		default:
			errMessages = append(errMessages, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}

	return Response{
		Status: StatusError,
		Error:  strings.Join(errMessages, ", "),
	}
}
