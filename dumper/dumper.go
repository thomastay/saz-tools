package sazdumper

import (
	"fmt"
	"time"

	sazparser "github.com/prantlf/sazdump/parser"
)

func Dump(sessions []sazparser.Session) error {
	fmt.Println("Number\tTimeline\tMethod\tCode\tURL\tBegin\tEnd\tDuration\tSize\tEncoding\tCache\tProcess")
	var startTime time.Time
	for _, session := range sessions {
		err := printResult(&startTime, &session)
		if err != nil {
			return err
		}
	}
	return nil
}

func printResult(startTime *time.Time, session *sazparser.Session) error {
	request := session.Response.Request
	response := session.Response
	clientBeginRequest, err := parseTime(session.Timers.ClientBeginRequest)
	if err != nil {
		return err
	}
	clientDoneResponse, err := parseTime(session.Timers.ClientDoneResponse)
	if err != nil {
		return err
	}
	if startTime.IsZero() {
		*startTime = clientBeginRequest
	}
	startOffset := clientBeginRequest.Sub(*startTime)
	duration := clientDoneResponse.Sub(clientBeginRequest)
	compression := response.Header.Get("Content-Encoding")
	if compression == "" {
		compression = "raw"
	}
	caching := response.Header.Get("Cache-Control")
	if caching == "" {
		compression = "N/A"
	}
	fmt.Printf("%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n",
		session.Number, formatDuration(startOffset),
		request.Method, response.StatusCode, request.URL.String(),
		formatTime(clientBeginRequest), formatTime(clientDoneResponse),
		formatDuration(duration), formatSize(int(response.ContentLength)),
		compression, caching, session.Flags["x-processinfo"])
	return nil
}
