package analyzer

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prantlf/saz-tools/internal/pluralizer"
	"github.com/prantlf/saz-tools/pkg/parser"
)

// ParseTime parses a network session timer in the maximum precision.
func ParseTime(dateTime string) (time.Time, error) {
	if dateTime == "0001-01-01T00:00:00" {
		return time.Time{}, nil
	}
	return time.Parse(time.RFC3339Nano, dateTime)
}

// GetExtras extracts additional information from  raw sessions, which was
// not returned by `Analyze`.
func GetExtras(session *parser.Session) *SessionExtras {
	return &SessionExtras{
		RequestExtras{
			Extras{
				convertHeader(&session.Request.Header),
				session.Request.TransferEncoding,
			},
			session.Request.Host,
			session.Request.RemoteAddr,
			session.Request.PostForm,
		},
		ResponseExtras{
			Extras{
				convertHeader(&session.Response.Header),
				session.Response.TransferEncoding,
			},
			session.Response.Proto,
			session.Response.Uncompressed,
		},
		TimerExtras{
			session.Timers.ClientConnected,
			session.Timers.ClientBeginRequest,
			session.Timers.GotRequestHeaders,
			session.Timers.ClientDoneRequest,
			session.Timers.GatewayTime,
			session.Timers.DNSTime,
			session.Timers.TCPConnectTime,
			session.Timers.HTTPSHandshakeTime,
			session.Timers.ServerConnected,
			session.Timers.FiddlerBeginRequest,
			session.Timers.ServerGotRequest,
			session.Timers.ServerBeginResponse,
			session.Timers.GotResponseHeaders,
			session.Timers.ServerDoneResponse,
			session.Timers.ClientBeginResponse,
		},
		convertFlags(session),
	}
}

// MergeExtras merges both basic information returned by `Analyze`
// and additional information returned by `GetExtras`.
func MergeExtras(rawSession *parser.Session, clientBeginSessions time.Time) (*MergedSession, error) {
	var basics Session
	err := analyzeSession(rawSession, &basics, clientBeginSessions)
	if err != nil {
		message := fmt.Sprintf("Analyzing %s network session failed.",
			pluralizer.FormatOrdinal(rawSession.Number))
		return nil, fmt.Errorf("%s\n%s", message, err.Error())
	}
	extras := GetExtras(rawSession)
	return &MergedSession{
		basics.Number,
		basics.Timeline,
		MergedRequest{
			basics.Request,
			extras.Request,
		},
		MergedResponse{
			basics.Response,
			extras.Response,
		},
		MergedTimers{
			basics.Timers,
			extras.Timers,
		},
		extras.Flags,
	}, nil
}

func analyzeSession(rawSession *parser.Session, fineSession *Session, clientBeginSessions time.Time) error {
	fineSession.Number = rawSession.Number
	fillSessionRequest(rawSession, fineSession)
	fillSessionResponse(rawSession, fineSession)
	err := fillSessionTimers(rawSession, fineSession, clientBeginSessions)
	if err != nil {
		return err
	}
	return nil
}

func fillSessionRequest(rawSession *parser.Session, fineSession *Session) {
	url := rawSession.Request.URL
	fineSession.Request.Method = rawSession.Request.Method
	fineSession.Request.URL.Full = url.String()
	fineSession.Request.URL.Scheme = url.Scheme
	fineSession.Request.URL.Host = url.Hostname()
	fineSession.Request.URL.HostAndPort = url.Host
	fineSession.Request.URL.Port = url.Port()
	fineSession.Request.URL.Path = url.Path
	fineSession.Request.URL.Query = url.RawQuery
	fineSession.Request.URL.PathAndQuery = url.RequestURI()
	fineSession.Request.ContentType = rawSession.Request.Header.Get("Content-Type")
	fineSession.Request.ContentLength = int(rawSession.Request.ContentLength)
	process, ok := getRawFlag(rawSession, "x-processinfo")
	if !ok {
		process = "unknown"
	}
	fineSession.Request.Process = process
}

func fillSessionResponse(rawSession *parser.Session, fineSession *Session) {
	var encoding, caching string
	if rawSession.Request.Method != http.MethodConnect {
		encoding = rawSession.Response.Header.Get("Content-Encoding")
		if encoding == "" {
			if rawSession.Response.Uncompressed {
				encoding = "raw"
			} else {
				encoding = "unspecified"
			}
		}
		caching = rawSession.Response.Header.Get("Cache-Control")
		if caching == "" {
			caching = "unspecified"
		}
	} else {
		encoding = "N/A"
		caching = "N/A"
	}
	fineSession.Response.StatusCode = rawSession.Response.StatusCode
	fineSession.Response.ContentType = rawSession.Response.Header.Get("Content-Type")
	fineSession.Response.ContentLength = int(rawSession.Response.ContentLength)
	fineSession.Response.Encoding = encoding
	fineSession.Response.Caching = caching
}

func fillSessionTimers(rawSession *parser.Session, fineSession *Session, clientBeginSessions time.Time) error {
	var clientBeginTimer string
	if rawSession.Request.Method == "CONNECT" {
		clientBeginTimer = rawSession.Timers.ClientConnected
	} else {
		clientBeginTimer = rawSession.Timers.ClientBeginRequest
	}
	clientBegin, err := ParseTime(clientBeginTimer)
	if err != nil {
		message := fmt.Sprintf("Parsing ClientBegin time from \"%s\" failed.", clientBeginTimer)
		return fmt.Errorf("%s\n%s", message, err.Error())
	}
	serverGotRequest, err := ParseTime(rawSession.Timers.ServerGotRequest)
	if err != nil {
		message := fmt.Sprintf("Parsing ServerGotRequest time from \"%s\" failed.",
			rawSession.Timers.ServerGotRequest)
		return fmt.Errorf("%s\n%s", message, err.Error())
	}
	serverBeginResponse, err := ParseTime(rawSession.Timers.ServerBeginResponse)
	if err != nil {
		message := fmt.Sprintf("Parsing ServerBeginResponse time from \"%s\" failed.",
			rawSession.Timers.ServerBeginResponse)
		return fmt.Errorf("%s\n%s", message, err.Error())
	}
	clientDoneResponse, err := ParseTime(rawSession.Timers.ClientDoneResponse)
	if err != nil {
		message := fmt.Sprintf("Parsing ClientDoneResponse time from \"%s\" failed.",
			rawSession.Timers.ClientDoneResponse)
		return fmt.Errorf("%s\n%s", message, err.Error())
	}
	fineSession.Timeline = formatDuration(computeDuration(clientBeginSessions, clientBegin))
	fineSession.Timers.ClientBegin = clientBeginTimer
	fineSession.Timers.RequestResponseTime = formatDuration(computeDuration(clientBegin, clientDoneResponse))
	fineSession.Timers.RequestSendTime = formatDuration(computeDuration(clientBegin, serverGotRequest))
	fineSession.Timers.ServerProcessTime = formatDuration(computeDuration(serverGotRequest, serverBeginResponse))
	fineSession.Timers.ResponseReceiveTime = formatDuration(computeDuration(serverBeginResponse, clientDoneResponse))
	fineSession.Timers.ClientDoneResponse = rawSession.Timers.ClientDoneResponse
	return nil
}

func computeDuration(start time.Time, stop time.Time) time.Duration {
	if stop.IsZero() {
		return 0
	}
	return stop.Sub(start)
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

func getRawFlag(session *parser.Session, name string) (string, bool) {
	source := session.Flags.Flags
	for index := range source {
		flag := &source[index]
		if name == flag.Name {
			return flag.Value, true
		}
	}
	return "", false
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

func convertHeader(input *http.Header) Header {
	output := make(Header)
	for name := range *input {
		values := input.Values(name)
		switch count := len(values); {
		case count == 1:
			output[name] = values[0]
		case count > 1:
			output[name] = values
		default:
			output[name] = nil
		}
	}
	return output
}
