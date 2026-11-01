package cors

import (
	"net/http"
	"os"
	"strings"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		origin := request.Header.Get("Origin")
		if origin != "" && isOriginAllowed(origin) {
			response.Header().Set("Access-Control-Allow-Origin", origin)
			response.Header().Set("Vary", "Origin")
			response.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			response.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		}

		if request.Method == http.MethodOptions {
			if origin == "" || !isOriginAllowed(origin) {
				response.WriteHeader(http.StatusForbidden)
				return
			}

			response.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(response, request)
	})
}

func isOriginAllowed(origin string) bool {
	allowedOrigins := strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ",")
	for _, allowedOrigin := range allowedOrigins {
		if strings.TrimSpace(allowedOrigin) == origin {
			return true
		}
	}

	return false
}
