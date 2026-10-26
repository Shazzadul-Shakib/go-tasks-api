package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// APIError is the body of every error response: {"error": {...}}.
type APIError struct {
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Details []FieldError `json:"details,omitempty"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// RespondError writes the error and stops the handler chain.
func RespondError(c *gin.Context, status int, code, message string, details ...FieldError) {
	c.AbortWithStatusJSON(status, gin.H{"error": APIError{Code: code, Message: message, Details: details}})
}

// bindError turns a ShouldBindJSON error into a 400 or 422 response.
func bindError(c *gin.Context, err error) {
	var (
		verrs     validator.ValidationErrors
		typeErr   *json.UnmarshalTypeError
		syntaxErr *json.SyntaxError
	)
	switch {
	case errors.As(err, &verrs):
		details := make([]FieldError, 0, len(verrs))
		for _, fe := range verrs {
			details = append(details, FieldError{Field: fe.Field(), Message: fieldMessage(fe)})
		}
		RespondError(c, http.StatusUnprocessableEntity, "validation_failed", "request body is invalid", details...)
	case errors.As(err, &typeErr):
		RespondError(c, http.StatusBadRequest, "invalid_json",
			fmt.Sprintf("field %q must be of type %s", typeErr.Field, typeErr.Type))
	case errors.As(err, &syntaxErr), errors.Is(err, io.ErrUnexpectedEOF):
		RespondError(c, http.StatusBadRequest, "invalid_json", "request body is not valid JSON")
	case errors.Is(err, io.EOF):
		RespondError(c, http.StatusBadRequest, "invalid_json", "request body is empty")
	default:
		RespondError(c, http.StatusBadRequest, "bad_request", err.Error())
	}
}

// fieldMessage turns a validator rule into a human sentence.
func fieldMessage(fe validator.FieldError) string {
	unit := ""
	if fe.Kind() == reflect.String {
		unit = " characters"
	}
	switch fe.Tag() {
	case "required":
		return "is required"
	case "min":
		return "must be at least " + fe.Param() + unit
	case "max":
		return "must be at most " + fe.Param() + unit
	default:
		return "failed the " + fe.Tag() + " rule"
	}
}