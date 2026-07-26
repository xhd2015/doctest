package build

import (
	"fmt"
	"io"
	"strings"
)

// vcsStatusMarker is the phrase go prints when -buildvcs=auto fails to stamp.
// See cmd/go internal/load setVCSError.
const vcsStatusMarker = "error obtaining VCS status"

// buildVCSStatusHint returns multi-line Error + hint guidance when goOutput
// indicates VCS stamping failed. Empty when the marker is absent.
//
// Doctest does not inject -buildvcs=false; users control stamping via GOFLAGS
// (or by fixing git trust). Shallow clone alone does not cause this error.
func buildVCSStatusHint(goOutput string) string {
	if !strings.Contains(goOutput, vcsStatusMarker) {
		return ""
	}
	return strings.Join([]string{
		"Error: go could not obtain VCS status (git failed while stamping; not caused by shallow clone alone)",
		"hint: set GOFLAGS=-buildvcs=false and re-run",
		"hint: CI (GitHub Actions) example:",
		"hint:   env:",
		"hint:     GOFLAGS: -buildvcs=false",
		"hint: or fix git trust: git config --global --add safe.directory '*'",
	}, "\n")
}

// FormatBuildVCSStatusHint is the exported form for callers outside this package
// (e.g. testbin) that append guidance to go build/test errors.
func FormatBuildVCSStatusHint(goOutput string) string {
	return buildVCSStatusHint(goOutput)
}

// printBuildVCSStatusHint writes Error/hint lines to w when goOutput matches.
// Color: Error red, hint gray (same policy as go-test timeout UX).
func printBuildVCSStatusHint(w io.Writer, goOutput string, style colorStyle) {
	if w == nil {
		return
	}
	msg := buildVCSStatusHint(goOutput)
	if msg == "" {
		return
	}
	printErrorHintLines(w, msg, style)
}

// printErrorHintLines prints multi-line Error:/hint: text with color when enabled.
func printErrorHintLines(w io.Writer, msg string, style colorStyle) {
	for _, line := range strings.Split(msg, "\n") {
		if line == "" {
			continue
		}
		if style.enabled {
			switch {
			case strings.HasPrefix(line, "Error:"):
				line = style.red(line)
			case strings.HasPrefix(line, "hint:"):
				line = style.gray(line)
			}
		}
		fmt.Fprintln(w, line)
	}
}
