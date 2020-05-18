// Package analyzer computes timings and other useful metrics for sessions
// parsed from SAZ files (Fiddler logs).
package analyzer

import (
	"fmt"

	"github.com/prantlf/saz-tools/internal/pluralizer"
	"github.com/prantlf/saz-tools/pkg/parser"
)

// Analyze converts raw sessions returned by `parser` to fine sessions
// with aggregated timings and other useful metrics.
func Analyze(rawSessions []parser.Session) ([]Session, error) {
	length := len(rawSessions)
	fineSessions := make([]Session, length)
	clientBeginSessions, err := ParseTime(rawSessions[0].Timers.ClientConnected)
	if err != nil {
		message := fmt.Sprintf("Parsing ClientConnected time from \"%s\" in the first network session failed.",
			rawSessions[0].Timers.ClientConnected)
		return nil, fmt.Errorf("%s\n%s", message, err.Error())
	}
	for i := 0; i < length; i++ {
		err := analyzeSession(&rawSessions[i], &fineSessions[i], clientBeginSessions)
		if err != nil {
			message := fmt.Sprintf("Analyzing %s network session failed.",
				pluralizer.FormatOrdinal(i+1))
			return nil, fmt.Errorf("%s\n%s", message, err.Error())
		}
	}
	return fineSessions, nil
}
