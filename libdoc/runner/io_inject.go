package runner

import (
	"io"

	"github.com/xhd2015/doctest/libdoc/core"
)

// applyWriters sets opts.Stdout/Stderr when non-nil. Used by TestWithWriters /
// processArgsWithWriters so nested in-process CLI can capture without package
// globals (callers pass buffers; concurrent leaves stay independent).
func applyWriters(opts *core.Options, stdout, stderr io.Writer) {
	if stdout != nil {
		opts.Stdout = stdout
	}
	if stderr != nil {
		opts.Stderr = stderr
	}
}
