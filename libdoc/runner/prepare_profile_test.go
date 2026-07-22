package runner

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestProfilePrepareDotDotDotBypass is opt-in prepare-only profiling for the
// monorepo ./... path (DOCTEST_DEBUG=bypass-go-test=1).
//
//	DOCTEST_PROFILE_PREPARE=1 go test ./libdoc/runner -run TestProfilePrepareDotDotDotBypass \
//	  -cpuprofile=/tmp/doctest-prepare.pprof -count=1 -timeout 30m
//
// Debug is injected only on a re-exec child via cmd.Env (never t.Setenv).
func TestProfilePrepareDotDotDotBypass(t *testing.T) {
	if os.Getenv("DOCTEST_PROFILE_PREPARE") == "" {
		t.Skip("set DOCTEST_PROFILE_PREPARE=1 to run prepare-only profile")
	}

	const childFlag = "DOCTEST_PROFILE_PREPARE_CHILD"
	if os.Getenv(childFlag) != "1" {
		// Re-exec so DOCTEST_DEBUG is process-initial for the child only.
		cmd := exec.Command(os.Args[0], "-test.run=^TestProfilePrepareDotDotDotBypass$", "-test.timeout=30m")
		base := make([]string, 0, len(os.Environ())+4)
		for _, e := range os.Environ() {
			k, _, _ := strings.Cut(e, "=")
			if k == "DOCTEST_DEBUG" || k == childFlag {
				continue
			}
			base = append(base, e)
		}
		cmd.Env = append(base, childFlag+"=1", "DOCTEST_DEBUG=bypass-go-test=1")
		if os.Getenv("DOCTEST_PROFILE_COLD") != "" {
			cmd.Env = append(cmd.Env, "DOCTEST_PROFILE_COLD=1")
		}
		// Keep profile prepare opt-in for the child.
		cmd.Env = append(cmd.Env, "DOCTEST_PROFILE_PREPARE=1")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			t.Fatalf("profile child: %v", err)
		}
		return
	}

	// Module root: libdoc/runner -> ../.. (absolute; no os.Chdir).
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(thisFile), "..", ".."))
	if _, err := os.Stat(filepath.Join(root, "go.mod")); err != nil {
		t.Fatalf("module root %s: %v", root, err)
	}
	// Prefer cold prepare when requested (child already has DOCTEST_DEBUG).
	pattern := filepath.Join(root, "...")
	args := []string{pattern}
	if os.Getenv("DOCTEST_PROFILE_COLD") != "" {
		args = append([]string{"--cold-cache", "--metrics-on"}, args...)
	} else {
		args = append([]string{"--metrics-on"}, args...)
	}
	if err := Test(args); err != nil {
		t.Fatalf("Test: %v", err)
	}
	runtime.GC()
}
