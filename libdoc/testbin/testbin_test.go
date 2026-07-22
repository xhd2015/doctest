package testbin

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const testbinHelperEnv = "DOCTEST_TESTBIN_HELPER"

func TestEnsureSharedAndIdempotent(t *testing.T) {
	// session.Once requires DOCTEST_SESSION_ID on the process. Set it only on a
	// re-exec child via cmd.Env (never t.Setenv / process mutation).
	if os.Getenv(testbinHelperEnv) != "1" {
		runEnsureHelper(t)
		return
	}

	modRoot, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(modRoot, "cmd", "doctest")); err != nil {
		t.Fatalf("module root %s: %v", modRoot, err)
	}

	b1 := Ensure(t, modRoot)
	st1, err := os.Stat(b1)
	if err != nil || st1.Size() == 0 {
		t.Fatalf("binary missing or empty: %v", err)
	}

	start := time.Now()
	b2 := Ensure(t, modRoot)
	elapsed := time.Since(start)
	if b1 != b2 {
		t.Fatalf("paths differ: %q vs %q", b1, b2)
	}
	if elapsed > time.Second {
		t.Fatalf("second Ensure took %v; expected fast reuse", elapsed)
	}
}

func runEnsureHelper(t *testing.T) {
	t.Helper()
	cmd := exec.Command(os.Args[0], "-test.run=^TestEnsureSharedAndIdempotent$", "-test.v=false")
	// Strip any outer DOCTEST_SESSION_ID then replace with a unit-test sid.
	base := make([]string, 0, len(os.Environ()))
	for _, e := range os.Environ() {
		k, _, _ := strings.Cut(e, "=")
		if k == "DOCTEST_SESSION_ID" || k == testbinHelperEnv {
			continue
		}
		base = append(base, e)
	}
	cmd.Env = append(base,
		"DOCTEST_SESSION_ID=testbin-unit-"+t.Name(),
		testbinHelperEnv+"=1",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("testbin helper: %v\n%s", err, out)
	}
}
