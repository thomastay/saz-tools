package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	cache "github.com/prantlf/saz-tools/internal/cache"
	sazanalyzer "github.com/prantlf/saz-tools/pkg/analyzer"
	sazparser "github.com/prantlf/saz-tools/pkg/parser"
)

type PostPayload struct {
	Key      string
	Sessions []sazanalyzer.Session
}

func postSaz(responseWriter http.ResponseWriter, request *http.Request) interface{} {
	if request.Method != http.MethodPost {
		responseWriter.Header().Set("Allow", "POST")
		http.Error(responseWriter, "Only POST allowed.", http.StatusMethodNotAllowed)
		return nil
	}
	if err := request.ParseMultipartForm(128 << 20); err != nil {
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return nil
	}
	fileReader, fileHeader, err := request.FormFile("saz")
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return nil
	}
	defer fileReader.Close()
	rawSessions, err := sazparser.ParseReader(fileReader, fileHeader.Size)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return nil
	}
	fineSessions, err := sazanalyzer.Analyze(rawSessions)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return nil
	}
	key := entryCache.Put(cache.Entry{rawSessions, fineSessions})
	return PostPayload{key, fineSessions}
}

func getSaz(responseWriter http.ResponseWriter, request *http.Request) interface{} {
	if request.Method != http.MethodGet {
		responseWriter.Header().Set("Allow", "GET")
		http.Error(responseWriter, "Only GET allowed.", http.StatusMethodNotAllowed)
		return nil
	}
	key := request.URL.Path[9:] // /api/saz/
	entry, ok := entryCache.Get(key)
	if !ok {
		http.Error(responseWriter, "Kay Not Found", http.StatusNotFound)
		return nil
	}
	return entry.FineSessions
}

type api struct{}

func (h *api) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	var payload interface{}
	path := request.URL.Path
	switch {
	case path == "/api/saz":
		payload = postSaz(responseWriter, request)
	case strings.HasPrefix(path, "/api/saz/"):
		payload = getSaz(responseWriter, request)
	default:
		http.Error(responseWriter, "Path Not Found", http.StatusNotFound)
	}
	if payload == nil {
		return
	}
	output, err := json.Marshal(payload)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	responseWriter.Header().Set("Content-Type", "application/json")
	io.WriteString(responseWriter, string(output))
}
