# saz-tools

[![Build Status](https://travis-ci.org/prantlf/saz-tools.svg?branch=master)](https://travis-ci.org/prantlf/saz-tools)
[![Go Report Card](https://goreportcard.com/badge/github.com/prantlf/saz-tools)](https://goreportcard.com/report/github.com/prantlf/saz-tools)

Tools for parsing SAZ files (Fiddler logs) and either printing their content on the console, or viewing them on a web page and offering basic analysis and export. Try the [on-line version].

## Tools

```
$ go install github.com/prantlf/saz-tools

$ sazdump.go -h
Usage: sazdump <file.saz>

$ sazdump foo.saz
Number Timeline Method Code URL Begin End Duration Size Encoding Cache Process
1 00:00.000 GET 200 https://example.com/foo 16:11:58.755 16:11:59.013 00:00.257 506 gzip max-age=31536000 chromium:18160
2 00:01.573 GET 200 https://example.com/bar 16:11:59.419 16:11:59.873 00:00.454 1,201 raw no-store,no-cache chromium:18160
...

$ sazserve -h
Usage: sazserve [options]
Options:
  -port string : port for the web server to listen to (default "7000")

$ sazserve
$ open http://localhost:7000
```

## API

* parser.ParseFile(fileName string) ([]parser.Sessions, error)
* parser.ParseReader(reader ReaderAt, size int64) ([]parser.Sessions, error)
* analyzer.Analyze(sessions []parser.Sessions) ([]analyzer.Sessions, error)
* dumper.Dump(sessions []parser.Sessions) error

```go
import (
  sazanalyzer "github.com/prantlf/saz-tools/pkg/analyzer"
  sazdumper "github.com/prantlf/saz-tools/pkg/dumper"
  sazparser "github.com/prantlf/saz-tools/pkg/parser"
)

func main() {
  rawSessions, _ := sazparser.ParseFile("foo.saz")
  fineSessions, _ := sazparser.Analyze(rawSessions)
  sazdumper.Dump(rawSessions)
}
```

## Docker

[![nodesource/node](http://dockeri.co/image/prantlf/saztools)](https://hub.docker.com/repository/docker/prantlf/saztools/)

[This image] allows you to execute the tools desribed above. It is built automatically on the top of the tag `latest` from the [scratch image].

The following [tags] are available for the `prantlf/saztools` image:

- `latest`

Download the latest image to your disk:

    docker pull prantlf/saztools
    # or
    docker pull prantlf/saztools:latest

Print usage description with command-line parameters:

    docker run --rm -it prantlf/saztools -h

For example, dump the context of the `foo.saz` file:

    docker run --rm -it -v ${PWD}:/work -w /work saztools foo.saz

You can also put a [`sazdump`] script to `PATH`:

    #!/bin/sh
    docker run --rm -it -v ${PWD}:/work -w /work saztools "$@"

and execute it from any location by supplying parameters to it, for example:

    sazdump foo.saz

The local image is built as `saztools` and pushed to the docker hub as `prantlf/saztools:latest`.

    # Remove an old local image:
    make clean
    #  Check the Dockerfile:
    make lint
    # Build a new local image:
    make build
    # Print the help for the diagram generator:
    make run-help
    # Generate an image for a diagram sample:
    make run-example
    # Tag the local image for pushing:
    make tag
    # Login to the docker hub:
    make login
    # Push the local image to the docker hub:
    make push

## License

Copyright (c) 2020 Ferdinand Prantl

Licensed under the MIT license.

[on-line version]: https://viewsaz.herokuapp.com/
[This image]: https://hub.docker.com/repository/docker/prantlf/saztools
[tags]: https://hub.docker.com/repository/docker/prantlf/saztools/tags
[scratch image]: https://hub.docker.com/_/scratch
[`sazdump`]: bin/sazdump
