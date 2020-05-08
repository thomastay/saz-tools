package main

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	cache "github.com/prantlf/saz-tools/internal/cache"
	analyzer "github.com/prantlf/saz-tools/pkg/analyzer"
	parser "github.com/prantlf/saz-tools/pkg/parser"
)

type sazPayload struct {
	Key      string
	Sessions []analyzer.Session
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
	rawSessions, err := parser.ParseReader(fileReader, fileHeader.Size)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return nil
	}
	fineSessions, err := analyzer.Analyze(rawSessions)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return nil
	}
	key := entryCache.Put(cache.Entry{rawSessions, fineSessions})
	return sazPayload{key, fineSessions}
}

var urlPath = regexp.MustCompile("([^/]+)")

type sessionPayload struct {
	Request  analyzer.RequestExtras
	Response analyzer.ResponseExtras
}

func getSaz(responseWriter http.ResponseWriter, request *http.Request) interface{} {
	if request.Method != http.MethodGet {
		responseWriter.Header().Set("Allow", "GET")
		http.Error(responseWriter, "Only GET allowed.", http.StatusMethodNotAllowed)
		return nil
	}
	pathSegments := urlPath.FindAllString(request.URL.Path, -1)
	segmentCount := len(pathSegments)
	if segmentCount < 3 {
		http.Error(responseWriter, "Missing Key", http.StatusBadRequest)
		return nil
	}
	key := pathSegments[2] // /api/saz/:key
	entry, ok := entryCache.Get(key)
	if !ok {
		http.Error(responseWriter, "Unknown Kay", http.StatusNotFound)
		return nil
	}
	switch segmentCount {
	case 3:
		return entry.FineSessions
	case 4:
		number, err := strconv.Atoi(pathSegments[3]) // /api/saz/:key/:number
		if err != nil {
			http.Error(responseWriter, "Invalid Key", http.StatusBadRequest)
			return nil
		}
		sessions := entry.RawSessions
		if number <= 0 || number > len(sessions) {
			http.Error(responseWriter, "Invalid Key", http.StatusBadRequest)
			return nil
		}
		requestExtras, responseExtras := analyzer.GetExtras(&sessions[number-1])
		return &sessionPayload{requestExtras, responseExtras}
	}
	http.Error(responseWriter, "Invalid Path", http.StatusNotFound)
	return nil
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
		http.Error(responseWriter, "Unrecognized Path", http.StatusNotFound)
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
