package build

import (
	"fmt"
	"io"
	"os"

	"github.com/xhd2015/doctest/libdoc/core"
)

// stripExposeFromCoverProfileSoft removes session-generated expose lines from
// opts.CoverProfile using ExposeMaterializedList under genRoots. Errors are
// reported as warnings on w (stderr) and do not fail the test run.
func stripExposeFromCoverProfileSoft(opts core.Options, genRoots []string, w io.Writer) {
	if opts.CoverProfile == "" || len(genRoots) == 0 {
		return
	}
	if w == nil {
		w = os.Stderr
	}
	removed, err := core.StripExposeFromCoverProfile(opts.CoverProfile, genRoots)
	if err != nil {
		fmt.Fprintf(w, "warning: doctest: strip expose from coverprofile: %v\n", err)
		return
	}
	if removed > 0 && opts.Verbose {
		fmt.Fprintf(w, "doctest: stripped %d expose coverprofile line(s)\n", removed)
	}
}
