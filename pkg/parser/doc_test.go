package sazparser_test

import (
	"fmt"
	"io"
	"os"

	sazparser "github.com/prantlf/saz-tools/pkg/parser"
)

// Parse the content of `foo.saz` and print the count of network sessions.
func ExampleParseFile() {
	sessions, err := sazparser.ParseFile("foo.saz")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("%d sessions found.", len(sessions))
	// Output: 42 sessions found.
}

// Parse the content of `foo.saz` and print the count of network sessions.
func ExampleParseReader() {
	var reader io.ReaderAt
	var size int64
	sessions, err := sazparser.ParseReader(reader, size)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var total int64 = 0
	for index := range sessions {
		total += sessions[index].Response.ContentLength
	}
	fmt.Printf("The total downloaded size was %d bytes.", total)
	// Output: The total downloaded size was 44040192 bytes.
}
