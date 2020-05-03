package sazanalyzer

type Request struct {
	Method string
	URL    string
}

type Response struct {
	StatusCode    int
	ContentLength int
}

type Timers struct {
	ClientConnected     string
	ClientBeginRequest  string
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
	ClientDoneResponse  string
}

type Flags struct {
	ClientIP string
	HostIP   string
	Process  string
}

type Session struct {
	Number   int
	Request  Request
	Response Response
	Timers   Timers
	Timeline string
	Duration string
	Encoding string
	Caching  string
	Flags    Flags
}
