// Views and analyzes SAZ files (Fiddler logs) in a web browser application.
//
//	$ sazserve -h
//	Usage: sazserve [options]
//	Options:
//	  -browser       : start the web browser automatically  (default false)
//	  -port <number> : port for the web server to listen on (default "7000")
//
//	$ sazserve
//	$ open http://localhost:7000/
package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/CAFxX/gziphandler"
	"github.com/skratchdot/open-golang/open"
	"github.com/thomastay/saz-tools/internal/cache"
)

var sessionCache *cache.Cache

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "7000"
	}
	browser := false
	printVersion := false
	flag.BoolVar(&browser, "browser", browser, "start the web browser automatically")
	flag.StringVar(&port, "port", port, "port for the web server to listen on")
	flag.BoolVar(&printVersion, "version", printVersion, "print the version of this tool and exit")
	flag.BoolVar(&printVersion, "v", printVersion, "print the version of this tool and exit (shorthand)")
	flag.Usage = func() {
		fmt.Println("Usage: sazserve [options]\nOptions:")
		flag.PrintDefaults()
	}
	flag.Parse()
	if printVersion {
		fmt.Printf("v%s\n", version)
		os.Exit(0)
	}
	sessionCache = cache.Create()
	gzipper, err := gziphandler.Middleware(gziphandler.Prefer(gziphandler.PreferGzip))
	if err != nil {
		fmt.Println("Initializing a gzip compressor for REST API responses failed.")
		fmt.Println(err)
		os.Exit(1)
	}
	brotler, err := gziphandler.Middleware(gziphandler.Prefer(gziphandler.PreferBrotli))
	if err != nil {
		fmt.Println("Initializing a brotli compressor for static assets failed.")
		fmt.Println(err)
		os.Exit(1)
	}
	apiHandler := gzipper(&api{})
	http.Handle("/api/saz/", apiHandler)
	http.Handle("/api/saz", apiHandler)
	http.Handle("/", brotler(http.FileServer(AssetFile())))
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Printf("Listening for TCP on localhost:%s failed.\n", port)
		fmt.Println(err)
		os.Exit(1)
	}
	origin := fmt.Sprintf("http://localhost:%s/", port)
	if browser {
		if err := open.Start(origin); err != nil {
			fmt.Printf("Starting the browser for %s failed.\n", origin)
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("Open %s in your web browser.\n", origin)
	}
	if err := http.Serve(listener, nil); err != nil {
		fmt.Printf("Serving HTTP for localhost:%s failed.\n", port)
		fmt.Println(err)
		os.Exit(1)
	}
}
