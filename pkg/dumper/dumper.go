// Package sazdumper prints a summary of SAZ files (Fiddler logs) on the console.
package sazdumper

import (
	"fmt"
	"strings"

	sazanalyzer "github.com/prantlf/saz-tools/pkg/analyzer"
	sazparser "github.com/prantlf/saz-tools/pkg/parser"
)

// Dump prints a summary line on the console for each session returned by `sazparser`.
func Dump(rawSessions []sazparser.Session) error {
	fmt.Println("Number\tTimeline\tMethod\tStatus\tURL\tBegin\tEnd\tDuration\tSize\tEncoding\tCaching\tProcess")
	fineSessions, err := sazanalyzer.Analyze(rawSessions)
	if err != nil {
		return err
	}
	lastTimeLine := fineSessions[len(fineSessions)-1].Timeline
	var durationPrecision int
	if strings.HasPrefix(lastTimeLine, "00:00") {
		durationPrecision = 6
	} else if strings.HasPrefix(lastTimeLine, "00") {
		durationPrecision = 3
	}
	for _, session := range fineSessions {
		err := printResult(&session, durationPrecision)
		if err != nil {
			return err
		}
	}
	return nil
}

func printResult(session *sazanalyzer.Session, durationPrecision int) error {
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
	duration, err := parseDuration(session.Timers.RequestResponseTime)
	if err != nil {
		return err
	}
	fmt.Printf("%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n",
		session.Number, formatDuration(timeline, durationPrecision),
		request.Method, response.StatusCode, request.URL.Full,
		formatTime(clientBeginRequest), formatTime(clientDoneResponse),
		formatDuration(duration, durationPrecision), formatSize(session.Response.ContentLength),
		session.Flags.Encoding, session.Flags.Caching, session.Flags.Process)
	return nil
}
