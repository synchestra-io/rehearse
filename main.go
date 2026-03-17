package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/synchestra-io/rehearse/pkg/cli"
)

func main() {
	fatal := func(err error) {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
		type exitCoder interface{ ExitCode() int }
		var ec exitCoder
		if errors.As(err, &ec) {
			os.Exit(ec.ExitCode())
			return
		}
		os.Exit(1)
	}
	cli.Run(os.Args, fatal)
}
