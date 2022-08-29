package craps

import (
	"flag"
	reallog "log"
)

var (
	verbose = flag.Bool("verbose", false, "verbose logging")
)

func log(tmpl string, args ...interface{}) {
	if *verbose {
		reallog.Printf(tmpl, args...)
	}
}
