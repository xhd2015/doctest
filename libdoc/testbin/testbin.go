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
	"sync"
	"syscall"
	"testing"

	"github.com/xhd2015/doctest/libdoc/build"
)

// processCache avoids re-entering go build within the same test process when
// multiple Setup calls share a package (rare) or helpers call Ensure twice.
var processCache sync.Map // absModuleRoot -> bin path

// Ensure returns a path to a doctest binary built from moduleRoot/cmd/doctest.
//
// The binary is written to a stable cache path keyed by absolute module root:
//
//	$CACHE/doctest/selftest-bin/<sha256(absRoot)>/doctest
//
// Concurrent callers for the same root serialize via flock; go build is then
// a no-op when sources are unchanged (shared -o path).
//
// Override with env DOCTEST_SELFTEST_BIN to force a prebuilt path (e.g. CI).
func Ensure(t testing.TB, moduleRoot string) string {
	t.Helper()

	if override := os.Getenv("DOCTEST_SELFTEST_BIN"); override != "" {
		if st, err := os.Stat(override); err == nil && !st.IsDir() {
			return override
		}
		t.Fatalf("DOCTEST_SELFTEST_BIN=%q is not a usable binary: %v", override, errStat(override))
	}

	absRoot, err := filepath.Abs(moduleRoot)
	if err != nil {
		t.Fatalf("abs module root: %v", err)
	}
	if cached, ok := processCache.Load(absRoot); ok {
		return cached.(string)
	}

	bin, err := ensureOnDisk(absRoot)
	if err != nil {
		t.Fatalf("build shared doctest binary: %v", err)
	}
	processCache.Store(absRoot, bin)
	return bin
}

func errStat(path string) error {
	_, err := os.Stat(path)
	return err
}

func ensureOnDisk(absRoot string) (string, error) {
	sum := sha256.Sum256([]byte(absRoot))
	key := hex.EncodeToString(sum[:16]) // 128-bit prefix is enough

	cacheBase, err := os.UserCacheDir()
	if err != nil {
		cacheBase = os.TempDir()
	}
	dir := filepath.Join(cacheBase, "doctest", "selftest-bin", key)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	bin := filepath.Join(dir, "doctest")
	lockPath := bin + ".lock"

	lock, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return "", fmt.Errorf("open lock: %w", err)
	}
	defer lock.Close()
	if err := syscall.Flock(int(lock.Fd()), syscall.LOCK_EX); err != nil {
		return "", fmt.Errorf("flock: %w", err)
	}
	defer func() { _ = syscall.Flock(int(lock.Fd()), syscall.LOCK_UN) }()

	// Re-check after lock: another process may have finished the build.
	if st, err := os.Stat(bin); err == nil && st.Mode().IsRegular() && st.Size() > 0 {
		// Still run go build so an outdated binary is refreshed when sources change.
		// With a stable -o path this is typically a sub-second up-to-date check.
	}

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
