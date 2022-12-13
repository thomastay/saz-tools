package dumper_test

import (
	"github.com/prantlf/saz-tools/pkg/dumper"
	"github.com/prantlf/saz-tools/pkg/parser"
)

// Pprint a summary line on the console for each session from `foo.saz`.
func ExampleDump() {
	sessions, err := parser.ParseFile("foo.saz", "")
	if err != nil {
		panic(err)
	}
	err = dumper.Dump(sessions)
	if err != nil {
		panic(err)
	}
	// Output: Number Timeline Method Code URL Begin End Duration Size Encoding Cache Process
	// 1 00:00.000 GET 200 https://example.com/foo 16:11:58.755 16:11:59.013 00:00.257 506 gzip max-age=31536000 chromium:18160
	// 2 00:01.573 GET 200 https://example.com/bar 16:11:59.419 16:11:59.873 00:00.454 1,201 raw no-store,no-cache chromium:18160
	// ...
}
