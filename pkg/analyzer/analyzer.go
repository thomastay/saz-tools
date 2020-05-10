// Package analyzer computes timings and other useful metrics for sessions
// parsed from SAZ files (Fiddler logs).
package analyzer

import (
	parser "github.com/prantlf/saz-tools/pkg/parser"
)

// Analyze converts raw sessions returned by `parser` to fine sessions
// with aggregated timings and other useful metrics.
func Analyze(rawSessions []parser.Session) ([]Session, error) {
	length := len(rawSessions)
	fineSessions := make([]Session, length)
	clienBeginFirstRequest, err := ParseTime(rawSessions[0].Timers.ClientBeginRequest)
	if err != nil {
		return nil, err
	}
	for i := 0; i < length; i++ {
		err := analyzeSession(&rawSessions[i], &fineSessions[i], clienBeginFirstRequest)
		if err != nil {
			return nil, err
		}
	}
	return fineSessions, nil
}
