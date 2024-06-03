package page

import "time"

type EventID string

const (
	EventThemeChanged EventID = "event.theme.changed"
)

type Event struct {
	ID EventID
}

type ServerEvent struct {
	Msg  string
	Time time.Time
}
