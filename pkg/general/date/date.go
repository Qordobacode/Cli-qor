package date

import (
	"time"
)

var (
	form = "2006-01-02 15:04:05"
)

// GetDateFromTimestamp convert timestamp value into human-readable string
func GetDateFromTimestamp(timestamp int64) string {
	tm := time.Unix(timestamp/1000, 0)
	return tm.Format(form)
}
