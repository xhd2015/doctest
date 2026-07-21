package runner

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// TestProfilePrepareDotDotDotBypass is opt-in prepare-only profiling for the
// monorepo ./... path (DOCTEST_DEBUG=bypass-go-test=1).
//
//	DOCTEST_PROFILE_PREPARE=1 go test ./libdoc/runner -run TestProfilePrepareDotDotDotBypass \
//	  -cpuprofile=/tmp/doctest-prepare.pprof -count=1 -timeout 30m
func TestProfilePrepareDotDotDotBypass(t *testing.T) {
	if os.Getenv("DOCTEST_PROFILE_PREPARE") == "" {
		t.Skip("set DOCTEST_PROFILE_PREPARE=1 to run prepare-only profile")
	}
	// Module root: libdoc/runner -> ../..
	root := filepath.Clean(filepath.Join("..", ".."))
	if _, err := os.Stat(filepath.Join(root, "go.mod")); err != nil {
		t.Fatalf("module root %s: %v", root, err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Setenv("DOCTEST_DEBUG", "bypass-go-test=1")
	// Prefer cold prepare when requested.
	args := []string{"./..."}
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
