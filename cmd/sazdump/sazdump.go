package main

import (
	"flag"
	"fmt"
	"os"

	sazdumper "github.com/prantlf/saz-tools/pkg/dumper"
	sazparser "github.com/prantlf/saz-tools/pkg/parser"
)

func main() {
	flag.Parse()
	sessions, err := sazparser.ParseFile(flag.Arg(0))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = sazdumper.Dump(sessions)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
