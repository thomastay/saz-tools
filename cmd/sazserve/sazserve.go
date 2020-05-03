package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	sazanalyzer "github.com/prantlf/saz-tools/pkg/analyzer"
	sazparser "github.com/prantlf/saz-tools/pkg/parser"
)

func handler(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		responseWriter.Header().Set("Allow", "POST")
		http.Error(responseWriter, "Only POST allowed.", http.StatusMethodNotAllowed)
		return
	}
	if err := request.ParseMultipartForm(128 << 20); err != nil {
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}
	fileReader, fileHeader, err := request.FormFile("saz")
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}
	defer fileReader.Close()
	rawSessions, err := sazparser.ParseReader(fileReader, fileHeader.Size)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}
	fineSessions, err := sazanalyzer.Analyze(rawSessions)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	output, err := json.Marshal(fineSessions)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	responseWriter.Header().Set("Content-Type", "application/json")
	io.WriteString(responseWriter, string(output))
}

func main() {
	var err error
	port := 7000
	portValue := os.Getenv("PORT")
	if portValue != "" {
		if port, err = strconv.Atoi(portValue); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	flag.IntVar(&port, "port", port, "port for the web server to listen to")
	flag.Usage = func() {
		fmt.Println("Usage: sazserve [options]\nOptions:")
		flag.PrintDefaults()
	}
	flag.Parse()
	http.HandleFunc("/saz", handler)
	http.Handle("/", http.FileServer(AssetFile()))
	if err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
