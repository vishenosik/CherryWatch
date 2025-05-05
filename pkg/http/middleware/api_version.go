package middleware

import (
	"context"
	"fmt"
	"net/http"

	pkghttp "github.com/vishenosik/CherryWatch/pkg/http"
)

// VersionHandler defines the interface for version handling
type VersionHandler interface {
	ParseRequest(r *http.Request) error
	WithContext(ctx context.Context) context.Context
}

// ApiVersionMiddleware validates the API version from request
func ApiVersionMiddleware(handler VersionHandler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := handler.ParseRequest(r); err != nil {
				pkghttp.SendErrors(w, http.StatusBadRequest, fmt.Sprintf("API version error: %s", err))
				return
			}

			ctx := handler.WithContext(r.Context())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
