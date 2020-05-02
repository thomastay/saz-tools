package sazparser

import (
	"encoding/xml"
	"net/http"
)

type Session struct {
	XMLName  xml.Name        `xml:"Session"`
	Timers   SessionTimers   `xml:"SessionTimers"`
	RawFlags RawSessionFlags `xml:"SessionFlags"`
	Flags    map[string]string
	Response *http.Response
	Number   int
}

type SessionTimers struct {
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

type RawSessionFlags struct {
	XMLName xml.Name         `xml:"SessionFlags"`
	Flags   []RawSessionFlag `xml:"SessionFlag"`
}

type RawSessionFlag struct {
	XMLName xml.Name `xml:"SessionFlag"`
	Name    string   `xml:"N,attr"`
	Value   string   `xml:"V,attr"`
}
