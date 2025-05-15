package time

import (
	"fmt"
	"time"
)

func FormatWithMeasurementUnit(t time.Duration) string {
	switch {
	case t.Milliseconds() != 0:
		return fmt.Sprintf("%dms", t.Milliseconds())
	case t.Microseconds() != 0:
		return fmt.Sprintf("%dmcs", t.Microseconds())
	default:
		return fmt.Sprintf("%dns", t.Nanoseconds())
	}
}
