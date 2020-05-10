package analyzer

import (
	"fmt"
	"net/http"
	"time"

	parser "github.com/prantlf/saz-tools/pkg/parser"
)

// ParseTime parses a network session timer in the maximum precision.
func ParseTime(dateTime string) (time.Time, error) {
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
func MergeExtras(rawSession *parser.Session, clienBeginFirstRequest time.Time) (*MergedSession, error) {
	var basics Session
	err := analyzeSession(rawSession, &basics, clienBeginFirstRequest)
	if err != nil {
		return nil, err
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

func analyzeSession(rawSession *parser.Session, fineSession *Session, clienBeginFirstRequest time.Time) error {
	method := rawSession.Request.Method
	url := rawSession.Request.URL
	clientBeginRequest, err := ParseTime(rawSession.Timers.ClientBeginRequest)
	if err != nil {
		return err
	}
	serverGotRequest, err := ParseTime(rawSession.Timers.ServerGotRequest)
	if err != nil {
		return err
	}
	serverBeginResponse, err := ParseTime(rawSession.Timers.ServerBeginResponse)
	if err != nil {
		return err
	}
	clientDoneResponse, err := ParseTime(rawSession.Timers.ClientDoneResponse)
	if err != nil {
		return err
	}
	var encoding, caching string
	if method != http.MethodConnect {
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
	fineSession.Number = rawSession.Number
	fineSession.Timeline = formatDuration(clientBeginRequest.Sub(clienBeginFirstRequest))
	fineSession.Request.Method = method
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
	fineSession.Response.StatusCode = rawSession.Response.StatusCode
	fineSession.Response.ContentType = rawSession.Response.Header.Get("Content-Type")
	fineSession.Response.ContentLength = int(rawSession.Response.ContentLength)
	fineSession.Response.Encoding = encoding
	fineSession.Response.Caching = caching
	fineSession.Timers.ClientBeginRequest = rawSession.Timers.ClientBeginRequest
	fineSession.Timers.RequestResponseTime = formatDuration(clientDoneResponse.Sub(clientBeginRequest))
	fineSession.Timers.RequestSendTime = formatDuration(serverGotRequest.Sub(clientBeginRequest))
	fineSession.Timers.ServerProcessTime = formatDuration(serverBeginResponse.Sub(serverGotRequest))
	fineSession.Timers.ResponseReceiveTime = formatDuration(clientDoneResponse.Sub(serverBeginResponse))
	fineSession.Timers.ClientDoneResponse = rawSession.Timers.ClientDoneResponse
	process, ok := getRawFlag(rawSession, "x-processinfo")
	if !ok {
		process = "unknown"
	}
	fineSession.Request.Process = process
	return nil
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
