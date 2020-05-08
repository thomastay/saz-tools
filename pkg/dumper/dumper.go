// Package dumper prints a summary of SAZ files (Fiddler logs) on the console.
package dumper

import (
	"fmt"
	"strings"

	analyzer "github.com/prantlf/saz-tools/pkg/analyzer"
	parser "github.com/prantlf/saz-tools/pkg/parser"
)

// Dump prints a summary line on the console for each session returned by `parser`.
func Dump(rawSessions []parser.Session) error {
	fmt.Println("Number\tTimeline\tMethod\tStatus\tURL\tBegin\tEnd\tDuration\tSize\tEncoding\tCaching\tProcess")
	fineSessions, err := analyzer.Analyze(rawSessions)
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
	for index := range fineSessions {
		err := printResult(&fineSessions[index], durationPrecision)
		if err != nil {
			return err
		}
	}
	return nil
}

func printResult(session *analyzer.Session, durationPrecision int) error {
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
	process := session.Flags["x-processinfo"]
	if process == "" {
		process = "unknown"
	}
	fmt.Printf("%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n",
		session.Number, formatDuration(timeline, durationPrecision),
		request.Method, response.StatusCode, request.URL.Full,
		formatTime(clientBeginRequest), formatTime(clientDoneResponse),
		formatDuration(duration, durationPrecision), formatSize(session.Response.ContentLength),
		session.Response.Encoding, session.Response.Caching, session.Request.Process)
	return nil
}
