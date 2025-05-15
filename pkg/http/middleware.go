package http

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/vishenosik/CherryWatch/pkg/api"
	attrs "github.com/vishenosik/CherryWatch/pkg/log"
	"github.com/vishenosik/CherryWatch/pkg/models"
	"github.com/vishenosik/CherryWatch/pkg/versions"
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

	if versionStr == "" {
		return errors.New("api version is not provided")
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
func ApiVersionFromContext[VersionType models.Version](ctx context.Context) (VersionType, error) {
	var empty VersionType
	val := ctx.Value(apiVersionKey{})
	if val == nil {
		return empty, errors.New("no API version in context")
	}

	version, ok := val.(VersionType)
	if !ok {
		return empty, errors.New("invalid API version type in context")
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

type VersionedHandlersMap = map[string]http.HandlerFunc

func VersionedHandler(handlers VersionedHandlersMap) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiVersion, err := ApiVersionFromContext[versions.DoubleVersion](r.Context())
		if err != nil {
			SendErrors(w, http.StatusBadRequest, err.Error())
			return
		}

		handler, ok := handlers[apiVersion.String()]
		if !ok || handler == nil {
			SendErrors(w, http.StatusNotImplemented, "api version unsupported")
			return
		}

		handler(w, r)
	}
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func RequestLogger(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lrw := newLoggingResponseWriter(w)
			timeStart := time.Now()

			next.ServeHTTP(lrw, r)

			log := logger.With(
				slog.String("method", fmt.Sprintf("%s %s", r.Method, r.URL.Path)),
				slog.Int("code", lrw.statusCode),
				attrs.Took(timeStart),
			)

			switch {
			case api.IsClientError(lrw.statusCode) || api.IsServerError(lrw.statusCode):
				log.Error("request failed with error")
			case api.IsRedirect(lrw.statusCode):
				log.Warn("request redirected")
			default:
				log.Info("request accepted")
			}
		})
	}
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(statusCode int) {
	lrw.statusCode = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}

func SetHeaders() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	}
}
