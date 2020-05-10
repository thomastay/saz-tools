// Views and analyses SAZ files (Fiddler logs) in a web browser application.
//
//   $ sazserve -h
//   Usage: sazserve [options]
//   Options:
//     -port string : port for the web server to listen to (default "7000")
//
//   $ sazserve
//   $ open http://localhost:7000
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	compressor "github.com/CAFxX/gziphandler"
	cache "github.com/prantlf/saz-tools/internal/cache"
)

var sessionCache *cache.Cache

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "7000"
	}
	flag.StringVar(&port, "port", port, "port for the web server to listen to")
	flag.Usage = func() {
		fmt.Println("Usage: sazserve [options]\nOptions:")
		flag.PrintDefaults()
	}
	flag.Parse()
	sessionCache = cache.Create()
	gzipper, err := compressor.Middleware(compressor.Prefer(compressor.PreferGzip))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	brotler, err := compressor.Middleware(compressor.Prefer(compressor.PreferBrotli))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	apiHandler := gzipper(&api{})
	http.Handle("/api/saz/", apiHandler)
	http.Handle("/api/saz", apiHandler)
	http.Handle("/", brotler(http.FileServer(AssetFile())))
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
