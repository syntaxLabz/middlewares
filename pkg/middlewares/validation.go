package middlewares

import (
	"net/http"
	"regexp"
	"strconv"

	"github.com/google/uuid"

	"github.com/syntaxLabz/errors/pkg/httperrors"
)

type HeaderMetaData struct {
	Type             string
	MinLength        int
	MaxLength        int
	Required         bool
	CustomValidation func(value string) bool
}

type HeaderValidation struct {
	Headers map[string]HeaderMetaData
}

func NewHeaderValidation(headers map[string]HeaderMetaData) *HeaderValidation {
	return &HeaderValidation{
		Headers: headers,
	}
}

const (
	Int    = "int"
	String = "string"
	Uuid   = "uuid"
	Email  = "email"
)

func (h *HeaderValidation) HeaderValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := httperrors.HeaderValidationError()

		for key, value := range h.Headers {
			headerValue := r.Header.Get(key)
			if headerValue == "" {
				if value.Required {
					err.AddDetail(httperrors.MissingHeader(key))
					continue
				}
			}

			if value.CustomValidation != nil {
				if !value.CustomValidation(headerValue) {
					err.AddDetail(httperrors.InvalidHeader(key))
				}
			}

			if !validateType(headerValue, value.Type, value.MinLength, value.MaxLength) {
				err.AddDetail(httperrors.InvalidHeader(key))
			}
		}

		if len(err.Details) > 0 {
			statusCode, err := err.ErrorResponse()
			w.WriteHeader(statusCode)
			w.Write(err.ToJSON())

			return
		}

		next.ServeHTTP(w, r)
	})
}

func validateType(value string, headerType string, ranges ...int) bool {
	switch headerType {
	case Int:
		return isInt(value)
	case String:
		if len(ranges) > 0 {
			return isString(value, ranges[0], ranges[1])
		}

		return true
	case Uuid:
		return isUUID(value)
	case Email:
		return isEmail(value)
	default:
		return false
	}
}

func isUUID(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}

func isInt(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func isEmail(s string) bool {
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(s)
}

func isString(s string, minLen, maxLength int) bool {
	return len(s) >= minLen && len(s) <= maxLength
}
