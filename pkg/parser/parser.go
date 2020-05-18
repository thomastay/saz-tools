// Package parser parses SAZ files (Fiddler logs) to an array of sessions,
// which contain all about network connections, requests and responses.
package parser

import (
	"archive/zip"
	"bufio"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/prantlf/saz-tools/internal/pluralizer"
)

// ParseFile prses a file to an array of network sessions.
func ParseFile(fileName string) ([]Session, error) {
	archiveReader, err := zip.OpenReader(fileName)
	if err != nil {
		message := fmt.Sprintf("Opening zipped \"%s\" failed.", fileName)
		return nil, fmt.Errorf("%s\n%s", message, err.Error())
	}
	defer archiveReader.Close()
	return parseArchive(&archiveReader.Reader)
}

// ParseReader parses a file content passed by a reader to an array of network sessions.
func ParseReader(reader io.ReaderAt, size int64) ([]Session, error) {
	archiveReader, err := zip.NewReader(reader, size)
	if err != nil {
		message := fmt.Sprintf("Opening zipped %d bytes failed.", size)
		return nil, fmt.Errorf("%s\n%s", message, err.Error())
	}
	return parseArchive(archiveReader)
}

func parseArchive(archiveReader *zip.Reader) ([]Session, error) {
	var session Session
	var sessions []Session

	for _, archivedFile := range archiveReader.File {
		match, number, fileType := parseArchivedFileName(archivedFile.Name)
		if !match {
			continue
		}

		switch fileType {
		case "c":
			err := parseRequest(archivedFile, &session, number)
			if err != nil {
				return nil, err
			}

		case "m":
			err := parseSessionAttributes(archivedFile, &session)
			if err != nil {
				return nil, err
			}

		case "s":
			err := parseResponse(archivedFile, &session, number)
			if err != nil {
				return nil, err
			}

			session.Number = number
			sessions = append(sessions, session)
			session = Session{}
		}
	}

	if len(sessions) == 0 {
		return nil, errors.New("saz/parser: no network sessions were found")
	}
	return sessions, nil
}

func parseRequest(archivedFile *zip.File, session *Session, number int) error {
	fileReader, err := archivedFile.Open()
	if err != nil {
		message := fmt.Sprintf("Opening \"%s\" failed.", archivedFile.Name)
		return fmt.Errorf("%s\n%s", message, err.Error())
	}
	defer fileReader.Close()

	requestReader := bufio.NewReader(fileReader)
	request, err := http.ReadRequest(requestReader)
	if err != nil {
		message := fmt.Sprintf("Reading request from \"%s\" failed.", archivedFile.Name)
		return fmt.Errorf("%s\n%s", message, err.Error())
	}
	defer request.Body.Close()
	session.Request = request

	err = slurpRequestBody(session)
	if err != nil {
		message := fmt.Sprintf("Buffering request body from %s network session failed.",
			pluralizer.FormatOrdinal(number))
		return fmt.Errorf("%s\n%s", message, err.Error())
	}

	return nil
}

func parseResponse(archivedFile *zip.File, session *Session, number int) error {
	fileReader, err := archivedFile.Open()
	if err != nil {
		message := fmt.Sprintf("Opening \"%s\" failed.", archivedFile.Name)
		return fmt.Errorf("%s\n%s", message, err.Error())
	}
	defer fileReader.Close()

	responseReader := bufio.NewReader(fileReader)
	response, err := http.ReadResponse(responseReader, session.Request)
	if err != nil {
		message := fmt.Sprintf("Reading response from \"%s\" failed.", archivedFile.Name)
		return fmt.Errorf("%s\n%s", message, err.Error())
	}
	defer response.Body.Close()
	session.Response = response

	err = slurpResponseBody(session)
	if err != nil {
		message := fmt.Sprintf("Buffering response body from %s network session failed.",
			pluralizer.FormatOrdinal(number))
		return fmt.Errorf("%s\n%s", message, err.Error())
	}

	return nil
}

func parseSessionAttributes(archivedFile *zip.File, session *Session) error {
	fileReader, err := archivedFile.Open()
	if err != nil {
		message := fmt.Sprintf("Opening \"%s\" failed.", archivedFile.Name)
		return fmt.Errorf("%s\n%s", message, err.Error())
	}
	defer fileReader.Close()

	bytes, err := ioutil.ReadAll(fileReader)
	if err != nil {
		message := fmt.Sprintf("Reading session timers and flags from \"%s\" failed.",
			archivedFile.Name)
		return fmt.Errorf("%s\n%s", message, err.Error())
	}

	err = xml.Unmarshal(bytes, &session)
	if err != nil {
		message := fmt.Sprintf("Unmarshaling session timers and flags from %d bytes of \"%s\" failed.",
			len(bytes), archivedFile.Name)
		return fmt.Errorf("%s\n%s", message, err.Error())
	}

	return nil
}
