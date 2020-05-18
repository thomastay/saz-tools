package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func main() {
	fileName := os.Args[2]
	originalContent, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Reading \"%s\" failed.\n", fileName)
		panic(err)
	}
	byLineBreaks := regexp.MustCompile(`\r?\n`)
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
		fmt.Printf("Creating \"%s\" failed.\n", fileName)
		panic(err)
	}
	_, err = file.WriteString(modifiedContent)
	if err != nil {
		file.Close()
		fmt.Printf("Writing %d characters to \"%s\" failed.\n", len(modifiedContent), fileName)
		panic(err)
	}
	err = file.Close()
	if err != nil {
		fmt.Printf("Closing \"%s\" failed.\n", fileName)
		panic(err)
	}
}
