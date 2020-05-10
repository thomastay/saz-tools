package analyzer

import (
	"net/url"
)

// URL contains a URL partsed to its parts.
type URL struct {
	Full         string
	Scheme       string
	Host         string
	HostAndPort  string
	Port         string
	Path         string
	Query        string
	PathAndQuery string
}

// Request contains information about a client request.
type Request struct {
	Method        string
	URL           URL
	ContentType   string
	ContentLength int
	Process       string
}

// Response contains information about a server response.
type Response struct {
	StatusCode    int
	ContentType   string
	ContentLength int
	Encoding      string
	Caching       string
}

// Timers contain begin and end times of important phases of a network session
// including computed durations of those phases.
type Timers struct {
	ClientBeginRequest  string
	RequestResponseTime string
	RequestSendTime     string
	ServerProcessTime   string
	ResponseReceiveTime string
	ClientDoneResponse  string
}

// Session represents an analyzed network session.
type Session struct {
	Number   int
	Timeline string
	Request  Request
	Response Response
	Timers   Timers
}

// Header contains request or response headers. Values are either strings
// or arrays (if a header occurs mutliple times).
type Header map[string]interface{}

// Extras is a base structure for additional information for requests and responses.
type Extras struct {
	Header           Header
	TransferEncoding []string
}

// RequestExtras contains additional information about a client request.
type RequestExtras struct {
	Extras
	Host          string
	RemoteAddress string
	Fields        url.Values
}

// ResponseExtras contains additional information about a server response.
type ResponseExtras struct {
	Extras
	Protocol     string
	Uncompressed bool
}

// TimerExtras contains additional network communication timers.
type TimerExtras struct {
	ClientConnected     string
	GotRequestHeaders   string
	ClientDoneRequest   string
	GatewayTime         string
	DNSTime             string
	TCPConnectTime      string
	HTTPSHandshakeTime  string
	ServerConnected     string
	FiddlerBeginRequest string
	ServerGotRequest    string
	ServerBeginResponse string
	GotResponseHeaders  string
	ServerDoneResponse  string
	ClientBeginResponse string
}

// Flags contain properties of a network session, which are not included
// in request or response headers.
type Flags map[string]string

// SessionExtras contains additional information about a network session.
type SessionExtras struct {
	Request  RequestExtras
	Response ResponseExtras
	Timers   TimerExtras
	Flags    Flags
}

// MergedRequest contains information about a client request including extras.
type MergedRequest struct {
	Request
	RequestExtras
}

// MergedResponse contains information about a server response including extras.
type MergedResponse struct {
	Response
	ResponseExtras
}

// MergedTimers contain begin and end times of all phases of a network session
// including computed durations of those phases.
type MergedTimers struct {
	Timers
	TimerExtras
}

// MergedSession represents an analyzed network session including extras.
type MergedSession struct {
	Number   int
	Timeline string
	Request  MergedRequest
	Response MergedResponse
	Timers   MergedTimers
	Flags    Flags
}
