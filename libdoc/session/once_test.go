package session

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"syscall"
	"testing"
)

func setSession(t *testing.T, id string) {
	t.Helper()
	t.Setenv(DoctestSessionIDEnv, id)
	// Ensure syscall.Getenv sees it (same process env).
	if v, ok := syscall.Getenv(DoctestSessionIDEnv); !ok || v != id {
		t.Fatalf("syscall.Getenv(%s)=%q,%v want %q", DoctestSessionIDEnv, v, ok, id)
	}
}

func TestOnceMissingSession(t *testing.T) {
	t.Setenv(DoctestSessionIDEnv, "")
	processMemo = sync.Map{}
	_, err := Once(t, "k", func(t testing.TB, cacheDir string) (string, error) {
		return "x", nil
	})
	if !errors.Is(err, errMissingSession) && err == nil {
		t.Fatalf("err=%v", err)
	}
}

func TestOnceEmptyKey(t *testing.T) {
	setSession(t, "sess-empty-key")
	processMemo = sync.Map{}
	_, err := Once(t, "", func(t testing.TB, cacheDir string) (string, error) {
		return "x", nil
	})
	if !errors.Is(err, errEmptyKey) {
		t.Fatalf("err=%v", err)
	}
}

func TestOnceRunsFnOncePerSessionKey(t *testing.T) {
	setSession(t, "sess-once-1")
	processMemo = sync.Map{}
	var calls atomic.Int32
	fn := func(t testing.TB, cacheDir string) (string, error) {
		calls.Add(1)
		if cacheDir == "" {
			t.Fatal("empty cacheDir")
		}
		// Prove cacheDir is usable.
		if err := os.WriteFile(filepath.Join(cacheDir, "marker"), []byte("1"), 0o644); err != nil {
			return "", err
		}
		return filepath.Join(cacheDir, "marker"), nil
	}

	v1, err := Once(t, "cli", fn)
	if err != nil {
		t.Fatal(err)
	}
	v2, err := Once(t, "cli", fn)
	if err != nil {
		t.Fatal(err)
	}
	if v1 != v2 {
		t.Fatalf("values differ: %q vs %q", v1, v2)
	}
	if calls.Load() != 1 {
		t.Fatalf("fn calls=%d want 1", calls.Load())
	}
	if _, err := os.Stat(v1); err != nil {
		t.Fatalf("marker: %v", err)
	}
}

func TestOnceDifferentKeysIndependent(t *testing.T) {
	setSession(t, "sess-keys")
	processMemo = sync.Map{}
	var a, b atomic.Int32
	va, err := Once(t, "a", func(t testing.TB, cacheDir string) (string, error) {
		a.Add(1)
		return "A", nil
	})
	if err != nil {
		t.Fatal(err)
	}
	vb, err := Once(t, "b", func(t testing.TB, cacheDir string) (string, error) {
		b.Add(1)
		return "B", nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if va != "A" || vb != "B" {
		t.Fatalf("got %q %q", va, vb)
	}
	if a.Load() != 1 || b.Load() != 1 {
		t.Fatalf("calls a=%d b=%d", a.Load(), b.Load())
	}
}

func TestOncePropagatesError(t *testing.T) {
	setSession(t, "sess-err")
	processMemo = sync.Map{}
	boom := errors.New("boom")
	var calls atomic.Int32
	fn := func(t testing.TB, cacheDir string) (string, error) {
		calls.Add(1)
		return "", boom
	}
	_, err := Once(t, "fail", fn)
	if err == nil || err.Error() != "boom" {
		t.Fatalf("err=%v", err)
	}
	_, err = Once(t, "fail", fn)
	if err == nil || err.Error() != "boom" {
		t.Fatalf("second err=%v", err)
	}
	// Second call should not re-invoke fn after error was persisted.
	if calls.Load() != 1 {
		t.Fatalf("fn calls=%d want 1", calls.Load())
	}
}

func TestSlugify(t *testing.T) {
	cases := map[string]string{
		"cli":           "cli",
		"go-binary":     "go-binary",
		"a/b":           "a-b",
		"a//b":          "a-b",
		"  x  ":         "x",
		"foo:bar":       "foo-bar",
	}
	for in, want := range cases {
		if got := slugify(in); got != want {
			t.Errorf("slugify(%q)=%q want %q", in, got, want)
		}
	}
}
