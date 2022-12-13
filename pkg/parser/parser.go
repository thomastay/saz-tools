// Package parser parses SAZ files (Fiddler logs) to an array of sessions,
// which contain all about network connections, requests and responses.
package parser

import (
	"bufio"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/alexmullins/zip"

	"github.com/thomastay/saz-tools/internal/pluralizer"
)

// ParseFile prses a file to an array of network sessions.
func ParseFile(fileName string, password string) ([]Session, error) {
	archiveReader, err := zip.OpenReader(fileName)
	if err != nil {
		message := fmt.Sprintf("Opening zipped \"%s\" failed.", fileName)
		return nil, fmt.Errorf("%s\n%s", message, err.Error())
	}
	defer archiveReader.Close()
	return parseArchive(&archiveReader.Reader, password)
}

// ParseReader parses a file content passed by a reader to an array of network sessions.
func ParseReader(reader io.ReaderAt, size int64) ([]Session, error) {
	archiveReader, err := zip.NewReader(reader, size)
	if err != nil {
		message := fmt.Sprintf("Opening zipped %d bytes failed.", size)
		return nil, fmt.Errorf("%s\n%s", message, err.Error())
	}
	return parseArchive(archiveReader, "")
}

func parseArchive(archiveReader *zip.Reader, password string) ([]Session, error) {
	count, err := countSessions(archiveReader)
	if err != nil {
		return nil, err
	}
	sessions := make([]Session, count/3)

	for _, archivedFile := range archiveReader.File {
		if password != "" {
			archivedFile.SetPassword(password)
		}
		match, number, fileType := parseArchivedFileName(archivedFile.Name)
		if !match {
			continue
		}
		session := &sessions[number-1]

		switch fileType {
		case "c":
			err := parseRequest(archivedFile, session)
			if err != nil {
				log.Printf("Could not parse request %d, skipping. Err: %s\n", number, err.Error())
			}

		case "m":
			session.Number = number
			err := parseSessionAttributes(archivedFile, session)
			if err != nil {
				log.Printf("Could not parse session attributes %d, skipping. Err: %s\n", number, err.Error())
			}

		case "s":
			err := parseResponse(archivedFile, session)
			if err != nil {
				log.Printf("Could not parse response %d, skipping. Err: %s\n", number, err.Error())
			}
		}
	}
	return checkSessions(sessions)
}

func countSessions(archiveReader *zip.Reader) (int, error) {
	count := 0
	for _, archivedFile := range archiveReader.File {
		match, _, _ := parseArchivedFileName(archivedFile.Name)
		if match {
			count++
		}
	}
	if count == 0 {
		return 0, errors.New("saz/parser: no network sessions were found")
	}
	if count%3 != 0 {
		return 0, errors.New("saz/parser: incomplete file triplet detected")
	}
	return count, nil
}

func parseRequest(archivedFile *zip.File, session *Session) error {
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
		message := fmt.Sprintf("Buffering request body from \"%s\" failed.", archivedFile.Name)
		return fmt.Errorf("%s\n%s", message, err.Error())
	}

	return nil
}

func parseResponse(archivedFile *zip.File, session *Session) error {
	if archivedFile.UncompressedSize == 0 {
		session.Response = &http.Response{
			Status: "Connection Closed", Request: session.Request,
		}
		return nil
	}

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
		message := fmt.Sprintf("Buffering response body from \"%s\" failed.", archivedFile.Name)
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

	bytes, err := io.ReadAll(fileReader)
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

func checkSessions(sessions []Session) ([]Session, error) {
	for i, session := range sessions {
		if session.Number == 0 {
			return nil, fmt.Errorf("saz/parser: attributes missing in %s network session",
				pluralizer.FormatOrdinal(i))
		}
		if session.Request == nil {
			return nil, fmt.Errorf("saz/parser: request data missing in %s network session",
				pluralizer.FormatOrdinal(i))
		}
		if session.Response == nil {
			return nil, fmt.Errorf("saz/parser: Response data missing in %s network session",
				pluralizer.FormatOrdinal(i))
		}
		if session.Request.URL.String() == "" {
			return nil, fmt.Errorf("saz/parser: request data missing in %s network session",
				pluralizer.FormatOrdinal(i))
		}
		if session.Response.Request == nil {
			return nil, fmt.Errorf("saz/parser: response data missing in %s network session",
				pluralizer.FormatOrdinal(i))
		}
	}
	return sessions, nil
}
