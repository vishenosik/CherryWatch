package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const (
	VersionParam      = "v"
	VersionHeader     = "X-API-Version"
	MinAPIVersion     = "1.0"
	CurrentAPIVersion = "2.1"
	DefaultAPIVersion = CurrentAPIVersion
)

// apiVersionKey is an unexported type for context keys to prevent collisions
type apiVersionKey struct{}

// APIVersion represents a parsed semantic version
type APIVersion struct {
	Major int
	Minor int
}

// ParseFirst parses the version from request (query param or header)
func (av *APIVersion) ParseFirst(r *http.Request) error {
	versionStr := r.URL.Query().Get(VersionParam)
	if versionStr == "" {
		versionStr = r.Header.Get(VersionHeader)
	}

	if versionStr == "" {
		defaultVer, _ := parseVersion(DefaultAPIVersion)
		*av = *defaultVer
		return nil
	}

	version, err := parseVersion(versionStr)
	if err != nil {
		return fmt.Errorf("invalid version format: %w", err)
	}

	minVer, _ := parseVersion(MinAPIVersion)
	currentVer, _ := parseVersion(CurrentAPIVersion)

	if compareVersions(version, minVer) < 0 || compareVersions(version, currentVer) > 0 {
		return fmt.Errorf("unsupported version. Min: %s, Max: %s", MinAPIVersion, CurrentAPIVersion)
	}

	*av = *version
	return nil
}

// WithContext adds the APIVersion to the context
func (av *APIVersion) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, apiVersionKey{}, av)
}

// String returns the version as "MAJOR.MINOR" string
func (av APIVersion) String() string {
	return fmt.Sprintf("%d.%d", av.Major, av.Minor)
}

// ApiVersionFromContext retrieves the APIVersion from context
func ApiVersionFromContext(ctx context.Context) (*APIVersion, error) {
	val := ctx.Value(apiVersionKey{})
	if val == nil {
		return nil, errors.New("no API version in context")
	}

	version, ok := val.(*APIVersion)
	if !ok {
		return nil, errors.New("invalid API version type in context")
	}
	return version, nil
}

// parseVersion parses a semantic version string (e.g., "2.1") into APIVersion
func parseVersion(versionStr string) (*APIVersion, error) {
	parts := strings.Split(versionStr, ".")
	if len(parts) != 2 {
		return nil, fmt.Errorf("version must be in format MAJOR.MINOR")
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid major version: %w", err)
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid minor version: %w", err)
	}

	return &APIVersion{Major: major, Minor: minor}, nil
}

// compareVersions compares two semantic versions
func compareVersions(v1, v2 *APIVersion) int {
	if v1.Major != v2.Major {
		if v1.Major < v2.Major {
			return -1
		}
		return 1
	}
	if v1.Minor != v2.Minor {
		if v1.Minor < v2.Minor {
			return -1
		}
		return 1
	}
	return 0
}

// VersionHandler defines the interface for version handling
type VersionHandler interface {
	ParseFirst(r *http.Request) error
	WithContext(ctx context.Context) context.Context
	String() string
}

// ApiVersionMiddleware validates the API version from request
func ApiVersionMiddleware(handler VersionHandler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := handler.ParseFirst(r); err != nil {
				sendError(w, fmt.Sprintf("API version error: %s", err), http.StatusBadRequest)
				return
			}
			ctx := handler.WithContext(r.Context())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ErrorResponse represents an API error
type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

// sendError sends a JSON error response
func sendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Error:   http.StatusText(statusCode),
		Details: message,
	})
}
