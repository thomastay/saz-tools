package parser

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"regexp"
	"strconv"
)

var archivedFileName = regexp.MustCompile("(\\d+)_(\\w)")

func parseArchivedFileName(name string) (bool, int, string, error) {
	match := archivedFileName.FindAllStringSubmatch(name, -1)
	if len(match) == 0 {
		return false, 0, "", nil
	}
	number, err := strconv.Atoi(match[0][1])
	if err != nil {
		return false, 0, "", err
	}
	return true, number, match[0][2], nil
}

func slurpRequestBody(session *Session) error {
	if session.Request.ContentLength <= 0 {
		return nil
	}
	var err error
	defer session.Request.Body.Close()
	if session.RequestBody, err = ioutil.ReadAll(session.Request.Body); err != nil {
		return err
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
			return err
		}
	default:
		reader = session.Response.Body
	}
	defer reader.Close()
	if session.ResponseBody, err = ioutil.ReadAll(reader); err != nil {
		return err
	}
	return nil
}
