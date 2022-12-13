package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/sourcegraph/syntaxhighlight"
	"github.com/thomastay/saz-tools/internal/pluralizer"
	"github.com/thomastay/saz-tools/pkg/analyzer"
	"github.com/thomastay/saz-tools/pkg/parser"
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
		message := fmt.Sprintf("Reading %d bytes of the request body of the type \"%s\" failed.",
			request.ContentLength, request.Header.Get("Content-Type"))
		err = fmt.Errorf("%s\n%s", message, err.Error())
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return nil
	}
	bodyReader := bytes.NewReader(bodyContent)
	rawSessions, err := parser.ParseReader(bodyReader, request.ContentLength)
	if err != nil {
		message := fmt.Sprintf("Parsing %d bytes of the SAZ file failed.", request.ContentLength)
		err = fmt.Errorf("%s\n%s", message, err.Error())
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return nil
	}
	sazKey, err := sessionCache.Put(rawSessions)
	if err != nil {
		message := fmt.Sprintf("Caching %d network sessions failed.", len(rawSessions))
		err = fmt.Errorf("%s\n%s", message, err.Error())
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return nil
	}
	fineSessions, err := analyzer.Analyze(rawSessions)
	if err != nil {
		message := fmt.Sprintf("Analyzing %d network sessions failed.", len(rawSessions))
		err = fmt.Errorf("%s\n%s", message, err.Error())
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return nil
	}
	return sazData{sazKey, fineSessions}
}

var urlPath = regexp.MustCompile("([^/]+)")

// nolint: funlen,gocyclo // dispatches all routes below /api/saz
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
		message := fmt.Sprintf("Unknown key \"%s\"", sazKey)
		http.Error(responseWriter, message, http.StatusNotFound)
		return nil
	}
	if segmentCount == 3 {
		if request.Method == http.MethodHead {
			responseWriter.WriteHeader(204)
			return nil
		}
		fineSessions, ok := getNetworkSessions(responseWriter, rawSessions, sazKey)
		if !ok {
			return nil
		}
		return fineSessions
	}
	sessionNumber, err := strconv.Atoi(pathSegments[3]) // /api/saz/:key/:number
	if err != nil {
		message := fmt.Sprintf("Malformed session number \"%s\" for %d network sessions with the key %s",
			pathSegments[3], len(rawSessions), sazKey)
		http.Error(responseWriter, message, http.StatusBadRequest)
		return nil
	}
	if sessionNumber <= 0 || sessionNumber > len(rawSessions) {
		message := fmt.Sprintf("Invalid session number %d for %d network sessions with the key %s",
			sessionNumber, len(rawSessions), sazKey)
		http.Error(responseWriter, message, http.StatusBadRequest)
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
		}
		response, ok := getNetworkSession(responseWriter, rawSessions, sazKey, sessionNumber, session)
		if !ok {
			return nil
		}
		return response
	}
	if segmentCount == 6 && pathSegments[5] == "body" {
		if pathSegments[4] == "request" { // /api/saz/:key/:number/request/body
			sendNetworkSessionRequestBody(responseWriter, sazKey, sessionNumber, session)
			return nil
		}
		if pathSegments[4] == "response" { // /api/saz/:key/:number/response/body
			sendNetworkSessionResponseBody(responseWriter, sazKey, sessionNumber, session)
			return nil
		}
	}
	message := fmt.Sprintf("Unrecognized path \"%s\"", urlPath)
	http.Error(responseWriter, message, http.StatusNotFound)
	return nil
}

func getNetworkSessions(responseWriter http.ResponseWriter, rawSessions []parser.Session, sazKey string) ([]analyzer.Session, bool) {
	fineSessions, err := analyzer.Analyze(rawSessions)
	if err != nil {
		message := fmt.Sprintf("Analyzing %d network sessions with the key %s failed.",
			len(rawSessions), sazKey)
		err = fmt.Errorf("%s\n%s", message, err.Error())
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return nil, false
	}
	return fineSessions, true
}

func getNetworkSession(responseWriter http.ResponseWriter, rawSessions []parser.Session, sazKey string, sessionNumber int, session *parser.Session) (*analyzer.MergedSession, bool) {
	clientBeginSessions, err := analyzer.ParseTime(rawSessions[0].Timers.ClientConnected)
	if err != nil {
		message := fmt.Sprintf("Parsing ClientConnected time from \"%s\" in the first network session with the key %s failed.",
			rawSessions[0].Timers.ClientConnected, sazKey)
		message = fmt.Sprintf("%s\n%s", message, err.Error())
		http.Error(responseWriter, message, http.StatusBadRequest)
		return nil, false
	}
	response, err := analyzer.MergeExtras(session, clientBeginSessions)
	if err != nil {
		message := fmt.Sprintf("Merging extra information to %s network session with the key %s failed.",
			pluralizer.FormatOrdinal(sessionNumber), sazKey)
		message = fmt.Sprintf("%s\n%s", message, err.Error())
		http.Error(responseWriter, message, http.StatusInternalServerError)
		return nil, false
	}
	return response, true
}

func sendNetworkSessionRequestBody(responseWriter http.ResponseWriter, sazKey string, sessionNumber int, session *parser.Session) {
	contentType := session.Request.Header.Get("Content-Type")
	if contentType == "application/x-www-form-urlencoded" {
		contentType = "text/plain"
	}
	responseWriter.Header().Set("Content-Type", contentType)
	if _, err := responseWriter.Write(session.RequestBody); err != nil {
		message := fmt.Sprintf("Sending %d bytes and \"%s\" type of the request body of the %s network session with the key %s failed.",
			session.Request.ContentLength, session.Request.Header.Get("Content-Type"),
			pluralizer.FormatOrdinal(sessionNumber), sazKey)
		message = fmt.Sprintf("%s\n%s", message, err.Error())
		http.Error(responseWriter, message, http.StatusInternalServerError)
	}
}

func sendNetworkSessionResponseBody(responseWriter http.ResponseWriter, sazKey string, sessionNumber int, session *parser.Session) {
	contentType := session.Response.Header.Get("Content-Type")
	content, _ := session.ResponseBody()
	if strings.HasPrefix(contentType, "text/html") ||
		strings.HasPrefix(contentType, "application/xhtml+xml") {
		highlighted, err := syntaxhighlight.AsHTML(content)
		if err != nil {
			fmt.Printf("Highlighting syntax of %d bytes and \"%s\" type of the response body of the %s network session with the key %s failed.\n",
				session.Response.ContentLength, session.Response.Header.Get("Content-Type"),
				pluralizer.FormatOrdinal(sessionNumber), sazKey)
			contentType = "text/plain"
		} else {
			// pre{background:#1d1f21;font-family:Menlo,Bitstream Vera Sans Mono,DejaVu Sans Mono,Monaco,Consolas,monospace;border:0!important}.pln{color:#c5c8c6}ol.linenums{margin-top:0;margin-bottom:0;color:#969896}li.L0,li.L1,li.L2,li.L3,li.L4,li.L5,li.L6,li.L7,li.L8,li.L9{padding-left:1em;background-color:#1d1f21;list-style-type:decimal}@media screen{.str{color:#b5bd68}.kwd{color:#b294bb}.com{color:#969896}.typ{color:#81a2be}.lit{color:#de935f}.pun{color:#c5c8c6}.opn{color:#c5c8c6}.clo{color:#c5c8c6}.tag{color:#c66}.atn{color:#de935f}.atv{color:#8abeb7}.dec{color:#de935f}.var{color:#c66}.fun{color:#81a2be}}
			prolog := `<style>
pre{background:#fff;font-family:Menlo,Bitstream Vera Sans Mono,DejaVu Sans Mono,Monaco,Consolas,monospace;border:0!important}.pln{color:#333}ol.linenums{margin-top:0;margin-bottom:0;color:#ccc}li.L0,li.L1,li.L2,li.L3,li.L4,li.L5,li.L6,li.L7,li.L8,li.L9{padding-left:1em;background-color:#fff;list-style-type:decimal}@media screen{.str{color:#183691}.kwd{color:#a71d5d}.com{color:#969896}.typ{color:#0086b3}.lit{color:#0086b3}.pun{color:#333}.opn{color:#333}.clo{color:#333}.tag{color:navy}.atn{color:#795da3}.atv{color:#183691}.dec{color:#333}.var{color:teal}.fun{color:#900}}
</style>
<pre>`
			if _, err = responseWriter.Write([]byte(prolog)); err != nil {
				message := fmt.Sprintf("Sending prolog of the response body of the %s network session with the key %s failed.",
					pluralizer.FormatOrdinal(sessionNumber), sazKey)
				message = fmt.Sprintf("%s\n%s", message, err.Error())
				http.Error(responseWriter, message, http.StatusInternalServerError)
			}
			content = highlighted
		}
	}
	responseWriter.Header().Set("Content-Type", contentType)
	if _, err := responseWriter.Write(content); err != nil {
		message := fmt.Sprintf("Sending %d bytes and \"%s\" type of the response body of the %s network session with the key %s failed.",
			session.Response.ContentLength, session.Response.Header.Get("Content-Type"),
			pluralizer.FormatOrdinal(sessionNumber), sazKey)
		message = fmt.Sprintf("%s\n%s", message, err.Error())
		http.Error(responseWriter, message, http.StatusInternalServerError)
	}
}

type api struct{}

func (h *api) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	var payload interface{}
	urlPath := request.URL.Path
	switch {
	case urlPath == "/api/saz":
		payload = postSaz(responseWriter, request)
	case strings.HasPrefix(urlPath, "/api/saz/"):
		payload = getSaz(responseWriter, request)
	default:
		message := fmt.Sprintf("Unknown path \"%s\"", urlPath)
		http.Error(responseWriter, message, http.StatusNotFound)
	}
	if payload != nil {
		sendJSON(urlPath, responseWriter, payload)
	}
}

func sendJSON(urlPath string, responseWriter http.ResponseWriter, payload interface{}) {
	output, err := json.Marshal(payload)
	if err != nil {
		message := fmt.Sprintf("Marshaling JSON response failed.\n%s", err.Error())
		http.Error(responseWriter, message, http.StatusInternalServerError)
		return
	}
	responseWriter.Header().Set("Content-Type", "application/json")
	_, err = io.WriteString(responseWriter, string(output))
	if err != nil {
		message := fmt.Sprintf("Sending %d bytes and \"application/json\" type of the response body for the request %s failed.",
			len(output), urlPath)
		message = fmt.Sprintf("%s\n%s", message, err.Error())
		http.Error(responseWriter, message, http.StatusInternalServerError)
		return
	}
}
