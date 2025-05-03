package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/vishenosik/CherryWatch/pkg/models"
)

const (
	VersionParam  = "v"
	VersionHeader = "X-API-Version"
)

var (
	ErrFormat = errors.New("version must be in format MAJOR.MINOR")
)

type Version struct {
	Major int
	Minor int
}

// String returns the version as "MAJOR.MINOR" string
func (v *Version) String() string {
	return fmt.Sprintf("%d.%d", v.Major, v.Minor)
}

// parseVersion parses a semantic version string (e.g., "2.1") into APIVersion
func parseVersion(versionStr string) (*Version, error) {
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

	return &Version{Major: major, Minor: minor}, nil
}

// Greater finds out if v1 >= v2
func (v1 *Version) GTE(v2 *Version) bool {
	if v1.Major != v2.Major {
		return v1.Major >= v2.Major
	}
	if v1.Minor != v2.Minor {
		return v1.Minor >= v2.Major
	}
	return true
}

// Greater finds out if v1 > v2
func (v1 *Version) GT(v2 *Version) bool {
	if v1.Major != v2.Major {
		return v1.Major > v2.Major
	}
	if v1.Minor != v2.Minor {
		return v1.Minor > v2.Major
	}
	return false
}

// apiVersionKey is an unexported type for context keys to prevent collisions
type apiVersionKey struct{}

// APIVersion represents a parsed semantic version
type APIVersion struct {
	*Version
	param   string
	header  string
	minimal *Version
	current *Version
}

type ApiVersionOption = func(*APIVersion)

func defaultApiVersion() *APIVersion {
	return &APIVersion{
		param:   VersionParam,
		header:  VersionHeader,
		minimal: &Version{},
	}
}

func NewApiVersion(current string, opts ...ApiVersionOption) (*APIVersion, error) {
	av := defaultApiVersion()
	currentVer, err := parseVersion(current)
	if err != nil {
		return nil, err
	}
	av.current = currentVer
	for _, opt := range opts {
		opt(av)
	}
	return av, nil
}

func MustInitApiVersion(current string, opts ...ApiVersionOption) *APIVersion {
	av, err := NewApiVersion(current, opts...)
	if err != nil {
		panic(err)
	}
	return av
}

// ParseFirst parses the version from request (query param or header)
func (av *APIVersion) ParseRequest(r *http.Request) error {

	versionStr := r.URL.Query().Get(VersionParam)
	if versionStr == "" {
		versionStr = r.Header.Get(VersionHeader)
	}

	version, err := parseVersion(versionStr)
	if err != nil {
		return errors.Wrap(err, "invalid version format")
	}

	if av.minimal.GTE(version) || version.GT(av.current) {
		return fmt.Errorf("unsupported version. Min: %s, Max: %s", av.minimal, av.current)
	}

	av.Version = version
	return nil
}

// WithContext adds the APIVersion to the context
func (av *APIVersion) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, apiVersionKey{}, av.Version)
}

// ApiVersionFromContext retrieves the APIVersion from context
func ApiVersionFromContext(ctx context.Context) (*Version, error) {
	val := ctx.Value(apiVersionKey{})
	if val == nil {
		return nil, errors.New("no API version in context")
	}

	version, ok := val.(*Version)
	if !ok {
		return nil, errors.New("invalid API version type in context")
	}
	return version, nil
}

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
				models.SendErrors(w, http.StatusBadRequest, fmt.Sprintf("API version error: %s", err))
				return
			}

			ctx := handler.WithContext(r.Context())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
