package main

import (
	"flag"
	"fmt"
	"os"

	sazdumper "github.com/prantlf/sazdump/dumper"
	sazparser "github.com/prantlf/sazdump/parser"
)

func main() {
	flag.Parse()
	sessions, err := sazparser.Parse(flag.Arg(0))
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	err = sazdumper.Dump(sessions)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
