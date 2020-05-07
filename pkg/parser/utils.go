package parser

import (
	"regexp"
	"strconv"
)

var archivedFileName *regexp.Regexp

func init() {
	archivedFileName, _ = regexp.Compile("(\\d+)_(\\w)")
}

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
