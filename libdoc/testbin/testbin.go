// Package testbin builds the doctest CLI once into a shared cache path for
// integration self-tests, so multi-package go test runs do not re-link the
// binary for every leaf (cold-start win).
package testbin

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/xhd2015/doctest/libdoc/build"
	"github.com/xhd2015/doctest/libdoc/session"
)

// Ensure returns a path to a doctest binary built from moduleRoot/cmd/doctest.
//
// Within one doctest test session, concurrent leaves share a single build
// attempt via session.Once. The binary itself is written to a durable path
// keyed by absolute module root so go build is an up-to-date no-op when sources
// are unchanged:
//
//	$CACHE/doctest/selftest-bin/<sha256(absRoot)>/doctest
//
// Override with env DOCTEST_SELFTEST_BIN (read via syscall.Getenv) for CI.
func Ensure(t testing.TB, moduleRoot string) string {
	t.Helper()

	if override, ok := syscall.Getenv("DOCTEST_SELFTEST_BIN"); ok && override != "" {
		if st, err := os.Stat(override); err == nil && !st.IsDir() {
			return override
		}
		t.Fatalf("DOCTEST_SELFTEST_BIN=%q is not a usable binary", override)
	}

	absRoot, err := filepath.Abs(moduleRoot)
	if err != nil {
		t.Fatalf("abs module root: %v", err)
	}

	// Key is static slug-friendly: module path is hashed for uniqueness.
	sum := sha256.Sum256([]byte(absRoot))
	onceKey := "go-binary-" + hex.EncodeToString(sum[:8])

	bin, err := session.Once(t, onceKey, func(t testing.TB, cacheDir string) (string, error) {
		// Durable output path (survives sessions); session.Once only serializes
		// concurrent first builds within a session.
		return buildDoctest(absRoot)
	})
	if err != nil {
		t.Fatalf("build shared doctest binary: %v", err)
	}
	return bin
}

func buildDoctest(absRoot string) (string, error) {
	sum := sha256.Sum256([]byte(absRoot))
	key := hex.EncodeToString(sum[:16])
	cacheBase, err := os.UserCacheDir()
	if err != nil {
		cacheBase = os.TempDir()
	}
	dir := filepath.Join(cacheBase, "doctest", "selftest-bin", key)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	bin := filepath.Join(dir, "doctest")

	args := []string{"build", "-o", bin}
	if build.NeedsBuildVCSFlag(absRoot) {
		args = append(args, "-buildvcs=false")
	}
	args = append(args, "./cmd/doctest")
	cmd := exec.Command("go", args...)
	cmd.Dir = absRoot
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("go %v: %w\n%s", args, err, out)
	}
	return bin, nil
}
