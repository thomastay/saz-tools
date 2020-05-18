package parser

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strconv"
)

var archivedFileName = regexp.MustCompile("(\\d+)_(\\w)")

func parseArchivedFileName(name string) (bool, int, string) {
	match := archivedFileName.FindAllStringSubmatch(name, -1)
	if len(match) == 0 {
		return false, 0, ""
	}
	number, _ := strconv.Atoi(match[0][1])
	return true, number, match[0][2]
}

func slurpRequestBody(session *Session) error {
	if session.Request.ContentLength <= 0 {
		return nil
	}
	var err error
	defer session.Request.Body.Close()
	if session.RequestBody, err = ioutil.ReadAll(session.Request.Body); err != nil {
		message := fmt.Sprintf("Reading %d bytes of the request body of the type \"%s\" failed.",
			session.Request.ContentLength, session.Request.Header.Get("Content-Type"))
		return fmt.Errorf("%s\n%s", message, err.Error())
	}
	return nil
}

func slurpResponseBody(session *Session) error {
	if session.Response.ContentLength <= 0 {
		return nil
	}
	var reader io.ReadCloser
	var err error
	switch session.Response.Header.Get("Content-Encoding") {
	case "gzip":
		if reader, err = gzip.NewReader(session.Response.Body); err != nil {
			message := fmt.Sprintf("Opening gzipped %d bytes of the response body of the type \"%s\" failed.",
				session.Response.ContentLength, session.Response.Header.Get("Content-Type"))
			return fmt.Errorf("%s\n%s", message, err.Error())
		}
	default:
		reader = session.Response.Body
	}
	defer reader.Close()
	if session.ResponseBody, err = ioutil.ReadAll(reader); err != nil {
		message := fmt.Sprintf("Reading %d bytes of the response body of the type \"%s\" failed.",
			session.Response.ContentLength, session.Response.Header.Get("Content-Type"))
		return fmt.Errorf("%s\n%s", message, err.Error())
	}
	return nil
}
