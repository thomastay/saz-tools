package analyzer

import (
	"net/http"
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
	ContentLength int
	Process       string
}

// Response contains information about a server response.
type Response struct {
	StatusCode    int
	ContentLength int
	Encoding      string
	Caching       string
}

// Timers contain begin and end times of phases of a network session including
// computed durations of those phases.
type Timers struct {
	ClientConnected     string
	ClientBeginRequest  string
	GotRequestHeaders   string
	ClientDoneRequest   string
	GatewayTime         string
	DNSTime             string
	TCPConnectTime      string
	HTTPSHandshakeTime  string
	RequestResponseTime string
	RequestSendTime     string
	ServerProcessTime   string
	ResponseReceiveTime string
	ServerConnected     string
	FiddlerBeginRequest string
	ServerGotRequest    string
	ServerBeginResponse string
	GotResponseHeaders  string
	ServerDoneResponse  string
	ClientBeginResponse string
	ClientDoneResponse  string
}

// Flags contain properties of a network session, which are not included
// in request or response headers.
type Flags map[string]string

// Session represents an analyzed network session.
type Session struct {
	Number   int
	Timeline string
	Request  Request
	Response Response
	Timers   Timers
	Flags    Flags
}

// Extras is a base structure for additional information for requests and responses.
type Extras struct {
	Header           http.Header
	TransferEncoding []string
}

// Request contains additional information about a client request.
type RequestExtras struct {
	Extras
	Host          string
	RemoteAddress string
	Fields        url.Values
}

// Response contains additional information about a server response.
type ResponseExtras struct {
	Extras
	Protocol     string
	Uncompressed bool
}
