package parser

import (
	"encoding/xml"
	"net/http"
)

// Session represents a deserialized network session.
// Originated at https://docs.telerik.com/fiddlercore/api/fiddler.session.
type Session struct {
	XMLName      xml.Name `xml:"Session"`
	Number       int
	Timers       Timers `xml:"SessionTimers"`
	Flags        Flags  `xml:"SessionFlags"`
	Request      *http.Request
	Response     *http.Response
	RequestBody  []byte
	ResponseBody []byte
}

// Timers contain begin and end times of phases of a deserialized network session.
// Originated at https://docs.telerik.com/fiddlercore/api/fiddler.sessiontimers.
type Timers struct {
	XMLName             xml.Name `xml:"SessionTimers"`
	ClientConnected     string   `xml:"ClientConnected,attr"`
	ClientBeginRequest  string   `xml:"ClientBeginRequest,attr"`
	GotRequestHeaders   string   `xml:"GotRequestHeaders,attr"`
	ClientDoneRequest   string   `xml:"ClientDoneRequest,attr"`
	GatewayTime         string   `xml:"GatewayTime,attr"`
	DNSTime             string   `xml:"DNSTime,attr"`
	TCPConnectTime      string   `xml:"TCPConnectTime,attr"`
	HTTPSHandshakeTime  string   `xml:"HTTPSHandshakeTime,attr"`
	ServerConnected     string   `xml:"ServerConnected,attr"`
	FiddlerBeginRequest string   `xml:"FiddlerBeginRequest,attr"`
	ServerGotRequest    string   `xml:"ServerGotRequest,attr"`
	ServerBeginResponse string   `xml:"ServerBeginResponse,attr"`
	GotResponseHeaders  string   `xml:"GotResponseHeaders,attr"`
	ServerDoneResponse  string   `xml:"ServerDoneResponse,attr"`
	ClientBeginResponse string   `xml:"ClientBeginResponse,attr"`
	ClientDoneResponse  string   `xml:"ClientDoneResponse,attr"`
}

// Flags contain properties of a deserialized network session, which
// are not included in request or response headers.
// Originated at https://docs.telerik.com/fiddlercore/api/fiddler.session.
type Flags struct {
	XMLName xml.Name `xml:"SessionFlags"`
	Flags   []Flag   `xml:"SessionFlag"`
}

// Flag contains a property of a deserialized network session, which
// are not included in request or response headers.
// Originated at https://docs.telerik.com/fiddlercore/api/fiddler.session.
type Flag struct {
	XMLName xml.Name `xml:"SessionFlag"`
	Name    string   `xml:"N,attr"`
	Value   string   `xml:"V,attr"`
}
