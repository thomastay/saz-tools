package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	analyzer "github.com/prantlf/saz-tools/pkg/analyzer"
	parser "github.com/prantlf/saz-tools/pkg/parser"
	"github.com/sourcegraph/syntaxhighlight"
)

type sazData struct {
	Key      string
	Sessions []analyzer.Session
}

func postSaz(responseWriter http.ResponseWriter, request *http.Request) interface{} {
	if request.Method != http.MethodPost {
		responseWriter.Header().Set("Allow", "POST")
		http.Error(responseWriter, "Disallowed Method", http.StatusMethodNotAllowed)
		return nil
	}
	defer request.Body.Close()
	bodyContent, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return nil
	}
	bodyReader := bytes.NewReader(bodyContent)
	rawSessions, err := parser.ParseReader(bodyReader, request.ContentLength)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return nil
	}
	sazKey := sessionCache.Put(rawSessions)
	fineSessions, err := analyzer.Analyze(rawSessions)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return nil
	}
	return sazData{sazKey, fineSessions}
}

var urlPath = regexp.MustCompile("([^/]+)")

func getSaz(responseWriter http.ResponseWriter, request *http.Request) interface{} {
	if request.Method != http.MethodGet && request.Method != http.MethodHead {
		responseWriter.Header().Set("Allow", "HEAD,GET")
		http.Error(responseWriter, "Disallowed Method", http.StatusMethodNotAllowed)
		return nil
	}
	pathSegments := urlPath.FindAllString(request.URL.Path, -1)
	segmentCount := len(pathSegments)
	if segmentCount < 3 {
		http.Error(responseWriter, "Missing Key", http.StatusBadRequest)
		return nil
	}
	sazKey := pathSegments[2] // /api/saz/:key
	rawSessions, ok := sessionCache.Get(sazKey)
	if !ok {
		http.Error(responseWriter, "Unknown Key", http.StatusNotFound)
		return nil
	}
	if segmentCount == 3 {
		if request.Method == http.MethodHead {
			responseWriter.WriteHeader(204)
			return nil
		}
		fineSessions, err := analyzer.Analyze(rawSessions)
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
			return nil
		}
		return fineSessions
	}
	sessionNumber, err := strconv.Atoi(pathSegments[3]) // /api/saz/:key/:number
	if err != nil {
		http.Error(responseWriter, "Invalid Key", http.StatusBadRequest)
		return nil
	}
	if sessionNumber <= 0 || sessionNumber > len(rawSessions) {
		http.Error(responseWriter, "Invalid Key", http.StatusBadRequest)
		return nil
	}
	if request.Method == http.MethodHead {
		responseWriter.WriteHeader(204)
		return nil
	}
	session := &rawSessions[sessionNumber-1]
	if segmentCount == 4 {
		if request.URL.Query().Get("scope") == "extras" {
			return analyzer.GetExtras(session)
		} else {
			clienBeginFirstRequest, err := analyzer.ParseTime(rawSessions[0].Timers.ClientBeginRequest)
			if err != nil {
				http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
				return nil
			}
			response, err := analyzer.MergeExtras(session, clienBeginFirstRequest)
			if err != nil {
				http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
				return nil
			}
			return response
		}
	}
	if segmentCount == 6 && pathSegments[5] == "body" {
		if pathSegments[4] == "request" { // /api/saz/:key/:number/request/body
			contentType := session.Request.Header.Get("Content-Type")
			if contentType == "application/x-www-form-urlencoded" {
				contentType = "text/plain"
			}
			responseWriter.Header().Set("Content-Type", contentType)
			if _, err = responseWriter.Write(session.RequestBody); err != nil {
				http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
			}
			return nil
		}
		if pathSegments[4] == "response" { // /api/saz/:key/:number/response/body
			contentType := session.Response.Header.Get("Content-Type")
			content := session.ResponseBody
			if strings.HasPrefix(contentType, "text/html") ||
				strings.HasPrefix(contentType, "application/xhtml+xml") {
				highlighted, err := syntaxhighlight.AsHTML(content)
				if err != nil {
					contentType = "text/plain"
				} else {
					// pre{background:#1d1f21;font-family:Menlo,Bitstream Vera Sans Mono,DejaVu Sans Mono,Monaco,Consolas,monospace;border:0!important}.pln{color:#c5c8c6}ol.linenums{margin-top:0;margin-bottom:0;color:#969896}li.L0,li.L1,li.L2,li.L3,li.L4,li.L5,li.L6,li.L7,li.L8,li.L9{padding-left:1em;background-color:#1d1f21;list-style-type:decimal}@media screen{.str{color:#b5bd68}.kwd{color:#b294bb}.com{color:#969896}.typ{color:#81a2be}.lit{color:#de935f}.pun{color:#c5c8c6}.opn{color:#c5c8c6}.clo{color:#c5c8c6}.tag{color:#c66}.atn{color:#de935f}.atv{color:#8abeb7}.dec{color:#de935f}.var{color:#c66}.fun{color:#81a2be}}
					prolog := `<style>
	pre{background:#fff;font-family:Menlo,Bitstream Vera Sans Mono,DejaVu Sans Mono,Monaco,Consolas,monospace;border:0!important}.pln{color:#333}ol.linenums{margin-top:0;margin-bottom:0;color:#ccc}li.L0,li.L1,li.L2,li.L3,li.L4,li.L5,li.L6,li.L7,li.L8,li.L9{padding-left:1em;background-color:#fff;list-style-type:decimal}@media screen{.str{color:#183691}.kwd{color:#a71d5d}.com{color:#969896}.typ{color:#0086b3}.lit{color:#0086b3}.pun{color:#333}.opn{color:#333}.clo{color:#333}.tag{color:navy}.atn{color:#795da3}.atv{color:#183691}.dec{color:#333}.var{color:teal}.fun{color:#900}}
</style>
<pre>`
					if _, err = responseWriter.Write([]byte(prolog)); err != nil {
						http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
					}
					content = highlighted
				}
			}
			responseWriter.Header().Set("Content-Type", contentType)
			if _, err = responseWriter.Write(content); err != nil {
				http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
			}
			return nil
		}
	}
	http.Error(responseWriter, "Invalid Path", http.StatusNotFound)
	return nil
}

type api struct{}

func (h *api) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	var payload interface{}
	switch {
	case request.URL.Path == "/api/saz":
		payload = postSaz(responseWriter, request)
	case strings.HasPrefix(request.URL.Path, "/api/saz/"):
		payload = getSaz(responseWriter, request)
	default:
		http.Error(responseWriter, "Unrecognized Path", http.StatusNotFound)
	}
	if payload != nil {
		sendJSON(responseWriter, payload)
	}
}

func sendJSON(responseWriter http.ResponseWriter, payload interface{}) {
	output, err := json.Marshal(payload)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	responseWriter.Header().Set("Content-Type", "application/json")
	io.WriteString(responseWriter, string(output))
}
