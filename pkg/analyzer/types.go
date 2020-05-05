package sazanalyzer

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

type Request struct {
	Method string
	URL    URL
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

type Flags struct {
	Encoding string
	Caching  string
	ClientIP string
	HostIP   string
	Process  string
}

type Session struct {
	Number   int
	Timeline string
	Request  Request
	Response Response
	Timers   Timers
	Flags    Flags
}
