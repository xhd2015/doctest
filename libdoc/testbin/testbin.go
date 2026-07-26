// Package testbin builds the doctest CLI once into a shared cache path for
// integration self-tests, so multi-package go test runs do not re-link the
// binary for every leaf (cold-start win).
package testbin

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/xhd2015/doctest/libdoc/build"
	"github.com/xhd2015/doctest/libdoc/core"
	"github.com/xhd2015/doctest/session"
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

	raw, err := session.Once(t, onceKey, func(t testing.TB, cacheDir string) (json.RawMessage, error) {
		// Durable output path (survives sessions); session.Once only serializes
		// concurrent first builds within a session.
		path, err := buildDoctest(absRoot)
		if err != nil {
			return nil, err
		}
		return json.Marshal(struct {
			Path string `json:"path"`
		}{Path: path})
	})
	if err != nil {
		t.Fatalf("build shared doctest binary: %v", err)
	}
	var got struct {
		Path string `json:"path"`
	}
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal session.Once value: %v (raw=%s)", err, raw)
	}
	if got.Path == "" {
		t.Fatalf("session.Once returned empty path (raw=%s)", raw)
	}
	return got.Path
}

func buildDoctest(absRoot string) (string, error) {
	sum := sha256.Sum256([]byte(absRoot))
	key := hex.EncodeToString(sum[:16])
	cacheBase, err := core.CacheHome()
	if err != nil {
		cacheBase = os.TempDir()
	}
	dir := filepath.Join(cacheBase, "doctest", "selftest-bin", key)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	bin := filepath.Join(dir, "doctest")

	// Cross-process lock: parallel trees may call Ensure with different
	// DOCTEST_SESSION_ID values and race on the same durable -o path.
	lockPath := filepath.Join(dir, ".build.lock")
	lf, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return "", fmt.Errorf("open build lock: %w", err)
	}
	defer lf.Close()
	if err := syscall.Flock(int(lf.Fd()), syscall.LOCK_EX); err != nil {
		return "", fmt.Errorf("flock build lock: %w", err)
	}
	defer syscall.Flock(int(lf.Fd()), syscall.LOCK_UN)

	// Always invoke go build so the toolchain rebuilds when sources change.
	// (A mere Stat on the -o path would pin a stale binary forever.)
	// Do not inject -buildvcs=false; users set GOFLAGS=-buildvcs=false if needed.
	args := []string{"build", "-o", bin, "./cmd/doctest"}
	cmd := exec.Command("go", args...)
	cmd.Dir = absRoot
	if out, err := cmd.CombinedOutput(); err != nil {
		msg := fmt.Sprintf("go %v: %v\n%s", args, err, out)
		if hint := build.FormatBuildVCSStatusHint(string(out)); hint != "" {
			return "", fmt.Errorf("%s\n%s", msg, hint)
		}
		return "", fmt.Errorf("%s", msg)
	}
	if st, err := os.Stat(bin); err != nil || st.IsDir() || st.Size() == 0 {
		return "", fmt.Errorf("go build produced no binary at %s", bin)
	}
	return bin, nil
}
