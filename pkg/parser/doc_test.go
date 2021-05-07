package parser_test

import (
	"fmt"
	"io"

	"github.com/prantlf/saz-tools/pkg/parser"
)

// Parse the content of `foo.saz` and print the count of network sessions.
func ExampleParseFile() {
	sessions, err := parser.ParseFile("foo.saz")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d network sessions found.", len(sessions))
	// Output: 42 network sessions found.
}

// Parse the content of `foo.saz` and print the total size of all responses.
func ExampleParseReader() {
	var reader io.ReaderAt
	var size int64
	sessions, err := parser.ParseReader(reader, size)
	if err != nil {
		panic(err)
	}
	var total int64
	for index := range sessions {
		total += sessions[index].Response.ContentLength
	}
	fmt.Printf("The total downloaded size was %d bytes.", total)
	// Output: The total downloaded size was 44040192 bytes.
}
