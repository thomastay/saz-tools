package sazanalyzer

import (
	"net/http"
	"time"

	sazparser "github.com/prantlf/saz-tools/pkg/parser"
)

func Analyze(rawSessions []sazparser.Session) ([]Session, error) {
	length := len(rawSessions)
	fineSessions := make([]Session, length)
	var clienBeginFirstRequest time.Time
	for i := 0; i < length; i++ {
		rawSession := &rawSessions[i]
		fineSession := &fineSessions[i]
		clientBeginRequest, err := parseTime(rawSession.Timers.ClientBeginRequest)
		if err != nil {
			return nil, err
		}
		clientDoneResponse, err := parseTime(rawSession.Timers.ClientDoneResponse)
		if err != nil {
			return nil, err
		}
		if clienBeginFirstRequest.IsZero() {
			clienBeginFirstRequest = clientBeginRequest
		}
		fineSession.Timeline = formatDuration(clientBeginRequest.Sub(clienBeginFirstRequest))
		fineSession.Duration = formatDuration(clientDoneResponse.Sub(clientBeginRequest))
		var encoding, caching string
		method := rawSession.Request.Method
		if method != http.MethodConnect {
			encoding = rawSession.Response.Header.Get("Content-Encoding")
			if encoding == "" {
				encoding = "raw"
			}
			caching = rawSession.Response.Header.Get("Cache-Control")
			if caching == "" {
				caching = "unpecified"
			}
		} else {
			encoding = "N/A"
			caching = "N/A"
		}
		fineSession.Encoding = encoding
		fineSession.Caching = caching
		fineSession.Number = rawSession.Number
		fineSession.Request.Method = method
		fineSession.Request.URL = rawSession.Request.URL.String()
		fineSession.Response.StatusCode = rawSession.Response.StatusCode
		fineSession.Response.ContentLength = int(rawSession.Response.ContentLength)
		fineSession.Timers.ClientConnected = rawSession.Timers.ClientConnected
		fineSession.Timers.ClientBeginRequest = rawSession.Timers.ClientBeginRequest
		fineSession.Timers.GotRequestHeaders = rawSession.Timers.GotRequestHeaders
		fineSession.Timers.ClientDoneRequest = rawSession.Timers.ClientDoneRequest
		fineSession.Timers.GatewayTime = rawSession.Timers.GatewayTime
		fineSession.Timers.DNSTime = rawSession.Timers.DNSTime
		fineSession.Timers.TCPConnectTime = rawSession.Timers.TCPConnectTime
		fineSession.Timers.HTTPSHandshakeTime = rawSession.Timers.HTTPSHandshakeTime
		fineSession.Timers.ServerConnected = rawSession.Timers.ServerConnected
		fineSession.Timers.FiddlerBeginRequest = rawSession.Timers.FiddlerBeginRequest
		fineSession.Timers.ServerGotRequest = rawSession.Timers.ServerGotRequest
		fineSession.Timers.ServerBeginResponse = rawSession.Timers.ServerBeginResponse
		fineSession.Timers.GotResponseHeaders = rawSession.Timers.GotResponseHeaders
		fineSession.Timers.ServerDoneResponse = rawSession.Timers.ServerDoneResponse
		fineSession.Timers.ClientBeginResponse = rawSession.Timers.ClientBeginResponse
		fineSession.Timers.ClientDoneResponse = rawSession.Timers.ClientDoneResponse
		fineSession.Flags.ClientIP = rawSession.Flags["x-clientip"]
		fineSession.Flags.HostIP = rawSession.Flags["x-hostip"]
		fineSession.Flags.Process = rawSession.Flags["x-processinfo"]
	}
	return fineSessions, nil
}
