package build

import (
	"strings"
	"testing"
)

func TestBuildVCSStatusHint_detectsMarker(t *testing.T) {
	goOut := "error obtaining VCS status: exit status 128\n\tUse -buildvcs=false to disable VCS stamping.\n"
	got := buildVCSStatusHint(goOut)
	if got == "" {
		t.Fatal("expected non-empty hint")
	}
	if !strings.HasPrefix(got, "Error:") {
		t.Fatalf("want Error: prefix, got:\n%s", got)
	}
	if !strings.Contains(got, "GOFLAGS=-buildvcs=false") {
		t.Fatalf("want GOFLAGS guidance, got:\n%s", got)
	}
	if !strings.Contains(got, "GOFLAGS: -buildvcs=false") {
		t.Fatalf("want CI yml example, got:\n%s", got)
	}
	if !strings.Contains(got, "safe.directory") {
		t.Fatalf("want safe.directory alternative, got:\n%s", got)
	}
	if !strings.Contains(got, "shallow clone") {
		t.Fatalf("want shallow-not-cause note, got:\n%s", got)
	}
}

func TestBuildVCSStatusHint_unrelatedError(t *testing.T) {
	cases := []string{
		"",
		"undefined: foo",
		"[build failed]",
		"exit status 1",
		"Use -buildvcs=false to disable VCS stamping.", // alone without marker
	}
	for _, in := range cases {
		if got := buildVCSStatusHint(in); got != "" {
			t.Fatalf("input %q: want empty hint, got:\n%s", in, got)
		}
	}
}

func TestFormatBuildVCSStatusHint_exported(t *testing.T) {
	in := "error obtaining VCS status: exit status 128\n"
	if FormatBuildVCSStatusHint(in) != buildVCSStatusHint(in) {
		t.Fatal("FormatBuildVCSStatusHint must match buildVCSStatusHint")
	}
}

func TestCaptureVCSStatusFromBuffers(t *testing.T) {
	res := goTestJSONResult{
		buildOutputLines: []string{
			"# example.com/x",
			"error obtaining VCS status: exit status 128",
			"Use -buildvcs=false to disable VCS stamping.",
		},
	}
	captureVCSStatusFromBuffers(&res)
	if res.vcsStatusError == "" {
		t.Fatal("expected vcsStatusError from buffers")
	}
	if !strings.Contains(res.vcsStatusError, "GOFLAGS=-buildvcs=false") {
		t.Fatalf("unexpected vcsStatusError:\n%s", res.vcsStatusError)
	}
	// Idempotent when already set.
	prev := res.vcsStatusError
	captureVCSStatusFromBuffers(&res)
	if res.vcsStatusError != prev {
		t.Fatal("should not overwrite existing vcsStatusError")
	}
}
