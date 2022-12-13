// Package dumper prints a summary of SAZ files (Fiddler logs) on the console.
package dumper

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/thomastay/saz-tools/internal/pluralizer"
	"github.com/thomastay/saz-tools/pkg/analyzer"
	"github.com/thomastay/saz-tools/pkg/parser"
)

// Dump prints a summary line on the console for each session returned by `parser`.
func Dump(rawSessions []parser.Session, writer io.Writer) error {
	fmt.Println("Number\tTimeline\tMethod\tStatus\tURL\tBegin\tEnd\tDuration\tSize\tEncoding\tCaching\tProcess")
	fineSessions, err := analyzer.Analyze(rawSessions)
	if err != nil {
		message := fmt.Sprintf("Analyzing %d network sessions failed.", len(rawSessions))
		return fmt.Errorf("%s\n%s", message, err.Error())
	}
	lastTimeLine := fineSessions[len(fineSessions)-1].Timeline
	var durationPrecision int
	if strings.HasPrefix(lastTimeLine, "00:00") {
		durationPrecision = 6
	} else if strings.HasPrefix(lastTimeLine, "00") {
		durationPrecision = 3
	}
	for index := range fineSessions {
		err := printResult(&fineSessions[index], durationPrecision, writer)
		if err != nil {
			message := fmt.Sprintf("Printing %s network session failed.",
				pluralizer.FormatOrdinal(index+1))
			return fmt.Errorf("%s\n%s", message, err.Error())
		}
	}
	return nil
}

func printResult(session *analyzer.Session, durationPrecision int, writer io.Writer) error {
	request := session.Request
	response := session.Response
	method := request.Method
	clientBegin, err := analyzer.ParseTime(session.Timers.ClientBegin)
	if err != nil {
		message := fmt.Sprintf("Parsing ClientBegin time from \"%s\" failed.",
			session.Timers.ClientBegin)
		return fmt.Errorf("%s\n%s", message, err.Error())
	}
	clientDoneResponse, err := analyzer.ParseTime(session.Timers.ClientDoneResponse)
	if err != nil {
		message := fmt.Sprintf("Parsing ClientDoneResponse time from \"%s\" failed.",
			session.Timers.ClientDoneResponse)
		return fmt.Errorf("%s\n%s", message, err.Error())
	}
	timeline, err := parseDuration(session.Timeline)
	if err != nil {
		message := fmt.Sprintf("Parsing Timeline duration from \"%s\" failed.",
			session.Timeline)
		log.Printf("%s\n%s", message, err.Error())
	}
	duration, err := parseDuration(session.Timers.RequestResponseTime)
	if err != nil {
		message := fmt.Sprintf("Parsing RequestResponseTime duration from \"%s\" failed.",
			session.Timers.RequestResponseTime)
		log.Printf("%s\n%s", message, err.Error())
	}
	message := fmt.Sprintf("%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n",
		session.Number, formatDuration(timeline, durationPrecision),
		method, response.StatusCode, request.URL.Full,
		formatTime(clientBegin), formatTime(clientDoneResponse),
		formatDuration(duration, durationPrecision), session.Response.ContentType,
		formatSize(session.Response.ContentLength), session.Response.Encoding,
		session.Response.Caching, session.Request.Process)
	writer.Write([]byte(message))
	return nil
}
