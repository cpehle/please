// package main implements various tools that Please finds useful for Go.
//
// Firstly, it implements a code templater for Go tests.
// This is essentially equivalent to what 'go test' does but it lifts some restrictions
// on file organisation and allows the code to be instrumented for coverage as separate
// build targets rather than having to repeat it for every test.
package main

import (
	"os"

	"gopkg.in/op/go-logging.v1"

	"cli"
	"tools/please_go_tool/testmain"
)

var log = logging.MustGetLogger("plz_go_test")

var opts = struct {
	Usage     string
	Dir       string   `short:"d" long:"dir" description:"Directory to search for Go package files for coverage"`
	Verbosity int      `short:"v" long:"verbose" default:"1" description:"Verbosity of output (higher number = more output, default 1 -> warnings and errors only)"`
	Exclude   []string `short:"x" long:"exclude" default:"third_party/go" description:"Directories to exclude from search"`
	Output    string   `short:"o" long:"output" description:"Output filename" required:"true"`
	Package   string   `short:"p" long:"package" description:"Package containing this test" env:"PKG"`
	Args      struct {
		Go      string   `positional-arg-name:"go" description:"Location of go command" required:"true"`
		Sources []string `positional-arg-name:"sources" description:"Test source files" required:"true"`
	} `positional-args:"true" required:"true"`
}{
	Usage: `
please_go_tool implements various tools that Please finds useful for building Go code.

Firstly, it implements a code templater for Go tests.
This is essentially equivalent to what 'go test' does but it lifts some restrictions
on file organisation and allows the code to be instrumented for coverage as separate
build targets rather than having to repeat it for every test.
`,
}

func main() {
	cli.ParseFlagsOrDie("plz_go_test", "7.2.0", &opts)
	cli.InitLogging(opts.Verbosity)
	coverVars, err := testmain.FindCoverVars(opts.Dir, opts.Exclude, opts.Args.Sources)
	if err != nil {
		log.Fatalf("Error scanning for coverage: %s", err)
	}
	if err = testmain.WriteTestMain(opts.Package, testmain.IsVersion18(opts.Args.Go), opts.Args.Sources, opts.Output, coverVars); err != nil {
		log.Fatalf("Error writing test main: %s", err)
	}
	os.Exit(0)
}
