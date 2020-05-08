package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func main() {
	var mimeTypeIcons []string
	err := filepath.Walk("cmd/sazserve/assets/png",
		func(filePath string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				fileName := info.Name()
				fileName = strings.TrimSuffix(fileName, path.Ext(fileName))
				mimeTypeIcons = append(mimeTypeIcons, fileName)
			}
			return nil
		})
	if err != nil {
		panic(err)
	}
	content := `const mimeTypeIcons = new Set()

setTimeout(() => {
`
	for _, mimeTypeIcon := range mimeTypeIcons {
		content += fmt.Sprintf("  mimeTypeIcons.add('%s')\n", mimeTypeIcon)
	}
	content += `})

export default mimeTypeIcons
`
	fileName := os.Args[2]
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	_, err = file.WriteString(content)
	if err != nil {
		file.Close()
		panic(err)
	}
	err = file.Close()
	if err != nil {
		panic(err)
	}
}
