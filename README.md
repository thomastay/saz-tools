# saz-tools

[![Build Status](https://travis-ci.org/prantlf/saz-tools.svg?branch=master)](https://travis-ci.org/prantlf/saz-tools)
[![Go Report Card](https://goreportcard.com/badge/github.com/prantlf/saz-tools)](https://goreportcard.com/report/github.com/prantlf/saz-tools)
[![Documentation](https://godoc.org/github.com/prantlf/saz-tools?status.svg)](http://godoc.org/github.com/prantlf/saz-tools)

Tools for parsing SAZ files (Fiddler logs) and either [printing their content] on the console, or [viewing them on a web page] and offering basic analysis and export. Try the [on-line version] of the SAZ Viewer.

## Installation

If you work on a Go module, using [`go get`] is the easiest way:

    $ go get -u github.com/prantlf/saz-tools/...

You can also install the latest version of the tools using [GoBinaries]:

    curl -sf https://gobinaries.com/prantlf/saz-tools | sh

Or choose a version from [GitHub releases].

## Tools

```
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

* [sazparser].[ParseFile](fileName string) ([] [sazparser.Session], error)
* [sazparser].[ParseReader](reader ReaderAt, size int64) ([] [sazparser.Session], error)
* [sazanalyzer].[Analyze](sessions [] [sazparser.Session]) ([] [sazanalyzer.Session], error)
* [sazdumper].[Dump](sessions [] [sazparser.Session]) error

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

## Build

You need [Go], [Make], [Node.js], [NPM] and [Patch] to build all parts of this module from sources.

    # Install build tools only once:
    make prepare
    # Build all targets:
    make
    # Dump a SAZ file analysis on the console:
    make run-dump SAZ="..."
    # Start a SAZ file viewer as a browser application with bundled assets:
    make run-serve
    # Start a SAZ file viewer as a browser application with filesystem assets:
    make debug-serve
    # Remove all build output:
    make clean

## Docker

[![prantlf/sazdump](http://dockeri.co/image/prantlf/sazdump)](https://hub.docker.com/repository/docker/prantlf/sazdump/) [![prantlf/sazserve](http://dockeri.co/image/prantlf/sazserve)](https://hub.docker.com/repository/docker/prantlf/sazserve/)

[The `sazdump` image] and [the `sazserve` image] allows you to execute the tools described above. They are built automatically on the top of the tag `latest` from the [scratch image].

The following [tags] are available for the `prantlf/saztools` image:

- `latest`

Download the latest image to your disk:

    docker pull prantlf/sazdump
    docker pull prantlf/sazserve
    # or
    docker pull prantlf/sazdump:latest
    docker pull prantlf/sazserve:latest

Print usage description with command-line parameters:

    docker run --rm -it prantlf/sazdump -h

For example, dump the context of the `foo.saz` file:

    docker run --rm -it -v ${PWD}:/work -w /work sazdump foo.saz

Or start the browser application to analyse SAZ files on your machine:

    docker run --rm -it -p 7000:7000 sazserve

You can also put [`sazdump`] and [`sazserve`] scripts to `PATH`:

    #!/bin/sh
    docker run --rm -it -v ${PWD}:/work -w /work sazdump "$@"

    #!/bin/sh
    docker run --rm -it -v -p 7000:7000 sazserve

and execute then from any location by supplying parameters to it, for example:

    sazdump foo.saz
    sazserve

Local images are built as `sazdump` and `sazserve` and they are pushed to the docker hub as `prantlf/sazdump:latest` and `prantlf/sazserve:latest`.

    # Check the Dockerfiles:
    make docker-lint
    # Build new local images:
    make docker-build
    # Print the help for the SAZ file dumper:
    make docker-run-help
    # Dump a SAZ file analysis on the console:
    make docker-dump-example SAZ="..."
    # Start a SAZ file viewer as a web application:
    make docker-serve-example
    # Tag local images for pushing:
    make docker-tag
    # Login to the docker hub:
    make docker-login
    # Push local images to the docker hub:
    make docker-push

## License

Copyright (c) 2020 Ferdinand Prantl

Licensed under the MIT license.

[on-line version]: https://viewsaz.herokuapp.com/
[`go get`]: https://golang.org/cmd/go/#hdr-Add_dependencies_to_current_module_and_install_them
[Go]: https://golang.org/
[golang repository]: https://hub.docker.com/_/golang
[Make]: https://www.gnu.org/software/make/
[Patch]: http://man7.org/linux/man-pages/man1/patch.1.html
[Node.js]: https://nodejs.org/
[NPM]: https://docs.npmjs.com/cli/npm
[GoBinaries]: https://gobinaries.com/
[GitHub releases]: https://github.com/prantlf/saz-tools/releases
[The `sazdump` image]: https://hub.docker.com/repository/docker/prantlf/sazdump
[the `sazserve` image]: https://hub.docker.com/repository/docker/prantlf/sazserve
[tags]: https://hub.docker.com/repository/docker/prantlf/saztools/tags
[scratch image]: https://hub.docker.com/_/scratch
[`sazdump`]: bin/sazdump
[`sazserve`]: bin/sazserve
[printing their content]: https://godoc.org/github.com/prantlf/saz-tools/cmd/sazdump
[viewing them on a web page]: https://godoc.org/github.com/prantlf/saz-tools/cmd/sazserve
[sazparser]: https://godoc.org/github.com/prantlf/saz-tools/pkg/parser
[sazparser.Session]: https://godoc.org/github.com/prantlf/saz-tools/pkg/parser#Session
[ParseFile]: https://godoc.org/github.com/prantlf/saz-tools/pkg/parser#ParseFile
[ParseReader]: https://godoc.org/github.com/prantlf/saz-tools/pkg/parser#ParseReader
[sazanalyzer]: https://godoc.org/github.com/prantlf/saz-tools/pkg/analyzer
[sazanalyzer.Session]: https://godoc.org/github.com/prantlf/saz-tools/pkg/analyzer#Session
[Analyze]: https://godoc.org/github.com/prantlf/saz-tools/pkg/analyzer#Analyze
[sazdumper]: https://godoc.org/github.com/prantlf/saz-tools/pkg/dumper
[Dump]: https://godoc.org/github.com/prantlf/saz-tools/pkg/dumper#Dump
