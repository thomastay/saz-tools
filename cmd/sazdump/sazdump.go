// Parses and prints an analysis of SAZ files (Fiddler logs) on the console.
//
//   $ sazdump foo.saz -h
//   Usage: sazdump <file.saz>
//
//   $ sazdump foo.saz
//   Number Timeline Method Code URL Begin End Duration Size Encoding Cache Process
//   1 00:00.000 GET 200 https://example.com/foo 16:11:58.755 16:11:59.013 00:00.257 506 gzip max-age=31536000 chromium:18160
//   2 00:01.573 GET 200 https://example.com/bar 16:11:59.419 16:11:59.873 00:00.454 1,201 raw no-store,no-cache chromium:18160
//   ...
package main

import (
	"flag"
	"fmt"
	"os"

	dumper "github.com/prantlf/saz-tools/pkg/dumper"
	parser "github.com/prantlf/saz-tools/pkg/parser"
)

func main() {
	flag.Usage = func() {
		fmt.Println("Usage: sazdump <file.saz>")
		flag.PrintDefaults()
	}
	flag.Parse()
	sazFile := flag.Arg(0)
	sessions, err := parser.ParseFile(sazFile)
	if err != nil {
		fmt.Printf("Parsing \"%s\" failed.\n", sazFile)
		fmt.Println(err)
		os.Exit(1)
	}
	err = dumper.Dump(sessions)
	if err != nil {
		fmt.Printf("Printing network sessions from \"%s\" failed.\n", sazFile)
		fmt.Println(err)
		os.Exit(1)
	}
}
