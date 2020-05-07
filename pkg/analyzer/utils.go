package analyzer

import (
	"fmt"
	"time"

	parser "github.com/prantlf/saz-tools/pkg/parser"
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

func getFlag(session *parser.Session, name string) (string, bool) {
	flags := session.Flags.Flags
	for index := range flags {
		flag := &flags[index]
		if flag.Name == name {
			return flag.Value, true
		}
	}
	return "", false
}
