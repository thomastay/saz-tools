# saz-tools

[![Build Status](https://travis-ci.org/prantlf/saz-tools.svg?branch=master)](https://travis-ci.org/prantlf/saz-tools)
[![Go Report Card](https://goreportcard.com/badge/github.com/prantlf/saz-tools)](https://goreportcard.com/report/github.com/prantlf/saz-tools)

Tools for parsing SAZ files (Fiddler logs) and either printing their content on the console, or viewing them on a web page and offering basic analysis and export. Try the [on-line version](https://viewsaz.herokuapp.com/).

## Tools

```
$ go install github.com/prantlf/saz-tools

$ sazdump foo.saz
Number Timeline Method Code URL Begin End Duration Size Encoding Cache Process
1 00:00.000 GET 200 https://example.com/foo 16:11:58.755 16:11:59.013 00:00.257 506 gzip max-age=31536000 chromium:18160
2 00:01.573 GET 200 https://example.com/bar 16:11:59.419 16:11:59.873 00:00.454 1,201 raw no-store,no-cache chromium:18160
...

$ sazserve
$ open http://localhost:7000
```

## API

```go
import (
  sazanalyzer "github.com/prantlf/saz-tools/pkg/analyzer"
  sazdumper "github.com/prantlf/saz-tools/pkg/dumper"
  sazparser "github.com/prantlf/saz-tools/pkg/parser"
)

func main() {
  rawSessions, err := sazparser.ParseFile("foo.saz")
  fineSessions, err := sazparser.Analyze(rawSessions)
  err = sazdumper.Dump(rawSessions)
}
```

## License

MIT
