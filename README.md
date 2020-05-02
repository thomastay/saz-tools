# sazdump

Parses a SAZ files (Fiddler logs) and prints their content on the console.

## Tool

```
$ go install github.com/prantlf/sazdump
$ sazdump foo.saz
Number Timeline Method Code URL Begin End Duration Size Encoding Cache Process
1 00:00.000 GET 200 https://example.com/foo 16:11:58.755 16:11:59.013 00:00.257 506 gzip max-age=31536000 chromium:18160
2 00:01.573 GET 200 https://example.com/bar 16:11:59.419 16:11:59.873 00:00.454 1,201 raw no-store,no-cache chromium:18160
...
```

## API

```go
import (
	sazdumper "github.com/prantlf/sazdump/dumper"
	sazparser "github.com/prantlf/sazdump/parser"
)

func main() {
	sessions, err := sazparser.Parse("foo.saz")
	err = sazdumper.Dump(sessions)
}
```

## License

MIT
