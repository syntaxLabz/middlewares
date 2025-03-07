package middlewares

import (
	"net/http"

	"github.com/syntaxLabz/errors/pkg/httperrors"
)

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				err := httperrors.NewServerError()
				statusCode, errResp := err.ErrorResponse()
				w.WriteHeader(statusCode)
				w.Write(errResp.ToJSON())

				return
			}
		}()

		next.ServeHTTP(w, r)
	})
}
