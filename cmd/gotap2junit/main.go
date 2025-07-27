package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/clvf/gotap2junit/transform"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(
			flag.CommandLine.Output(),
			"SYNOPSIS\n\t%v < $tap13.report > $junit.xml\n\n",
			filepath.Base(os.Args[0]))
		fmt.Fprintf(flag.CommandLine.Output(),
			"DESCRIPTION\n\tParse TAP version 13 report on stdin and transform it "+
				"to Junit XML format on stdout.\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	transform.Run(os.Stdin, os.Stdout)
}
