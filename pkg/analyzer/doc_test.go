package sazanalyzer_test

import (
	"fmt"
	"os"

	sazanalyzer "github.com/prantlf/saz-tools/pkg/analyzer"
	sazparser "github.com/prantlf/saz-tools/pkg/parser"
)

// Analyze the content of `foo.saz` and print the duration of the network
// session with the biggest response.
func ExampleAnalyze() {
	rawSessions, err := sazparser.ParseFile("foo.saz")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fineSessions, err := sazanalyzer.Analyze(rawSessions)
	if err != nil {
		fmt.Println(err)
		return
	}
	var biggest *sazanalyzer.Session
	for index := range sessions {
		if biggest == nil || sessions[index].Response.ContentLength > biggest.Response.ContentLength {
			biggest = &sessions[index]
		}
	}
	fmt.Printf("The biggest response was obtained in $s.", biggest.Timers.RequestResponseTime)
	// Output: The biggest response was obtained in 00:01:42:042001.
}
