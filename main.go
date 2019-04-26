package main

import (
	"fmt"
	"os"

	"github.com/mui87/atcsubmit/cli"
)

const (
	codeOK = iota
	codeErr
)

func main() {
	c, err := cli.New(os.Args, os.Stdout, os.Stderr)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "ERROR: %s\n", err.Error())
		os.Exit(codeErr)
	}
	if err := c.Run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "ERROR: %s\n", err.Error())
		os.Exit(codeErr)
	}
	os.Exit(codeOK)
}
