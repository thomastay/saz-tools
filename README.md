# Fork of saz-tools

# Changelog:

1. More lenient parsing even if the HTTP requests don't meet spec
1. Enable parsing of password protected SAZ files
1. Upgrade go.mod
1. Print Response body / headers to output file or stdout, per curl syntax

### Installation

Requires Go 1.19

```
go install https://github.com/thomastay/saz-tools/cmd/sazdump
go install https://github.com/thomastay/saz-tools/cmd/sazserve
```

### Usage:

```
sazdump archive.saz
    Dumps archive.sav to stdout

sazdump -n 10 -o body.json archive.saz
    Saves response body number 10 to body.json

sazdump -n 10 -o body.json -p password1234 archive.saz
    Saves response body number 10 to body.json with password protected archive.sav

sazdump -n 10 -I archive.saz
    Prints headers to stdout
```

# Old README

[![Build Status](https://github.com/prantlf/saz-tools/workflows/Test/badge.svg)](https://github.com/prantlf/saz-tools/actions)
[![Dependency Status](https://david-dm.org/prantlf/saz-tools.svg)](https://david-dm.org/prantlf/saz-tools)
[![devDependency Status](https://david-dm.org/prantlf/saz-tools/dev-status.svg)](https://david-dm.org/prantlf/saz-tools#info=devDependencies)
[![Go Report Card](https://goreportcard.com/badge/github.com/prantlf/saz-tools)](https://goreportcard.com/report/github.com/prantlf/saz-tools)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/prantlf/saz-tools?color=teal)
[![Documentation](https://godoc.org/github.com/prantlf/saz-tools?status.svg)](http://godoc.org/github.com/prantlf/saz-tools)
[![Homebrew](https://img.shields.io/badge/dynamic/json.svg?url=https://raw.githubusercontent.com/prantlf/homebrew-tap/master/Info/saz-tools.json&query=$.versions.stable&label=homebrew)](https://github.com/prantlf/homebrew-tap#readme)
[![Snap](https://img.shields.io/badge/dynamic/json.svg?url=https://raw.githubusercontent.com/prantlf/saz-tools/master/package.json&query=$.version&label=snap)](https://snapcraft.io/saz-tools)
[![Scoop](https://img.shields.io/badge/dynamic/json.svg?url=https://raw.githubusercontent.com/prantlf/scoop-bucket/master/saz-tools.json&query=$.version&label=scoop)](https://github.com/prantlf/scoop-bucket#readme)
[![GitHub release (latest SemVer including pre-releases)](https://img.shields.io/github/v/release/prantlf/saz-tools?include_prereleases&label=github%2Fdeb%2Frpm)](https://github.com/prantlf/saz-tools/releases)
[![npm](https://img.shields.io/npm/v/saz-tools)](https://www.npmjs.com/package/saz-tools#top)
![Docker Image Version (latest by date)](https://img.shields.io/docker/v/prantlf/sazdump?color=cyan&label=docker)

Tools for parsing SAZ files (Fiddler logs) and either [printing their content] on the console, or [viewing them on a web page] and offering basic analysis and export. Try the [on-line version] of the SAZ Viewer.

## Installation

If you have [Go] installed, using [`go get`] to install a global module is the easiest way:

    $ GO111MODULE=off go get -u github.com/prantlf/saz-tools/...

If you have [Node.js] installed, you can use [NPM], [Yarn] or [PNPM] to install a global module easily:

    npm i -g saz-tools
    yarn global add saz-tools
    pnpm i -g saz-tools

If you have the standard `sh` available, you can use the installation script from [GoBinaries]:

    curl -sf https://gobinaries.com/prantlf/saz-tools | sh

If you manage your software using [Homebrew], you can install the tools using their formula:

    brew install prantlf/tap/saz-tools

Windows users can install using the [Scoop manifest]:

    scoop bucket add prantlf https://github.com/prantlf/scoop-bucket.git
    scoop install prantlf/saz-tools

Ubuntu users can install the [Snap package]:

    sudo snap install saz-tools

If you work on Linux which uses `deb` or `rpm` packages, you can download and install a package from [GitHub releases].

Or download and unpack a binary archive for your operation system from [GitHub releases] directly.

If you want to install a specific commit or the latest master and you do not have the development environment to build it, you can use Docker to [`build`]:

    git clone https://github.com/prantlf/saz-tools.git
    cd saz-tools
    docker run --rm -it -v ${PWD}:/work -w /work \
      prantlf/golang-make-nodejs-git clean prepare all DOCKER=1

If you want to run the tools using [Docker] images, see the [instructions below](#docker).

## Tools

```
$ sazdump.go -h
Usage: sazdump [options] <file.saz>
Options:
  -version | -v : print the version of this tool and exit

$ sazdump foo.saz
Number Timeline Method Code URL Begin End Duration Size Encoding Cache Process
1 00:00.000 GET 200 https://example.com/foo 16:11:58.755 16:11:59.013 00:00.257 506 gzip max-age=31536000 chromium:18160
2 00:01.573 GET 200 https://example.com/bar 16:11:59.419 16:11:59.873 00:00.454 1,201 raw no-store,no-cache chromium:18160
...

$ sazserve -h
Usage: sazserve [options]
Options:
  -browser       : start the web browser automatically  (default false)
  -port <number> : port for the web server to listen on (default "7000")
  -version | -v  : print the version of this tool and exit

$ sazserve
$ open http://localhost:7000/
```

## API

- [parser].[ParseFile](fileName string) ([] [parser.Session], error)
- [parser].[ParseReader](reader ReaderAt, size int64) ([] [parser.Session], error)
- [analyzer].[Analyze](sessions [] [parser.Session]) ([] [analyzer.Session], error)
- [dumper].[Dump](sessions [] [parser.Session]) error

```go
import (
  "github.com/prantlf/saz-tools/pkg/analyzer"
  "github.com/prantlf/saz-tools/pkg/dumper"
  "github.com/prantlf/saz-tools/pkg/parser"
)

func main() {
  rawSessions, _ := parser.ParseFile("foo.saz")
  fineSessions, _ := analyzer.Analyze(rawSessions)
  dumper.Dump(rawSessions)
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

The following tags are available for the `prantlf/sazdump` and `prantlf/sazserve` images:

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

## Release

    # make all changes, bump the version and tag it
    make clean prepare lint all
    conventional-changelog -p angular -i CHANGELOG.md -s
    # update the tag
    snapcraft login --with snap.login
    goreleaser release --rm-dist
    # pull and update brew, update scoop
    # push to heroku

## License

Copyright (c) 2020-2021 Ferdinand Prantl

Licensed under the MIT license.

[on-line version]: https://viewsaz.herokuapp.com/
[`go get`]: https://golang.org/cmd/go/#hdr-Add_dependencies_to_current_module_and_install_them
[go]: https://golang.org/
[golang repository]: https://hub.docker.com/_/golang
[homebrew]: https://brew.sh/
[snap package]: https://snapcraft.io/saz-tools
[scoop manifest]: https://github.com/prantlf/scoop-bucket#prantlfscoop-bucket
[docker]: https://www.docker.com/
[make]: https://www.gnu.org/software/make/
[patch]: http://man7.org/linux/man-pages/man1/patch.1.html
[node.js]: https://nodejs.org/
[npm]: https://docs.npmjs.com/cli/npm
[yarn]: https://classic.yarnpkg.com/docs/cli/
[pnpm]: https://pnpm.js.org/pnpm-cli
[gobinaries]: https://gobinaries.com/
[github releases]: https://github.com/prantlf/saz-tools/releases
[the `sazdump` image]: https://hub.docker.com/repository/docker/prantlf/sazdump
[the `sazserve` image]: https://hub.docker.com/repository/docker/prantlf/sazserve
[scratch image]: https://hub.docker.com/_/scratch
[`build`]: bin/build
[`sazdump`]: bin/sazdump
[`sazserve`]: bin/sazserve
[printing their content]: https://godoc.org/github.com/prantlf/saz-tools/cmd/sazdump
[viewing them on a web page]: https://godoc.org/github.com/prantlf/saz-tools/cmd/sazserve
[parser]: https://godoc.org/github.com/prantlf/saz-tools/pkg/parser
[parser.session]: https://godoc.org/github.com/prantlf/saz-tools/pkg/parser#Session
[parsefile]: https://godoc.org/github.com/prantlf/saz-tools/pkg/parser#ParseFile
[parsereader]: https://godoc.org/github.com/prantlf/saz-tools/pkg/parser#ParseReader
[analyzer]: https://godoc.org/github.com/prantlf/saz-tools/pkg/analyzer
[analyzer.session]: https://godoc.org/github.com/prantlf/saz-tools/pkg/analyzer#Session
[analyze]: https://godoc.org/github.com/prantlf/saz-tools/pkg/analyzer#Analyze
[dumper]: https://godoc.org/github.com/prantlf/saz-tools/pkg/dumper
[dump]: https://godoc.org/github.com/prantlf/saz-tools/pkg/dumper#Dump
