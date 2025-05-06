package http

import (
	"context"
	"fmt"
	"net/http"

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

// VersionHandler defines the interface for version handling
type VersionHandler interface {
	ParseRequest(r *http.Request) error
	WithContext(ctx context.Context) context.Context
}

// apiVersionKey is an unexported type for context keys to prevent collisions
type apiVersionKey struct{}

// APIVersion represents a parsed semantic version
type APIVersion struct {
	Version models.Version
	param   string
	header  string
	minimal models.Version
	current models.Version
}

type ApiVersionOption = func(*APIVersion)

func defaultApiVersion() *APIVersion {
	return &APIVersion{
		param:  VersionParam,
		header: VersionHeader,
	}
}

func newApiVersion(version models.Version, current string, opts ...ApiVersionOption) (*APIVersion, error) {
	av := defaultApiVersion()
	av.Version = version
	av.minimal = version
	currentVer, err := version.Parse_(current)
	if err != nil {
		return nil, err
	}

	av.current = currentVer
	for _, opt := range opts {
		opt(av)
	}
	return av, nil
}

func mustInitApiVersion(version models.Version, current string, opts ...ApiVersionOption) *APIVersion {
	av, err := newApiVersion(version, current, opts...)
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

	version, err := av.Version.Parse_(versionStr)
	if err != nil {
		return errors.Wrap(err, "invalid version format")
	}

	if !version.In_(av.minimal, av.current) {
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
func ApiVersionFromContext(ctx context.Context) (models.Version, error) {
	val := ctx.Value(apiVersionKey{})
	if val == nil {
		return nil, errors.New("no API version in context")
	}

	version, ok := val.(models.Version)
	if !ok {
		return nil, errors.New("invalid API version type in context")
	}
	return version, nil
}

// ApiVersionMiddleware validates the API version from request
func ApiVersionMiddleware(
	version models.Version,
	current string,
	opts ...ApiVersionOption,
) func(http.Handler) http.Handler {

	av := mustInitApiVersion(version, current, opts...)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := av.ParseRequest(r); err != nil {
				SendErrors(w, http.StatusBadRequest, fmt.Sprintf("API version error: %s", err))
				return
			}
			ctx := context.WithValue(r.Context(), apiVersionKey{}, av.Version)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
