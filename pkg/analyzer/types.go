package analyzer

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
	Method string
	URL    URL
}

// Response contains information about a server response.
type Response struct {
	StatusCode    int
	ContentLength int
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
type Flags struct {
	Encoding string
	Caching  string
	ClientIP string
	HostIP   string
	Process  string
}

// Session represents an analyzed network session.
type Session struct {
	Number   int
	Timeline string
	Request  Request
	Response Response
	Timers   Timers
	Flags    Flags
}
