package testbin

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestEnsureSharedAndIdempotent(t *testing.T) {
	// session.Once requires DOCTEST_SESSION_ID.
	t.Setenv("DOCTEST_SESSION_ID", "testbin-unit-"+t.Name())

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
