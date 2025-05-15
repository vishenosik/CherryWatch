// The colors package provides text highlighting functionality through the Higlighter struct,
// which can highlight numbers and keywords in strings using configurable color codes.
// The package uses regular expressions for pattern matching and supports a functional options pattern
// for flexible configuration.
package colors

import (
	// builtin
	"regexp"
	"strings"
)

var (
	// NumberRegex matches standalone numbers (both integers and floats) while ignoring:
	// 	- Numbers embedded in words (e.g., "abc123")
	// 	- Numbers in UUIDs (e.g., "d72be13c-9a0a-4df9-8475-d0f4b0701248")
	// Pattern breakdown:
	//   (^|\s)    - Start of string or whitespace
	//   (\d+      - One or more digits
	//   (?:\.\d+)? - Optional decimal part
	//   )         - End of number capture
	//   (\s|$|[^\w-]) - Followed by whitespace, end of string, or non-word character (except hyphen)
	NumberRegex = regexp.MustCompile(`(^|\s)(\d+(?:\.\d+)?)(\s|$|[^\w-])`)
)

// Higlighter provides text highlighting capabilities for numbers and keywords.
type Higlighter struct {
	// Compiled regex for number matching
	numbers *regexp.Regexp
	// Flag to enable/disable number highlighting
	doNumbers bool
	// Color code for number highlighting
	numbersColor ColorCode
	// Compiled regex for keyword matching
	keywords *regexp.Regexp
	// Map of keywords to their corresponding color codes
	keywordsToColors map[string]ColorCode
}

// optsFunc is a function type used for configuring a Higlighter instance.
// It follows the functional options pattern to provide flexible configuration.
//
// Parameters:
//
//	h - Pointer to the Higlighter instance to be configured
type optsFunc func(*Higlighter)

// NewHighlighter creates and returns a new Higlighter instance with optional configurations.
//
// Parameters:
//
//	opts - Variadic list of option functions to configure the highlighter
//
// Returns:
//
//	*Higlighter - Newly created and configured Higlighter instance
//
// Example:
//
//	h := NewHighlighter(
//	    WithNumbersHighlight(ColorRed),
//	    WithKeyWordsHighlight(map[string]ColorCode{"error": ColorRed}),
//	)
func NewHighlighter(opts ...optsFunc) *Higlighter {
	h := &Higlighter{
		numbers: NumberRegex,
	}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

// Modify updates an existing Higlighter or creates a new one if the input is nil.
// It applies the provided option functions to configure the Higlighter.
//
// Parameters:
//
//	h    - Pointer to the Higlighter to modify (nil creates a new instance)
//	opts - Variadic list of option functions to apply
//
// Returns:
//
//	*Higlighter - Modified or newly created Higlighter instance
func Modify(h *Higlighter, opts ...optsFunc) *Higlighter {
	if h == nil {
		h = NewHighlighter(opts...)
	}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

// HighlightNumbers applies color highlighting to numbers in the input string if enabled.
//
// Parameters:
//
//	src - Input string to process
//
// Returns:
//
//	string - Processed string with numbers highlighted if enabled, otherwise original string
func (h *Higlighter) HighlightNumbers(src string) string {
	if h.doNumbers {
		return h.numbers.ReplaceAllStringFunc(src, func(s string) string {
			return colors[h.numbersColor](s)
		})
	}
	return src
}

// HighlightKeyWords applies color highlighting to specified keywords in the input string.
//
// Parameters:
//
//	src - Input string to process
//
// Returns:
//
//	string - Processed string with keywords highlighted if any are defined, otherwise original string
func (h *Higlighter) HighlightKeyWords(src string) string {
	if len(h.keywordsToColors) == 0 {
		return src
	}
	return h.keywords.ReplaceAllStringFunc(src, func(s string) string {
		return colors[h.keywordsToColors[s]](s)
	})
}

// WithNumbersHighlight returns an option function that enables number highlighting.
//
// Parameters:
//
//	color - ColorCode to use for highlighting numbers
func WithNumbersHighlight(color ColorCode) optsFunc {
	return func(h *Higlighter) {
		h.doNumbers = true
		h.numbersColor = color
	}
}

// WithKeyWordsHighlight returns an option function that enables keyword highlighting.
//
// Parameters:
//
//	keywordsToColors - Map of keywords to their ColorCodes
func WithKeyWordsHighlight(keywordsToColors map[string]ColorCode) optsFunc {
	return func(h *Higlighter) {
		keywords := make([]string, 0, len(keywordsToColors))
		for key := range keywordsToColors {
			keywords = append(keywords, key)
		}
		h.keywords = compileKeyWordsRegex(keywords...)
		h.keywordsToColors = keywordsToColors
	}
}

// compileKeyWordsRegex compiles a regular expression for matching keywords.
// Each keyword is wrapped with word boundaries to ensure whole-word matching.
//
// Parameters:
//
//	keywords - Variadic list of keywords to match
//
// Returns:
//
//	*regexp.Regexp - Compiled regular expression for keyword matching
func compileKeyWordsRegex(keywords ...string) *regexp.Regexp {
	for i := range keywords {
		keywords[i] = `\b` + regexp.QuoteMeta(keywords[i]) + `\b`
	}
	return regexp.MustCompile(strings.Join(keywords, "|"))
}
