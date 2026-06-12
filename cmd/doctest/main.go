package main

import (
	"fmt"
	"os"

	"github.com/xhd2015/doctest/run"
)

func main() {
	err := run.Run(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
