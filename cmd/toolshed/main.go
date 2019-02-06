package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/belt-sh/toolshed"
)

const helpText = `usage: toolshed [options]

  --listen  The address to listen on
`

var (
	version = flag.Bool("version", false, "")
	listen  = flag.String("listen", ":8080", "")
)

func usage() {
	fmt.Fprintf(os.Stderr, helpText)
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if *version {
		fmt.Fprintf(os.Stdout, "toolshed %s (%s)\n", Version, Commit)
		os.Exit(0)
	}

	logger := log.New(os.Stderr, "[toolshed] ", log.LstdFlags)

	if err := toolshed.Run(*listen, logger); err != nil {
		logger.Fatalf("err: %v\n", err)
	}
}
