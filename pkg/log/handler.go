package log

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/vishenosik/CherryWatch/pkg/colors"
	"gopkg.in/yaml.v2"
)

const (
	timeFormat = "15:04:05.000"
)

type customAttrs struct {
	component string
}

type Handler struct {
	handler slog.Handler
	writer  io.Writer
	rec     nextFunc
	buf     *bytes.Buffer
	mutex   *sync.Mutex

	// syntax highlighter
	highlight *colors.Higlighter

	// marshaller type
	marshalType uint8

	attrs customAttrs
}

// The signature of the function for setting parameters
type HandlerOption func(*Handler)

func defaultHandler() *Handler {
	handlerOptions := &slog.HandlerOptions{Level: slog.LevelDebug}
	buf := &bytes.Buffer{}

	h := &Handler{
		handler: slog.NewJSONHandler(buf, &slog.HandlerOptions{
			Level:       handlerOptions.Level,
			AddSource:   handlerOptions.AddSource,
			ReplaceAttr: suppressDefaultAttrs(handlerOptions.ReplaceAttr),
		}),
		buf:       buf,
		writer:    os.Stdout,
		highlight: colors.NewHighlighter(),
		rec:       handlerOptions.ReplaceAttr,
		mutex:     &sync.Mutex{},
	}

	return h
}

func NewHandler(opts ...HandlerOption) *Handler {
	h := defaultHandler()
	for _, opt := range opts {
		opt(h)
	}
	return h
}

func (h *Handler) Handle(ctx context.Context, rec slog.Record) error {

	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("[%s] ", rec.Time.Format(timeFormat)))

	if h.attrs.component != "" {
		builder.WriteString(fmt.Sprintf("[%s] ", color.GreenString(h.attrs.component)))
	}

	builder.WriteString(fmt.Sprintf("%s: %s\n", level(rec), color.CyanString(rec.Message)))

	attrs, err := h.computeAttrs(ctx, rec)
	if err != nil {
		return err
	}

	attrsStr, err := h.marshal(attrs)
	if err != nil {
		return err
	}

	attrsStr = h.highlight.HighlightNumbers(attrsStr)
	attrsStr = h.highlight.HighlightKeyWords(attrsStr)

	if attrsStr != "" {
		builder.WriteString(fmt.Sprintf("%s\n", attrsStr))
	}

	_, err = io.WriteString(h.writer, builder.String())
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	out := copy(h)
	out.handler = h.handler.WithAttrs(attrs)

	for _, attr := range attrs {
		switch attr.Key {
		case AttrAppComponent:
			out.attrs.component = attr.Value.String()
		}
	}

	return out
}

func (h *Handler) WithGroup(name string) slog.Handler {
	out := copy(h)
	out.handler = h.handler.WithGroup(name)
	return out
}

func WithWriter(writer io.Writer) HandlerOption {
	return func(h *Handler) {
		if writer != nil {
			h.writer = writer
		}
	}
}

func WithLevel(level slog.Level) HandlerOption {
	return func(h *Handler) {
		handlerOptions := &slog.HandlerOptions{
			Level: level,
		}
		h.handler = slog.NewJSONHandler(h.buf, &slog.HandlerOptions{
			Level:       handlerOptions.Level,
			AddSource:   handlerOptions.AddSource,
			ReplaceAttr: suppressDefaultAttrs(handlerOptions.ReplaceAttr),
		})
	}
}

func WithNumbersHighlight(color colors.ColorCode) HandlerOption {
	return func(h *Handler) {
		h.highlight = colors.Modify(h.highlight, colors.WithNumbersHighlight(color))
	}
}

func WithKeyWordsHighlight(keywordsToColors map[string]colors.ColorCode) HandlerOption {
	return func(h *Handler) {
		h.highlight = colors.Modify(h.highlight, colors.WithKeyWordsHighlight(keywordsToColors))
	}
}

func WithYamlMarshaller() HandlerOption {
	return func(h *Handler) {
		h.marshalType = yaml_marshaller
	}
}

const (
	json_marshaller uint8 = iota
	yaml_marshaller
)

type nextFunc func([]string, slog.Attr) slog.Attr

type attrs = map[string]any

func suppressDefaultAttrs(next nextFunc) nextFunc {
	return func(groups []string, a slog.Attr) slog.Attr {
		switch a.Key {
		case slog.TimeKey, slog.LevelKey, slog.MessageKey, AttrAppComponent:
			return slog.Attr{}
		}
		if next == nil {
			return a
		}
		return next(groups, a)
	}
}

func (h *Handler) computeAttrs(ctx context.Context, rec slog.Record) (attrs, error) {

	h.mutex.Lock()
	defer func() {
		h.buf.Reset()
		h.mutex.Unlock()
	}()

	if err := h.handler.Handle(ctx, rec); err != nil {
		return nil, fmt.Errorf("error when calling inner handler's Handle: %w", err)
	}

	var attrs attrs
	err := json.Unmarshal(h.buf.Bytes(), &attrs)
	if err != nil {
		return nil, fmt.Errorf("error when unmarshaling inner handler's Handle result: %w", err)
	}

	return attrs, nil
}

func (h *Handler) marshal(attrs attrs) (string, error) {
	var (
		data []byte
		err  error
	)
	if len(attrs) > 0 {

		switch h.marshalType {
		case json_marshaller:
			data, err = json.MarshalIndent(attrs, "", "  ")

		case yaml_marshaller:
			data, err = yaml.Marshal(attrs)
		}

		if err != nil {
			return "", err
		}
	}
	return string(data), nil
}

func copy(h *Handler) *Handler {
	return &Handler{
		handler:     h.handler,
		writer:      h.writer,
		highlight:   h.highlight,
		buf:         h.buf,
		rec:         h.rec,
		mutex:       h.mutex,
		marshalType: h.marshalType,
		attrs:       h.attrs,
	}
}

func level(rec slog.Record) string {
	level := rec.Level.String()

	switch rec.Level {

	case slog.LevelDebug:
		level = color.MagentaString(level)

	case slog.LevelInfo:
		level = color.BlueString(level)

	case slog.LevelWarn:
		level = color.YellowString(level)

	case slog.LevelError:
		level = color.RedString(level)

	}
	return level
}
