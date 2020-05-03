package sazanalyzer

import (
	"fmt"
	"time"
)

func parseTime(dateTime string) (time.Time, error) {
	return time.Parse(time.RFC3339Nano, dateTime)
}

func formatDuration(duration time.Duration) string {
	duration = duration.Round(time.Microsecond)
	hours := duration / time.Hour
	duration -= hours * time.Hour
	minutes := duration / time.Minute
	duration -= minutes * time.Minute
	seconds := duration / time.Second
	duration -= seconds * time.Second
	microseconds := duration / time.Microsecond
	return fmt.Sprintf("%02d:%02d:%02d.%06d", hours, minutes, seconds, microseconds)
}
