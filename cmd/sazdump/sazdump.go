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
	"io"
	"log"
	"net/http"
	"os"

	"github.com/thomastay/saz-tools/pkg/dumper"
	"github.com/thomastay/saz-tools/pkg/parser"
)

func main() {
	printVersion := false
	printNumber := -1
	outFile := ""
	password := ""
	headers := false
	onlyHeaders := false
	flag.BoolVar(&printVersion, "version", printVersion, "print the version of this tool and exit")
	flag.BoolVar(&printVersion, "v", printVersion, "print the version of this tool and exit (shorthand)")
	flag.IntVar(&printNumber, "n", printNumber, "Session ID to dump (dumps response body)")
	flag.StringVar(&password, "password", password, "Password")
	// flags based on https://curl.se/docs/manpage.html
	flag.StringVar(&outFile, "o", outFile, "Out File")
	flag.BoolVar(&headers, "i", headers, "Include the response headers")
	flag.BoolVar(&onlyHeaders, "I", onlyHeaders, "Include only the response headers")
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
		// os.Exit(1)
		return
	}

	var out io.Writer
	if outFile != "" {
		f, err := os.Create(outFile)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		out = f
	} else {
		out = os.Stdout
	}

	if printNumber == -1 {
		err = dumper.Dump(sessions, out)
		if err != nil {
			fmt.Printf("Printing network sessions from \"%s\" failed.\n", sazFile)
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		// User specified a specific option
		if printNumber == 0 {
			log.Println("Fiddler logs start with 1")
			os.Exit(1)
		}
		session := sessions[printNumber-1]

		if headers || onlyHeaders {
			writeHeaders(out, session.Response)
		}
		if !onlyHeaders {
			if headers {
				// Two CRLF delimiter
				out.Write([]byte("\r\n\r\n"))
			}
			body, _ := session.ResponseBody()
			out.Write(body)
		}
	}
}

func writeHeaders(out io.Writer, r *http.Response) {
	statusLine := fmt.Sprintf("%s %s\n", r.Proto, r.Status)
	out.Write([]byte(statusLine))
	// Loop over header names
	for name, values := range r.Header {
		// Loop over all values for the name.
		for _, value := range values {
			m := fmt.Sprintf("%s: %s\n", name, value)
			out.Write([]byte(m))
		}
	}
}
