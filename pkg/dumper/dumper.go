package sazdumper

import (
	"fmt"

	sazanalyzer "github.com/prantlf/saz-tools/pkg/analyzer"
	sazparser "github.com/prantlf/saz-tools/pkg/parser"
)

func Dump(rawSessions []sazparser.Session) error {
	fmt.Println("Number\tTimeline\tMethod\tStatus\tURL\tBegin\tEnd\tDuration\tSize\tEncoding\tCaching\tProcess")
	fineSessions, err := sazanalyzer.Analyze(rawSessions)
	if err != nil {
		return err
	}
	for _, session := range fineSessions {
		err := printResult(&session)
		if err != nil {
			return err
		}
	}
	return nil
}

func printResult(session *sazanalyzer.Session) error {
	request := session.Request
	response := session.Response
	clientBeginRequest, err := parseTime(session.Timers.ClientBeginRequest)
	if err != nil {
		return err
	}
	clientDoneResponse, err := parseTime(session.Timers.ClientDoneResponse)
	if err != nil {
		return err
	}
	timeline, err := parseDuration(session.Timeline)
	if err != nil {
		return err
	}
	duration, err := parseDuration(session.Duration)
	if err != nil {
		return err
	}
	fmt.Printf("%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n",
		session.Number, formatDuration(timeline),
		request.Method, response.StatusCode, request.URL,
		formatTime(clientBeginRequest), formatTime(clientDoneResponse),
		formatDuration(duration), formatSize(session.Response.ContentLength),
		session.Encoding, session.Caching, session.Flags.Process)
	return nil
}
