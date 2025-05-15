package colors

import (
	// pkg
	"github.com/fatih/color"
)

type ColorCode uint8

const (
	Red ColorCode = iota
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

// colorFunc is a function type that represents a color formatting function.
// It takes a format string and optional arguments, and returns a colored string.
//
// Parameters:
//   - format: A string that contains optional verbs for formatting.
//   - a: Optional arguments to be formatted according to the format string.
//
// Returns:
//
//	A string that has been formatted and colored according to the specific color function.
type colorFunc func(format string, a ...any) string

var colors = map[ColorCode]colorFunc{
	Red:     color.RedString,
	Green:   color.GreenString,
	Yellow:  color.YellowString,
	Blue:    color.BlueString,
	Magenta: color.MagentaString,
	Cyan:    color.CyanString,
	White:   color.WhiteString,
}
