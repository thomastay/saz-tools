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

func convertFlags(session *parser.Session) Flags {
	target := make(Flags)
	source := session.Flags.Flags
	for index := range source {
		flag := &source[index]
		target[flag.Name] = flag.Value
	}
	return target
}

func GetExtras(session *parser.Session) (RequestExtras, ResponseExtras) {
	return RequestExtras{
			Extras{
				session.Request.Header,
				session.Request.TransferEncoding,
			},
			session.Request.Host,
			session.Request.RemoteAddr,
			session.Request.PostForm,
		}, ResponseExtras{
			Extras{
				session.Response.Header,
				session.Response.TransferEncoding,
			},
			session.Response.Proto,
			session.Response.Uncompressed,
		}
}
