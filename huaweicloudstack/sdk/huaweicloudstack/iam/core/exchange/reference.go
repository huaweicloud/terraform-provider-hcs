package exchange

import "time"

type ApiReference struct {
	Host          string
	Method        string
	Path          string
	Raw           string
	UserAgent     string
	StartedTime   time.Time
	DurationMs    time.Duration
	RequestId     string
	StatusCode    int
	ContentLength int64
}
