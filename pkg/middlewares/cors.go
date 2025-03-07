package middlewares

import (
	"net/http"
	"strings"
)

type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

func NewCORSConfig(origins, methods, headers []string) *CORSConfig {
	return &CORSConfig{
		AllowedOrigins: origins,
		AllowedMethods: methods,
		AllowedHeaders: headers,
	}
}

func (c *CORSConfig) CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if contains(c.AllowedHeaders, "*") || contains(c.AllowedOrigins, origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", joinSlice(c.AllowedMethods, "GET, POST, PUT, DELETE, OPTIONS"))
			w.Header().Set("Accescs-Control-Allow-Headers", joinSlice(c.AllowedHeaders, "Content-Type, Authorization"))

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}

	return false
}

func joinSlice(slice []string, defaultValue string) string {
	if len(slice) == 0 {
		return defaultValue
	}

	return strings.Join(slice, ", ")
}
