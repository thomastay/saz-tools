package main

import (
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func main() {
	fileName := os.Args[2]
	originalContent, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	byLineBreaks := regexp.MustCompile("\\r?\\n")
	lines := byLineBreaks.Split(string(originalContent), -1)
	var firstCodeLineIndex int
	for lineIndex, line := range lines {
		line = strings.TrimLeft(line, " \t")
		if !(line == "" || strings.HasPrefix(line, "//")) {
			firstCodeLineIndex = lineIndex
			break
		}
	}
	modifiedContent := lines[firstCodeLineIndex] + "\n" +
		strings.Join(lines[:firstCodeLineIndex], "\n") +
		strings.Join(lines[firstCodeLineIndex+1:], "\n")
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	_, err = file.WriteString(modifiedContent)
	if err != nil {
		file.Close()
		panic(err)
	}
	err = file.Close()
	if err != nil {
		panic(err)
	}
}
