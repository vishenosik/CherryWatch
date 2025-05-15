package log

import (
	// builtin
	"log/slog"
	"time"

	// internal
	time_helper "github.com/vishenosik/CherryWatch/pkg/time"
)

const (
	AttrError        = "err"
	AttrOperation    = "operation"
	AttrTook         = "took"
	AttrUserID       = "user_id" // Assuming User struct has field "ID"
	AttrAppID        = "app_id"  // Assuming App struct has field
	AttrAppComponent = "app_component"
)

func Error(err error) slog.Attr {
	return slog.String(AttrError, err.Error())
}

func Operation(op string) slog.Attr {
	return slog.String(AttrOperation, op)
}

func Took(timeStart time.Time) slog.Attr {
	return slog.String(AttrTook, time_helper.FormatWithMeasurementUnit(time.Since(timeStart)))
}

func UserID(userID string) slog.Attr {
	return slog.String(AttrUserID, userID)
}

func AppID(appID string) slog.Attr {
	return slog.String(AttrAppID, appID)
}

func AppComponent(component string) slog.Attr {
	return slog.String(AttrAppComponent, component)
}
