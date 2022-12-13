// Parses and prints an analysis of SAZ files (Fiddler logs) on the console.
//
//	$ sazdump -h
//	Usage: sazdump <file.saz>
//
//	$ sazdump foo.saz
//	Number Timeline Method Code URL Begin End Duration Size Encoding Cache Process
//	1 00:00.000 GET 200 https://example.com/foo 16:11:58.755 16:11:59.013 00:00.257 506 gzip max-age=31536000 chromium:18160
//	2 00:01.573 GET 200 https://example.com/bar 16:11:59.419 16:11:59.873 00:00.454 1,201 raw no-store,no-cache chromium:18160
//	...
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/prantlf/saz-tools/pkg/dumper"
	"github.com/prantlf/saz-tools/pkg/parser"
)

func main() {
	printVersion := false
	printNumber := -1
	outFile := ""
	password := ""
	flag.BoolVar(&printVersion, "version", printVersion, "print the version of this tool and exit")
	flag.BoolVar(&printVersion, "v", printVersion, "print the version of this tool and exit (shorthand)")
	flag.IntVar(&printNumber, "n", printNumber, "Session ID to dump (dumps response body)")
	flag.StringVar(&outFile, "o", outFile, "Out File")
	flag.StringVar(&password, "password", password, "Password")
	flag.Usage = func() {
		fmt.Println("Usage: sazdump <file.saz>")
		flag.PrintDefaults()
	}
	flag.Parse()
	if printVersion {
		fmt.Printf("v%s\n", version)
		os.Exit(0)
	}
	sazFile := flag.Arg(0)
	if sazFile == "" {
		fmt.Println("sazdump: missing .saz file name")
		flag.Usage()
		os.Exit(1)
	}
	sessions, err := parser.ParseFile(sazFile, password)
	if err != nil {
		fmt.Printf("Parsing \"%s\" failed.\n", sazFile)
		fmt.Println(err)
		os.Exit(1)
	}
	if printNumber == -1 {
		err = dumper.Dump(sessions)
		if err != nil {
			fmt.Printf("Printing network sessions from \"%s\" failed.\n", sazFile)
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		session := sessions[printNumber]
		if outFile != "" {
			os.WriteFile(outFile, session.ResponseBody, 0644)
			log.Println("Wrote to", outFile)
		} else {
			fmt.Println(session.ResponseBody)

		}
	}
}
