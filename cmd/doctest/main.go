package main

import (
	"fmt"
	"os"

	"github.com/xhd2015/doctest/libdoc/debug"
	"github.com/xhd2015/doctest/run"
)

func main() {
	// Engine-internal DOCTEST_DEBUG: parse once at process start so host
	// profiles cover every subcommand (test, vet, agent, metrics, …).
	settings, err := debug.FromEnv()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	stop, err := debug.StartProfiles(settings)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer stop()

	err = run.Run(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
