package analyzer_test

import (
	"fmt"

	"github.com/thomastay/saz-tools/pkg/analyzer"
	"github.com/thomastay/saz-tools/pkg/parser"
)

// Analyze the content of `foo.saz` and print the duration of the network
// session with the biggest response.
func ExampleAnalyze() {
	rawSessions, err := parser.ParseFile("foo.saz", "")
	if err != nil {
		panic(err)
	}
	fineSessions, err := analyzer.Analyze(rawSessions)
	if err != nil {
		panic(err)
	}
	var biggest *analyzer.Session
	for index := range fineSessions {
		session := &fineSessions[index]
		if biggest == nil || session.Response.ContentLength > biggest.Response.ContentLength {
			biggest = session
		}
	}
	fmt.Printf("The biggest response was obtained in %s.", biggest.Timers.RequestResponseTime)
	// Output: The biggest response was obtained in 00:01:42:042001.
}
